import qrcode
import argparse

# Crie o analisador de argumentos
parser = argparse.ArgumentParser(description='Gere um QR Code a partir de um texto fornecido.')
parser.add_argument('data', type=str, help='O texto a ser convertido em QR Code')

# Analise os argumentos
args = parser.parse_args()

# Crie uma instância qr
qr = qrcode.QRCode(
    version=1,
    error_correction=qrcode.constants.ERROR_CORRECT_H,
    box_size=10,
    border=4,
)

# Adicione os dados ao qr
qr.add_data(args.data)
qr.make(fit=True)

# Imprima o QR Code no terminal
qr.print_ascii()
