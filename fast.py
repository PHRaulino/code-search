from fastapi import FastAPI
from typing import List
from sqlalchemy.orm import Session
from . import crud, models, schemas
from .database import engine, SessionLocal

app = FastAPI()

# Dependency
def get_db():
    try:
        db = SessionLocal()
        yield db
    finally:
        db.close()

# Create routes for all models
models_and_prefixes = [
    {"model": models.Empresa, "prefix": "empresas"},
    {"model": models.Fundo, "prefix": "fundos"}
]

for item in models_and_prefixes:
    prefix = item["prefix"]
    router = crud.CRUDRouter(schema=schemas.get_schema_for_model(item["model"]), db_model=item["model"], prefix=prefix)
    app.include_router(router, prefix=f"/{prefix}", tags=[prefix])

