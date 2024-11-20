from sqlalchemy.ext.asyncio import AsyncSession
from urllib.parse import urlparse
from services.redis_utils import RedisUtil
from dependency_injector.wiring import inject
import logging


class UncollectedUrlRepository:
    @inject
    def __init__(self, logger: logging.Logger, rdb_session: AsyncSession, redis_util: RedisUtil):
        self.logger = logger
        self.rdb_session = rdb_session
        self.redis_util = redis_util

    async def i_lpop_url_list(self) -> str | None:
        """
        URL 리스트에서 URL을 가져오는 함수
        :return: str | None
        """
        url = await self.redis_util.lpop_url_list()
        if url:
            self.logger.info(f"Get URL: {url} from URL list")
            return url
        else:
            self.logger.info("No URL in URL list")
            return None

    async def assign_to_domain_queue(self, url: str) -> None:
        """
        넘어온 URL 도메인만 추출하여 큐에 넣는 함수
        :param url: 큐에 넣을 URL
        :return: None
        """
        if urlparse(url).scheme != "https":
            self.logger.warning(f"Not https: {url}")
            return None

        domain: str = urlparse(url).hostname
        if domain is None:
            self.logger.warning(f"Invalid URL: {url}")
            return None

        self.logger.info(f"Assign to domain queue: {domain}")
        await self.redis_util.lpush_domain(domain_name=domain, url=url)
        await self.redis_util.sadd_domain(domain_list_name="domain_queue_list", domain=domain)
