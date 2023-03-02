import re

def valida_cnpj(cnpj: str) -> bool:
    # Remover caracteres não numéricos do CNPJ
    cnpj = re.sub(r'\D', '', cnpj)

    # Verificar se o CNPJ tem 14 dígitos
    if len(cnpj) != 14:
        return False

    # Verificar se todos os dígitos são iguais
    if cnpj == cnpj[0] * 14:
        return False

    # Verificar o primeiro dígito verificador
    peso = [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2]
    soma = sum(int(cnpj[i]) * peso[i] for i in range(12))
    if (soma % 11) < 2:
        digito_1 = 0
    else:
        digito_1 = 11 - (soma % 11)
    
    if int(cnpj[12]) != digito_1:
        return False

    # Verificar o segundo dígito verificador
    peso = [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2]
    soma = sum(int(cnpj[i]) * peso[i] for i in range(13))
    if (soma % 11) < 2:
        digito_2 = 0
    else:
        digito_2 = 11 - (soma % 11)
    
    if int(cnpj[13]) != digito_2:
        return False

    return True
