import threading
import sqlalchemy as db
from sqlalchemy.orm import sessionmaker

class DatabaseSingleton:
    _instance = None
    _engine = None
    _Session = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._engine = db.create_engine("postgresql://user:password@localhost/database")
            cls._Session = sessionmaker(bind=cls._engine, autoflush=False, expire_on_commit=False, scopefunc=threading.get_ident)
        return cls._instance

    def get_session(self):
        return self._Session()
