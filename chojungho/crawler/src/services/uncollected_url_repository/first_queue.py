from sqlalchemy.ext.asyncio import AsyncSession
from typing import AsyncIterator
from redis import asyncio as aioredis
from urllib.parse import urlparse
import logging


class FirstQueue:
    def __init__(self, logger: logging.Logger, rdb_session: AsyncSession, redis_db: AsyncIterator[aioredis.Redis]):
        self.logger = logger
        self.rdb_session = rdb_session
        self.redis_db = redis_db

    async def assign_to_domain_queue(self, url: str) -> None:
        """
        넘어온 URL 도메인만 추출하여 큐에 넣는 함수
        :param url: 큐에 넣을 URL
        :return: None
        """
        domain: str = urlparse(url).hostname
        if domain is None:
            return

        async with self.redis_db as redis:
            await redis.lpush(domain, url)
            await redis.sadd("domain_queue_list", domain)
