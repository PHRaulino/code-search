
# obtenha os nomes das colunas do cursor
col_names = [desc[0] for desc in cur.description]

# inicialize uma lista vazia para armazenar os resultados
results = []

# itere sobre as linhas do cursor e crie um dicionário para cada linha
for row in cur.fetchall():
    result_dict = {}

    # adicione cada valor da linha ao dicionário usando o nome da coluna como chave e o valor da coluna como valor
    for i in range(len(col_names)):
        result_dict[col_names[i]] = row[i]

    results.append(result_dict)

# feche o cursor e a conexão com o banco de dados
cur.close()
conn.close()

# exiba os resultados
print(results)
