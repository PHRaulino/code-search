from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, scoped_session

class DatabaseManager:
    _instance = None
    
    @classmethod
    def instance(cls):
        if cls._instance is None:
            cls._instance = cls()
        return cls._instance
    
    def __init__(self):
        self.engine = create_engine("sqlite:///./test.db")
        self.session = scoped_session(sessionmaker(bind=self.engine))
    
    def get_session(self):
        return self.session
