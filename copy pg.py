import subprocess

# Dados do banco de origem
origem_host = "localhost"
origem_port = 5432
origem_database = "banco_origem"
origem_user = "usuario_origem"
origem_password = "senha_origem"

# Dados do banco de destino
destino_host = "localhost"
destino_port = 5433
destino_database = "banco_destino"
destino_user = "usuario_destino"
destino_password = "senha_destino"

# Cria o comando para gerar um arquivo de backup do banco de origem
backup_command = f"pg_dump -h {origem_host} -p {origem_port} -U {origem_user} -w {origem_database} > backup.sql"
subprocess.call(backup_command, shell=True)

# Cria o comando para restaurar o arquivo de backup no banco de destino
restore_command = f"pg_restore -h {destino_host} -p {destino_port} -U {destino_user} -w -d {destino_database} backup.sql"
subprocess.call(restore_command, shell=True)

print("Replicação de dados realizada com sucesso!")
