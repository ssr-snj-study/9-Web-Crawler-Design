from dependency_injector import containers, providers
from services import UncollectedUrlRepository, RedisUtil


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # PostgreSQL 리소스
    postgres_session = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()

    # redis_util
    redis_util = providers.Factory(
        RedisUtil,
        logger=logger,
        redis_db=redis_client.provided,
    )

    # uncollected_url_repository
    uncollected_url_repository = providers.Factory(
        UncollectedUrlRepository,
        logger=logger,
        rdb_session=postgres_session.provided,
        redis_util=redis_util,
    )
