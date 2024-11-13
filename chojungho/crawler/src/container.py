from dependency_injector import containers, providers
from common import conf, setup_logging
from infrastructure.rdb.rdb_postgresql import AsyncEngine, get_pg_session
from infrastructure.nosql.redis_client import init_redis_pool
from services import Container as ServicesContainer


class Container(containers.DeclarativeContainer):
    wiring_config = containers.WiringConfiguration(
        packages=[
            "services",
        ],
    )

    _config = conf()
    config = providers.Configuration()
    config.from_dict(_config.dict())

    # log 의존성 주입
    logger = providers.Singleton(setup_logging)

    # PostgreSQL 리소스
    postgres_engine = providers.Resource(AsyncEngine, config=config)

    postgres_session = providers.Resource(
        get_pg_session,
        session_factory=postgres_engine.provided.session_factory,
    )

    # Redis 리소스
    redis_client = providers.Resource(
        init_redis_pool,
        host=config.REDIS_HOST,
        port=config.REDIS_PORT,
        password=config.REDIS_PASSWORD,
    )

    # Service
    services_container = providers.Container(
        ServicesContainer, logger=logger, postgres_session=postgres_session, redis_client=redis_client
    )
