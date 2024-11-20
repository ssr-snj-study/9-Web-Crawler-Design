from dependency_injector.wiring import Provide, inject
from . import UncollectedUrlRepository
from .crawling import run_crawling
import logging


@inject
async def run_crawler(
    logger: logging.Logger = Provide["services_container.logger"],
    uncollected_url_repository: UncollectedUrlRepository = Provide["services_container.uncollected_url_repository"],
) -> None:
    """
    미수집 URL 저장소 실행
    :param logger: 로거
    :param uncollected_url_repository: 미수집 저장소
    :return: None
    """
    logger.info("Uncollected Url Repository Started")
    while True:
        # todo 후면큐 구현후 전면큐 로직구현 예정
        url = await uncollected_url_repository.i_lpop_url_list()
        # 도메인큐(후면큐)
        if url:
            await uncollected_url_repository.assign_to_domain_queue(url)
            await run_crawling()
