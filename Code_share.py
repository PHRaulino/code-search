import qrcode
import argparse
import os

# Crie o analisador de argumentos
parser = argparse.ArgumentParser(description='Gere um QR Code a partir de um texto fornecido.')
parser.add_argument('data', type=str, nargs='?', default='', help='O texto a ser convertido em QR Code')

# Analise os argumentos
args = parser.parse_args()

data = args.data

# Verifique se algum argumento foi fornecido
if data == '':
    # Tente ler de um arquivo de texto
    try:
        with open('input.txt', 'r') as file:
            data = file.read().replace('\n', '')
    except FileNotFoundError:
        print("Por favor, forneça um argumento ou coloque os dados em um arquivo chamado 'input.txt'.")

# Crie uma instância qr
qr = qrcode.QRCode(
    version=1,
    error_correction=qrcode.constants.ERROR_CORRECT_H,
    box_size=10,
    border=4,
)

# Adicione os dados ao qr
qr.add_data(data)
qr.make(fit=True)

# Imprima o QR Code no terminal
qr.print_ascii()
