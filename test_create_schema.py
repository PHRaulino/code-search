from sqlalchemy import DDL, create_engine

ddl_gestao = DDL("CREATE SCHEMA IF NOT EXISTS gestao")
ddl_cadastro = DDL("CREATE SCHEMA IF NOT EXISTS cadastro")

with engine.connect() as conn:
    conn.execute(ddl_gestao)
    conn.execute(ddl_cadastro)
