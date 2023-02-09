from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session, sessionmaker
import threading

class DatabaseSingleton:
    engine_write = None
    engine_read = None
    session_factory_write = None
    session_factory_read = None
    session_write = None
    session_read = None

    @classmethod
    def initialize(cls, connection_string_write, connection_string_read):
        cls.engine_write = create_engine(connection_string_write)
        cls.engine_read = create_engine(connection_string_read)
        cls.session_factory_write = sessionmaker(bind=cls.engine_write)
        cls.session_factory_read = sessionmaker(bind=cls.engine_read)

    @classmethod
    def get_session_write(cls):
        if not cls.session_write:
            cls.session_write = scoped_session(cls.session_factory_write, scopefunc=threading.get_ident)
        return cls.session_write

    @classmethod
    def get_session_read(cls):
        if not cls.session_read:
            cls.session_read = scoped_session(cls.session_factory_read, scopefunc=threading.get_ident)
        return cls.session_read
