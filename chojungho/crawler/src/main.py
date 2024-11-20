import asyncio
from container import Container
from services import run_crawler


async def main() -> None:
    """
    메인 함수
    :return: None
    """
    container = Container()
    await container.init_resources()

    try:
        await run_crawler()  # 미수집 URL 저장소 실행
    finally:
        await container.shutdown_resources()


if __name__ == "__main__":
    asyncio.run(main())
