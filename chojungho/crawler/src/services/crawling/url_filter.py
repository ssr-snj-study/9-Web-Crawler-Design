from urllib.parse import urlparse
import aiohttp
import re
import logging


class UrlFilter:
    def __init__(
        self,
        logger: logging.Logger,
        session: aiohttp.ClientSession,
        exclude_extensions: list[str],
        allowed_types: list[str],
        exclusion_patterns: list[str],
    ):
        self._url = None
        self.logger = logger
        self.session = session
        self._exclude_extensions = exclude_extensions
        self._allowed_types = allowed_types
        self._exclusion_patterns = exclusion_patterns

    @property
    def url(self) -> str:
        return self._url

    @url.setter
    def url(self, url: str) -> None:
        self._url = url

    async def filter_by_extension(self) -> bool:
        """
        파일 확장자를 기반으로 URL을 필터링
        :return: 제외 대상(True)인지 여부
        """
        # URL 경로에서 확장자를 추출
        path = urlparse(self.url).path
        for ext in self._exclude_extensions:
            if path.lower().endswith(f".{ext}"):
                self.logger.info(f"URL {self.url} excluded by extension")
                return True
        return False

    async def filter_by_content_type(self) -> bool:
        """
        콘텐츠 타입을 확인하여 URL을 필터링
        :return: 제외 대상(True)인지 여부
        """
        try:
            async with aiohttp.ClientSession() as session:  # ClientSession 인스턴스 생성
                async with session.head(self.url, allow_redirects=True) as response:
                    if response.status >= 400:
                        self.logger.info(f"URL {self.url} excluded by status code {response.status}")
                        return True
                    content_type = response.headers.get("Content-Type", "")
                    for allowed in self._allowed_types:
                        if allowed in content_type:
                            return False
                    self.logger.info(f"URL {self.url} excluded by content type {content_type}")
                    return True
        except aiohttp.ClientError as e:
            self.logger.error(f"Error checking URL {self.url} content type: {e}")
            return True

    async def filter_by_exclusion_list(self) -> bool:
        """
        접근 제외 목록에 포함된 URL인지 확인
        :return: 제외 대상(True)인지 여부
        """
        for pattern in self._exclusion_patterns:
            if re.search(pattern, self.url):
                self.logger.info(f"URL {self.url} excluded by exclusion pattern {pattern}")
                return True
        return False

    async def is_crawlable_url(self) -> bool:
        """
        URL이 크롤링 가능한지 판단
        :return: 크롤링 가능(True) 여부
        """

        return (
            not await self.filter_by_extension()
            and not await self.filter_by_content_type()
            and not await self.filter_by_exclusion_list()
        )
