Olá, time!
Segue a entrega do case técnico conforme solicitado.

Visão Geral
A solução foi desenvolvida em Go, com foco em escalabilidade e modularidade. Mesmo utilizando uma única base de código, a arquitetura é orientada a domínios e preparada para migração futura para microserviços.
As funcionalidades principais de solicitação, reserva, compra de ingressos e adição de produtos estão disponíveis na aplicação e documentadas na collection do Insomnia.

Anexos
Arquivo ZIP com o executável da aplicação (cinetuber.exe) e o codigo fonte do projeto
Collection do Insomnia com todas as rotas organizadas por domínio
Link do repositório no GitHub com todo o código-fonte: desafio-tech-itau

Como utilizar a aplicação

Opção 1: Executar o binário 
.exe
 (recomendado)
Extraia o .zip em uma pasta de sua preferência.
Abra o terminal (CMD ou PowerShell) e navegue até a pasta extraída.
Execute (no git bash):
./cinetuber.exe
A aplicação será iniciada na porta 8080.
No Insomnia, importe a collection e teste os endpoints da API.
Opção 2: Rodar a aplicação via código-fonte (caso o .exe não funcione)
Instale o Go (versão 1.21 ou superior) via Central de Software ou pelo site oficial: https://go.dev/dl/
Clone o repositório:
git clone https://github.com/PHRaulino/desafio-tech-itau
cd desafio-tech-itau
Instale as dependências:
go mod tidy
Execute a aplicação:
go run cmd/main.go

Opção 3: Build manual do executável
Caso queira gerar o binário localmente (Windows):
go build -o cinetuber.exe cmd/main.go
./cinetuber.exe
 Observações Técnicas
O banco de dados utilizado é SQLite, criado automaticamente com os scripts de schema e seed embutidos no binário.
Todas as respostas da API seguem o padrão JSON com o campo "data" para facilitar o consumo.
A estrutura do projeto é modular e organizada por domínio (produtos, pedidos, sessões etc.).
Documentação OpenAPI disponível no repositório.
Fico à disposição para esclarecer qualquer dúvida.
Obrigado pela oportunidade!

Atenciosamente,
PH Raulino
phraulino@outlook.com
