from dependency_injector.wiring import Provide, inject
from . import UncollectedUrlRepository
import logging


@inject
async def run_crawler(
    logger: logging.Logger = Provide["services_container.logger"],
    uncollected_url_repository: UncollectedUrlRepository = Provide["services_container.uncollected_url_repository"],
):
    logger.info("Crawler Started")
    while True:
        # todo 후면큐 구현후 전면큐 로직구현 예정
        url = "https://www.naver.com"

        # 도메인큐(후면큐)
        await uncollected_url_repository.assign_to_domain_queue(url)
