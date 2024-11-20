import asyncio
from dependency_injector.wiring import inject, Provide
import logging
from services.crawling.worker import worker as crawling_worker


@inject
async def run(logger: logging.Logger = Provide["services_container.logger"]):
    logger.info("Crawler Started")
    task = [asyncio.create_task(crawling_worker(worker_id)) for worker_id in range(5)]
    await asyncio.gather(*task)
