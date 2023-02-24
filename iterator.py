from typing import List
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from database import get_db
from .crud_base import CRUDBase

def crud_router(models: List, prefix: str) -> APIRouter:
    router = APIRouter()
    for model in models:
        schema = model.__name__ + "Schema"
        crud = CRUDBase(model=model, schema=schema)
        router.include_router(
            crud.router,
            prefix=prefix + "/" + model.__tablename__,
            tags=[model.__name__],
            dependencies=[Depends(get_db)],
        )
    return router
