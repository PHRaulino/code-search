# IA First 🚀

Automatização Inteligente de Geração de Código com IA para Workflows no GitHub

---

## 📌 Visão Geral

O **IA First** é um orquestrador escrito em **Go** que integra **GitHub Issues** com **Inteligência Artificial (StackSpot AI)** para automatizar a geração de código, testes, documentação e validação seguindo conceitos como **Clean Architecture** e **Design Patterns**.

Através da abertura de Issues em um repositório, o sistema interpreta os pedidos dos usuários, executa múltiplos agentes de IA especializados e entrega um **Pull Request pronto para revisão** — tudo isso com execução paralela para máxima eficiência.

---

## 🎯 Objetivos do Projeto

- Automatizar a geração de código a partir de pedidos em linguagem natural ou estruturada.
- Padronizar entregas seguindo arquiteturas e práticas previamente definidas.
- Reduzir o tempo de desenvolvimento manual através de múltiplos agentes IA.
- Permitir fácil integração com qualquer repositório GitHub através de **GitHub Actions**.

---

## ⚙️ Tecnologias Utilizadas

- **Go (Golang)** — Linguagem principal do orquestrador.
- **StackSpot AI** — Geração de código, testes e documentação com IA.
- **GitHub API** — Integração para Issues, Pull Requests e repositórios.
- **GitHub Actions** — Automação e execução do pipeline CI/CD.
- **YAML** — Comunicação estruturada entre IA e orquestrador.

---

## 🔗 Fluxo de Funcionamento

1. **Abertura da Issue**
   - O usuário descreve o que deseja (funcionalidade, ajuste ou dúvida).
   - Pode ser linguagem natural ou preenchimento de um **Template de Issue**.

2. **Interpretação e Análise**
   - Um **Agente de Interpretação** da StackSpot AI transforma o pedido em um **YAML estruturado** contendo:
     - Nome da feature
     - Linguagem
     - Arquitetura sugerida
     - Padrões de projeto
     - Flags de teste e documentação

3. **Orquestração no Go**
   - O YAML é parseado pelo orquestrador Go.
   - Com base no YAML, o pipeline dinâmico é montado, chamando os agentes necessários:
     - Geração de código
     - Geração de testes
     - Geração de documentação

4. **Execução Paralela**
   - Todos os agentes são executados em paralelo usando **goroutines** e **channels**.

5. **Entrega**
   - O resultado consolidado é adicionado a um **Pull Request** no repositório, pronto para revisão.

---

## 🧩 Estrutura de Pastas (Proposta)
```text
ia-first/
├── cmd/
│   └── cli/
│       └── main.go          # CLI principal
├── internal/
│   ├── orchestrator.go      # Orquestração das tasks
│   ├── parser.go            # Leitura e parse do YAML
│   └── agents/
│       ├── generate.go      # Agente de Geração de Código
│       ├── test.go          # Agente de Geração de Testes
│       └── doc.go           # Agente de Geração de Documentação
├── pkg/
│   ├── github/              # Módulo de integração com GitHub API
│   └── stackspot/           # Módulo de integração com StackSpot AI
├── README.md
└── go.mod
```
---

## 🔑 Agentes Especializados

| Agente | Função |
|--------|--------|
| **Entendimento do Pedido** | Interpreta a descrição e transforma em YAML estruturado |
| **Geração de Código** | Produz o código-fonte seguindo o padrão definido |
| **Geração de Testes** | Cria testes unitários ou de integração |
| **Geração de Documentação** | Produz documentação de código ou README |

---

## 📝 Exemplo de YAML Gerado

```yaml
feature: "Adicionar endpoint de criação de usuário"
language: "Go"
architecture: "Clean Architecture"
patterns:
  - "Repository Pattern"
  - "DTO"
testing: true
documentation: true
```

---

🧠 Estratégias de Inteligência
	•	Canary Sources: Treinar e alimentar o StackSpot AI com projetos e padrões existentes.
	•	Pipeline Condicional: O YAML define dinamicamente quais agentes serão executados.
	•	Paralelismo: Execução simultânea para reduzir tempo total.
	•	Fallback e Validação: Possibilidade de reprocessar ou pedir mais contexto em caso de falha.

---
🚀 Próximos Passos
	1.	Montar o esqueleto do CLI em Go.
	2.	Criar o parser de YAML.
	3.	Implementar o primeiro Agente IA (mock ou real).
	4.	Implementar o paralelismo com goroutines.
	5.	Integrar com GitHub API para abertura de Pull Requests.
	6.	Conectar com o StackSpot AI para chamadas reais.