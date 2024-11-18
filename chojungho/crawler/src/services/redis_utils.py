from typing import AsyncIterator
from redis import asyncio as aioredis
import logging


class RedisUtil:
    def __init__(self, logger: logging.Logger, redis_db: AsyncIterator[aioredis.Redis]):
        self.logger = logger
        self.redis_db = redis_db

    async def lpush_domain(self, domain_name: str, url: str) -> None:
        """
        도메인 큐에 URL 넣기
        :param domain_name: 도메인 이름
        :param url: URL
        :return: None
        """
        async with self.redis_db as redis:
            await redis.lpush(domain_name, url)
            self.logger.info(f"Push URL: {url} to domain: {domain_name}")

    async def sadd_domain(self, domain_list_name: str, domain: str) -> None:
        """
        도메인 목록에 도메인 추가
        :param domain_list_name: 도메인 목록 이름
        :param domain: 도메인
        :return: None
        """
        async with self.redis_db as redis:
            await redis.sadd(domain_list_name, domain)
            self.logger.info(f"Add domain: {domain} to domain list: {domain_list_name}")

    async def get_url_to_domain_queue(self, domain: bytes) -> str | None:
        """
        domain queue에서 url 가져오는 함수
        :param domain: 도메인
        :return: str | None
        """
        async with self.redis_db as redis:
            url = await redis.lpop(domain)
            url = url.decode("utf-8") if url else None
            self.logger.info(f"Get URL: {url} from domain: {domain}")
            return url

    async def get_domain_list(self, domain_name: str) -> set:
        """
        도메인 목록 가져오기
        :return: set
        """
        async with self.redis_db as redis:
            return await redis.smembers(domain_name)
