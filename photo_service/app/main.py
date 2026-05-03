from __future__ import annotations

import uuid

from fastapi import FastAPI, HTTPException, Request, Response
from sqlalchemy import LargeBinary, Column, DateTime, func, select
from sqlalchemy.dialects.postgresql import UUID as PG_UUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, Session

from app.db import create_shard_engines
from app.sharding import shard_index


class Base(DeclarativeBase):
    pass


class Image(Base):
    __tablename__ = "images"

    id: Mapped[uuid.UUID] = mapped_column(PG_UUID(as_uuid=True), primary_key=True)
    data: Mapped[bytes] = mapped_column(LargeBinary, nullable=False)
    created_at: Mapped[object] = mapped_column(DateTime(timezone=True), server_default=func.now())


app = FastAPI(title="Photo Service", version="1.0.0")
shards = create_shard_engines()


@app.on_event("startup")
def _startup() -> None:
    # Ensure table exists on all shards (demo-friendly)
    for eng in shards.engines:
        Base.metadata.create_all(bind=eng)


def _engine_for(image_id: uuid.UUID):
    return shards.engines[shard_index(image_id)]


@app.get("/health")
def health() -> dict:
    return {"status": "ok"}


@app.put("/images/{image_id}")
async def put_image(image_id: uuid.UUID, request: Request) -> dict:
    data = await request.body()
    if not data:
        raise HTTPException(status_code=400, detail="Image body is required")

    eng = _engine_for(image_id)
    with Session(eng) as s:
        existing = s.get(Image, image_id)
        if existing is None:
            s.add(Image(id=image_id, data=data))
        else:
            existing.data = data
        s.commit()
    return {"id": str(image_id), "shard": shard_index(image_id)}


@app.post("/images")
async def create_image(request: Request) -> dict:
    data = await request.body()
    if not data:
        raise HTTPException(status_code=400, detail="Image body is required")

    image_id = uuid.uuid4()
    eng = _engine_for(image_id)
    with Session(eng) as s:
        s.add(Image(id=image_id, data=data))
        s.commit()
    return {"id": str(image_id), "shard": shard_index(image_id)}


@app.get("/images/{image_id}")
def get_image(image_id: uuid.UUID) -> Response:
    eng = _engine_for(image_id)
    with Session(eng) as s:
        img = s.get(Image, image_id)
        if img is None:
            raise HTTPException(status_code=404, detail="Image not found")
        return Response(
            content=img.data,
            media_type="application/octet-stream",
            headers={"X-Photo-Shard": str(shard_index(image_id))},
        )


@app.get("/images/{image_id}/debug")
def debug_image(image_id: uuid.UUID) -> dict:
    idx = shard_index(image_id)
    return {"id": str(image_id), "shard": idx, "database": f"photos{idx}"}


@app.delete("/images/{image_id}", status_code=204)
def delete_image(image_id: uuid.UUID) -> Response:
    eng = _engine_for(image_id)
    with Session(eng) as s:
        img = s.get(Image, image_id)
        if img is None:
            return Response(status_code=204)
        s.delete(img)
        s.commit()
    return Response(status_code=204)

