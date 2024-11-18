from pydantic_settings import BaseSettings
from os import environ


class Config(BaseSettings):
    POSTGRES_USER: str = ""
    POSTGRES_PASSWORD: str = ""
    POSTGRES_HOST: str = ""
    POSTGRES_PORT: int
    POSTGRES_DB: str = ""
    REDIS_HOST: str = ""
    REDIS_PORT: int
    REDIS_PASSWORD: str = ""

    class Config:
        env_file = ".env"


class LocalConfig(Config):
    DEBUG: bool = True
    SQL_PRINT: bool = True

    POSTGRES_SERVER: str = ""
    POSTGRES_SCHEMA: str = "public"

    REDIS_SERVER: str = ""


def conf():
    c = dict(local=LocalConfig)
    return c[environ.get("API_ENV", "local")]()
