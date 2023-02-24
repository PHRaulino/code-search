
from typing import List, Type, Union
from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy.orm import Session
from sqlalchemy.exc import NoResultFound
from sqlalchemy import update

from database import Base

class CRUDBase:
    def __init__(self, model: Type[Base], schema: Type[BaseModel]):
        self.model = model
        self.schema = schema

        self.router = APIRouter()
        self._register_routes()

    def _register_routes(self):
        @self.router.get("/", response_model=List[self.schema])
        def read_all(db: Session = Depends()):
            return db.query(self.model).all()

        @self.router.get("/{id}", response_model=self.schema)
        def read_one(id: int, db: Session = Depends()):
            instance = db.query(self.model).get(id)
            if not instance:
                raise HTTPException(status_code=404, detail="Resource not found")
            return instance

        @self.router.post("/", response_model=self.schema)
        def create(payload: self.schema, db: Session = Depends()):
            instance = self.model(**payload.dict())
            db.add(instance)
            db.commit()
            db.refresh(instance)
            return instance

        @self.router.put("/{id}", response_model=self.schema)
        def update(id: int, payload: self.schema, db: Session = Depends()):
            instance = db.query(self.model).get(id)
            if not instance:
                raise HTTPException(status_code=404, detail="Resource not found")
            for field, value in payload:
                setattr(instance, field, value)
            db.add(instance)
            db.commit()
            db.refresh(instance)
            return instance

        @self.router.delete("/{id}", response_model=self.schema)
        def delete(id: int, db: Session = Depends()):
            instance = db.query(self.model).get(id)
            if not instance:
                raise HTTPException(status_code=404, detail="Resource not found")
            db.delete(instance)
            db.commit()
            return instance
