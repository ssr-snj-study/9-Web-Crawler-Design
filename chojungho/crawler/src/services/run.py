from dependency_injector.wiring import Provide, inject
from . import FirstQueue
import logging


@inject
async def run_crawler(
    logger: logging.Logger = Provide["services_container.logger"],
    first_queue_service: FirstQueue = Provide["services_container.first_queue_service"],
):
    logger.info("Crawler Started")
    while True:
        # todo 후면큐 구현후 전면큐 로직구현 예정
        url = "https://www.naver.com"

        # 도메인큐(후면큐)
        await first_queue_service.assign_to_domain_queue(url)
