from infrastructure.schema.urls import Urls
from sqlalchemy import exists, select, insert
import logging
from contextlib import AbstractAsyncContextManager
from typing import Callable
from sqlalchemy.ext.asyncio import AsyncSession
from dependency_injector.wiring import inject


class SqlAlchemyUtil:
    @inject
    def __init__(self, logger: logging, rdb_session: Callable[..., AbstractAsyncContextManager[AsyncSession]]):
        self.logger = logger
        self.rdb_session = rdb_session

    async def check_url(self, url: str) -> bool:
        """
        URL 중복 체크
        :param url: URL
        :return: 중복 여부
        """
        async with self.rdb_session() as session:
            result = await session.execute(select(exists().where(Urls.url == url)))
            is_exists = result.scalar()

        if is_exists:
            self.logger.info(f"URL {url} is already collected")
            return True
        else:
            self.logger.info(f"URL {url} is not collected")
            return False

    async def save_url(self, url: str) -> None:
        """
        URL 저장
        :param url: URL
        """
        async with self.rdb_session() as session:
            await session.execute(insert(Urls).values(url=url))
            await session.commit()
            self.logger.info(f"URL {url} saved")
