from typing import List, TypeVar
from sqlalchemy.orm import Session
from .database import SessionLocal

ModelType = TypeVar("ModelType")

class CRUD:
    def __init__(self, model: ModelType):
        self.model = model
    
    def get_all(self, skip: int = 0, limit: int = 100, db: Session = None) -> List[ModelType]:
        if not db:
            db = SessionLocal()
        return db.query(self.model).offset(skip).limit(limit).all()
