import psycopg2

# Dados do banco de origem
origem_host = "localhost"
origem_port = 5432
origem_database = "banco_origem"
origem_user = "usuario_origem"
origem_password = "senha_origem"

# Dados do banco de destino
destino_host = "localhost"
destino_port = 5432
destino_database = "banco_destino"
destino_user = "usuario_destino"
destino_password = "senha_destino"

try:
    # Conexão com o banco de origem
    origem_conn = psycopg2.connect(
        host=origem_host,
        port=origem_port,
        database=origem_database,
        user=origem_user,
        password=origem_password
    )
    origem_cursor = origem_conn.cursor()

    # Conexão com o banco de destino
    destino_conn = psycopg2.connect(
        host=destino_host,
        port=destino_port,
        database=destino_database,
        user=destino_user,
        password=destino_password
    )
    destino_cursor = destino_conn.cursor()

    # Cria uma cópia do banco de origem para o destino
    origem_cursor.execute(f"CREATE DATABASE {destino_database} TEMPLATE {origem_database}")

    # Confirma a transação
    origem_conn.commit()

    # Encerra a conexão com o banco de origem
    origem_cursor.close()
    origem_conn.close()

    # Encerra a conexão com o banco de destino
    destino_cursor.close()
    destino_conn.close()

    print("Cópia do banco de dados realizada com sucesso!")

except (Exception, psycopg2.DatabaseError) as error:
    print(error)
