import re

def valida_cpf(cpf: str) -> bool:
    # Remover caracteres não numéricos do CPF
    cpf = re.sub(r'\D', '', cpf)

    # Verificar se o CPF tem 11 dígitos
    if len(cpf) != 11:
        return False

    # Verificar se todos os dígitos são iguais
    if cpf == cpf[0] * 11:
        return False

    # Verificar o primeiro dígito verificador
    soma = sum(int(cpf[i]) * (10 - i) for i in range(9))
    if (soma % 11) < 2:
        digito_1 = 0
    else:
        digito_1 = 11 - (soma % 11)
    
    if int(cpf[9]) != digito_1:
        return False

    # Verificar o segundo dígito verificador
    soma = sum(int(cpf[i]) * (11 - i) for i in range(10))
    if (soma % 11) < 2:
        digito_2 = 0
    else:
        digito_2 = 11 - (soma % 11)
    
    if int(cpf[10]) != digito_2:
        return False

    return True
