Aqui está o prompt reestruturado com boas práticas de prompt engineering, otimizado para um fluxo de desenvolvimento full stack com Python + Angular, e claro, pensado para uso com uma IA como o Copilot Agent ou ChatGPT com foco em geração de código:

⸻

🧠 Prompt: Criação de Sistema Full Stack Python + Angular com Automação Playwright

Você é um assistente especializado em desenvolvimento full stack com Python, FastAPI, Angular e automação com Playwright.

Sua tarefa é gerar dois repositórios separados, além de estruturar recursos essenciais para rodar uma automação com interface gráfica web e lógica 100% em Python.

⸻

✅ Objetivo da Aplicação
Criar um executável Python que:
	•	Sirva uma interface web Angular por meio de um servidor FastAPI.
	•	Abra automaticamente o navegador com duas abas:
	•	Aba 1: Carrega a interface Angular.
	•	Aba 2: Abre a página que será raspada.
	•	Permita que o usuário:
	•	Faça login manualmente uma única vez na aba de raspagem.
	•	Configure e execute a automação pela interface Angular.
	•	Observe logs e status da automação na própria página Angular.
	•	Mantenha o navegador aberto enquanto o executável estiver em execução.
	•	Pare a aplicação se o navegador for fechado.

⸻

🛠️ Instruções Técnicas — Gerar Código e Estruturas para:
1. Organização dos Repositórios
	•	Frontend (Angular):
	•	Apenas inicializar o projeto com Angular CLI.
	•	Build será feito e servido pelo backend Python.
	•	Backend (Python):
	•	FastAPI com rota para servir os arquivos estáticos do build Angular.
	•	Endpoints REST ou WebSocket para acionar automação e enviar status para interface.

2. FastAPI
	•	Rota /app servindo o index.html e arquivos estáticos do Angular build.
	•	Endpoint /run-automation que chama a automação Playwright.
	•	Endpoint /status para logs em tempo real (opcional: WebSocket).

3. Automação com Playwright
	•	Abrir navegador Chromium com duas abas:

- Aba 1 → http://localhost:<porta>/app   (interface Angular)
- Aba 2 → Página-alvo para raspagem


	•	Manter sessão ativa na aba 2 enquanto o browser estiver aberto.
	•	Automação realiza ações de raspagem somente quando chamada pelo frontend via FastAPI.

4. Build e Execução
	•	Incluir um Makefile com comandos como:
	•	make build-frontend
	•	make build-backend
	•	make build-exe
	•	Gerar executável Python com PyInstaller ou equivalente.

⸻

⚠️ Restrições Importantes
	•	Não implementar sistema de credenciais. O login da página de raspagem será feito manualmente pelo usuário.
	•	A aplicação não deve fechar o browser automaticamente.
	•	Somente utilize Python + Angular + Playwright + FastAPI. Nada de Electron ou Node adicional.

⸻

📦 Entregáveis Esperados
	1.	Estrutura de pastas dos dois repositórios.
	2.	Código inicial para backend FastAPI.
	3.	Script para automação Playwright com funcionalidades descritas.
	4.	Makefile com targets funcionais para build.
	5.	Configuração básica de comunicação Angular → FastAPI.

⸻

🎯 Gere o código necessário nos formatos adequados, com explicações breves onde necessário. Use boas práticas de clean architecture, modularização e documentação interna.

⸻

Se quiser que eu gere esse prompt num formato pronto pra colar no Copilot ou adaptar, é só avisar!