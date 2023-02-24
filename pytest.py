from in_memory_db.sqlite import SQLite
from sqlalchemy.orm import Session
from sqlalchemy import create_engine

@pytest.fixture(scope="session")
def db_engine():
    engine = create_engine("sqlite:///:memory:")
    Base.metadata.create_all(bind=engine)
    yield engine

@pytest.fixture(scope="function")
def db(db_engine):
    conn = db_engine.connect()
    trans = conn.begin()
    session = Session(bind=conn)

    # Substitua a linha abaixo pelo nome do seu schema
    session.execute("CREATE SCHEMA IF NOT EXISTS cadastro")
    session.execute("CREATE SCHEMA IF NOT EXISTS gestao")

    yield session

    session.close()
    trans.rollback()
    conn.close()

@pytest.fixture(scope="function")
def client(db):
    app.dependency_overrides[database.get_session_write] = lambda: db

    # Insira aqui as tabelas do seu banco de dados
    db.execute("CREATE TABLE cadastro.TBA (id INTEGER PRIMARY KEY, name TEXT, email TEXT)")
    db.execute("CREATE TABLE gestao.TBB (id INTEGER PRIMARY KEY, name TEXT, email TEXT)")

    with TestClient(app) as c:
        yield c
