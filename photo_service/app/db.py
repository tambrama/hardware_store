from __future__ import annotations

import os
from dataclasses import dataclass

from sqlalchemy import create_engine
from sqlalchemy.engine import Engine


@dataclass(frozen=True)
class ShardEngines:
    engines: tuple[Engine, Engine, Engine, Engine]


def _require_env(name: str) -> str:
    v = os.getenv(name)
    if not v:
        raise RuntimeError(f"Missing required env: {name}")
    return v


def create_shard_engines() -> ShardEngines:
    dsns = (
        _require_env("PHOTO_SHARD_0_DSN"),
        _require_env("PHOTO_SHARD_1_DSN"),
        _require_env("PHOTO_SHARD_2_DSN"),
        _require_env("PHOTO_SHARD_3_DSN"),
    )
    engines = tuple(create_engine(dsn, pool_pre_ping=True) for dsn in dsns)
    return ShardEngines(engines=engines)  # type: ignore[arg-type]

