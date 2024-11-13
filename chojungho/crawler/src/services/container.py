from dependency_injector import containers, providers
from services import UncollectedUrlRepository


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # PostgreSQL 리소스
    postgres_session = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()

    # uncollected_url_repository
    uncollected_url_repository = providers.Factory(
        UncollectedUrlRepository,
        logger=logger,
        rdb_session=postgres_session.provided,
        redis_db=redis_client.provided,
    )
