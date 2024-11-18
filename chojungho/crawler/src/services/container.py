from dependency_injector import containers, providers
from services import UncollectedUrlRepository, RedisUtil, SqlAlchemyUtil
from services.crawling.url_filter import UrlFilter
from common import FilterRule
import aiohttp


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # filter url
    filter_rule: FilterRule = providers.Singleton(FilterRule)

    # PostgreSQL 리소스
    postgres_engine = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()

    # aiohttp.ClientSession 리소스
    aiohttp_session = providers.Resource(aiohttp.ClientSession)

    # redis_util
    redis_util = providers.Factory(
        RedisUtil,
        logger=logger,
        redis_db=redis_client.provided,
    )

    # sqlalchemy_utl
    sqlalchemy_util = providers.Factory(
        SqlAlchemyUtil,
        logger=logger,
        rdb_session=postgres_engine.provided.get_pg_session,
    )

    # uncollected_url_repository
    uncollected_url_repository = providers.Factory(
        UncollectedUrlRepository,
        logger=logger,
        rdb_session=postgres_engine.provided.get_pg_session,
        redis_util=redis_util,
    )

    # url_filter
    url_filter = providers.Factory(
        UrlFilter,
        logger=logger,
        session=aiohttp_session,
        exclude_extensions=filter_rule().EXCLUDE_EXTENSIONS,
        allowed_types=filter_rule().ALLOWED_TYPES,
        exclusion_patterns=filter_rule().EXCLUSION_PATTERNS,
    )
