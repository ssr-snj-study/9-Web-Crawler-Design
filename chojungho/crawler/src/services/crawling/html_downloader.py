import aiohttp
import asyncio
from bs4 import BeautifulSoup
from dependency_injector.wiring import Provide, inject
import logging


@inject
async def download_html(
    url: str, logger: logging.Logger = Provide["services_container.logger"]
) -> BeautifulSoup | None:
    """
    HTML 다운로더
    :param url: 도메인큐에서 넘어온 URL
    :param logger: 로거
    :return:
    """
    timeout = aiohttp.ClientTimeout(total=5)
    async with aiohttp.ClientSession(timeout=timeout) as session:
        try:
            async with session.get(url) as response:
                if response.status != 200:
                    raise aiohttp.ClientResponseError(
                        response.request_info,
                        response.history,
                        status=response.status,
                        message=response.reason,
                        headers=response.headers,
                    )
                html = await response.text()
                soup = BeautifulSoup(html, "html.parser")
                return soup
        except aiohttp.ClientResponseError as e:
            logger.error(f"HTTP error: {e.status} {e.message} for {url}")
        except aiohttp.ClientConnectionError as e:
            logger.error(f"Connection error: {e} for {url}")
        except aiohttp.ClientPayloadError as e:
            logger.error(f"Payload error: {e} for {url}")
        except asyncio.TimeoutError:
            logger.error(f"Timeout for {url}")


# test
# async def main():
#     urls = ["https://www.google.com", "https://www.naver.com"]
#     tasks = [download_html(url) for url in urls]
#     await asyncio.gather(*tasks)
#
#
# # 실행
# print(asyncio.run(main()))
