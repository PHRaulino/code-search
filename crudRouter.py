from typing import Any, Dict, List, Optional, Type, TypeVar
from fastapi import APIRouter, HTTPException, Path
from pydantic import BaseModel
from sqlalchemy.orm import Session
from .database import Base


def CRUDRouter(db_model: Type[Base], schema: Type[BaseModel], prefix: str = ""):
    router = APIRouter()

    # Create
    @router.post("/", response_model=schema)
    def create_item(item: schema, db: Session = Depends(get_db)):
        db_item = db_model(**item.dict())
        db.add(db_item)
        db.commit()
        db.refresh(db_item)
        return db_item

    # Read all
    @router.get("/", response_model=List[schema])
    def read_items(skip: int = 0, limit: int = 100, db: Session = Depends(get_db)):
        items = db.query(db_model).offset(skip).limit(limit).all()
        return items

    # Read one
    @router.get("/{item_id}", response_model=schema)
    def read_item(item_id: int, db: Session = Depends(get_db)):
        item = db.query(db_model).filter(db_model.id == item_id).first()
        if not item:
            raise HTTPException(status_code=404, detail=f"{db_model.__name__} not found")
        return item

    # Update
    @router.put("/{item_id}", response_model=schema)
    def update_item(item_id: int, item: schema, db: Session = Depends(get_db)):
        db_item = db.query(db_model).filter(db_model.id == item_id).first()
        if not db
