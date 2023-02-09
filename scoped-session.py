from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session

engine = create_engine("sqlite:///./test.db")
session = scoped_session(sessionmaker(bind=engine), scopefunc=<function_that_returns_a_unique_id>)
