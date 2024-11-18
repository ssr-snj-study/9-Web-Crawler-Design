from bs4 import BeautifulSoup
import asyncio
import logging
from dependency_injector.wiring import inject, Provide
from services.redis_utils import RedisUtil
from services.crawling import download_html


@inject
async def worker(
    worker_id: int,
    logger: logging.Logger = Provide["services_container.logger"],
    redis_util: RedisUtil = Provide["services_container.redis_util"],
):
    while True:
        # 도메인 목록 가져오기
        domain_list: set[bytes] = await redis_util.get_domain_list("domain_queue_list")
        if not domain_list:
            logger.info("No domain queue")
            await asyncio.sleep(5)
            continue

    # 각 도메인 큐에서 URL을 순서대로 가져와 처리
    for domain in domain_list:
        url: str = await redis_util.get_url_to_domain_queue(domain)
        if url:
            logger.info(f"Worker {worker_id} processing URL: {url} from domain: {domain}")
            # todo: URL 처리
            soup: BeautifulSoup = await download_html(url)
            print(soup.title.text)
            await asyncio.sleep(1)
        else:
            logger.info(f"No URL in queue for domain: {domain}")
