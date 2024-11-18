import asyncio
from dependency_injector.wiring import inject, Provide
import logging
from services.crawling.worker import worker


@inject
async def run(logger: logging.Logger = Provide["services_container.logger"]):
    logger.info("Crawler Started")
    while True:
        task = [asyncio.create_task(worker(worker_id)) for worker_id in range(1)]
        await asyncio.gather(*task)
