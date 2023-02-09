import time
from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session

class DatabaseManager:
    _instance = None
    _session_expiration = 3 * 60 * 60
    
    @classmethod
    def instance(cls):
        if cls._instance is None:
            cls._instance = cls()
        return cls._instance
    
    def __init__(self):
        self.engine = create_engine("sqlite:///./test.db")
        self.session = scoped_session(sessionmaker(bind=self.engine))
        self.last_update = 0
    
    def get_session(self):
        current_time = time.time()
        if current_time - self.last_update > self._session_expiration:
            self.session.remove()
            self.session = scoped_session(sessionmaker(bind=self.engine))
            self.last_update = current_time
        return self.session
