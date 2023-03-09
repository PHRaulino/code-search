def snake_to_pascal_case(input_str):
    words = input_str.split("_")
    capitalized_words = [word.capitalize() for word in words]
    return "".join(capitalized_words)

# Exemplo de uso
input_string = "exemplo_de_texto_em_snake_case"
output_string = snake_to_pascal_case(input_string)
print(output_string)  # Saída: "ExemploDeTextoEmSnakeCase"
