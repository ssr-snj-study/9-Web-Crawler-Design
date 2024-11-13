from dependency_injector import containers, providers
from services import FirstQueue


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # PostgreSQL 리소스
    postgres_session = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()

    # uncollected_url_repository
    first_queue_service = providers.Factory(
        FirstQueue,
        logger=logger,
        rdb_session=postgres_session.provided,
        redis_db=redis_client.provided,
    )
