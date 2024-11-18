from bs4 import BeautifulSoup
import logging
from dependency_injector.wiring import inject, Provide
from services.redis_utils import RedisUtil
from services.crawling import download_html
from services.crawling.url_parser import extract_urls_from_soup


@inject
async def worker(
    worker_id: int,
    logger: logging.Logger = Provide["services_container.logger"],
    redis_util: RedisUtil = Provide["services_container.redis_util"],
):
    """
    크롤링 워커
    :param worker_id: Worker ID
    :param logger: 로거
    :param redis_util: 레디스 유틸
    :return:
    """
    # 도메인 목록 가져오기
    domain_list: set[bytes] = await redis_util.get_domain_list("domain_queue_list")
    if not domain_list:
        logger.info("No domain queue")
        return None

    # 각 도메인 큐에서 URL을 순서대로 가져와 처리
    for domain in domain_list:
        url: str = await redis_util.get_url_to_domain_queue(domain)
        if url:
            logger.info(f"Worker {worker_id} processing URL: {url} from domain: {domain}")

            # URL 다운로드
            soup: BeautifulSoup = await download_html(url)

            # todo: 중복 콘텐츠확인

            # URL 추출
            urls: set[str] = await extract_urls_from_soup(soup, domain.decode())

            # todo: 추출한 URL필터 작업 진행

            # todo: 방문한 URL체크

            # todo: 미방문 URL 미수집, URL저장소에 저장
            return urls
        else:
            logger.info(f"No URL in queue for domain: {domain}")
