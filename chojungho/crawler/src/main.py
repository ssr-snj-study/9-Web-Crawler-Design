import asyncio
from container import Container
from services.run import run_crawler


async def main():
    container = Container()
    await container.init_resources()

    try:
        await run_crawler()
    finally:
        await container.shutdown_resources()


if __name__ == "__main__":
    asyncio.run(main())
