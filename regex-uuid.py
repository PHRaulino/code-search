import re

def is_uuid(string):
    pattern = r'^[0-9a-f]{8}-?[0-9a-f]{4}-?4[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$'
    return re.match(pattern, string.lower()) is not None

# Exemplo de uso:
uuid_string = "550e8400-e29b-41d4-a716-446655440000"
if is_uuid(uuid_string):
    print("É um UUID válido.")
else:
    print("Não é um UUID válido.")
