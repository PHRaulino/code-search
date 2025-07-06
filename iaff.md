# IA First 🚀

Automatização Inteligente de Geração de Código com IA para Workflows no GitHub

---

## 📌 Visão Geral

O **IA First** é um orquestrador que integra **GitHub Issues** com **Inteligência Artificial (StackSpot AI)** para automatizar a geração de código, testes, documentação e validação seguindo conceitos como **Clean Architecture** e **Design Patterns**.

Através da abertura de Issues em um repositório, o sistema interpreta os pedidos dos usuários, executa múltiplos agentes de IA especializados e entrega um **Pull Request pronto para revisão** — tudo isso com execução paralela para máxima eficiência.

---

## 🎯 Objetivos do Projeto

- **Automatizar a geração de código** a partir de pedidos em linguagem natural ou estruturada
- **Padronizar entregas** seguindo arquiteturas e práticas previamente definidas
- **Reduzir o tempo de desenvolvimento** através de múltiplos agentes IA executando em paralelo
- **Permitir fácil integração** com qualquer repositório GitHub através de GitHub Actions
- **Introduzir Go** como linguagem backend na organização de forma prática e segura

---

## 🔗 Fluxo de Funcionamento

### **1. Trigger & Interpretação**
```
GitHub Issue → Issue Parser → Interpretation Agent → JSON-RPC Structure
```

### **2. Orquestração & Execução**
```
JSON-RPC → Pipeline Manager → Agent Coordinator → Parallel Execution
                                              ├─ Code Generation Agent
                                              ├─ Test Generation Agent
                                              └─ Documentation Agent
```

### **3. Consolidação & Entrega**
```
Parallel Results → Result Consolidator → GitHub PR Creation
```

---

## 🔑 Agentes Especializados

| Agente | Função | Input | Output |
|--------|--------|-------|---------|
| **Issue Interpreter** | Interpreta descrição e gera estrutura JSON-RPC | Issue Description | JSON-RPC Config |
| **Code Generator** | Produz código-fonte seguindo padrões | JSON-RPC + Context | Source Code |
| **Test Generator** | Cria testes unitários/integração | Code + Requirements | Test Files |
| **Documentation Generator** | Produz documentação técnica | Code + Context | Documentation |

---

## 📝 Exemplo de Estrutura JSON-RPC

```json
{
  "jsonrpc": "2.0",
  "method": "generateFeature",
  "params": {
    "feature": {
      "name": "Adicionar endpoint de criação de usuário",
      "description": "Endpoint REST para criação de usuários com validação"
    },
    "architecture": {
      "pattern": "Clean Architecture",
      "language": "Go",
      "framework": "Gin"
    },
    "design_patterns": [
      "Repository Pattern",
      "DTO Pattern",
      "Builder Pattern"
    ],
    "requirements": {
      "testing": true,
      "documentation": true,
      "validation": true
    },
    "agents": [
      {
        "name": "code_generator",
        "priority": 1,
        "config": {
          "template": "rest_endpoint",
          "validation": true
        }
      },
      {
        "name": "test_generator",
        "priority": 2,
        "dependencies": ["code_generator"]
      },
      {
        "name": "doc_generator",
        "priority": 3,
        "dependencies": ["code_generator"]
      }
    ]
  },
  "id": 1
}
```
---

## 🔄 Estratégia de Migração

### **Flexibilidade de Protocolos**
O sistema suporta múltiplos protocolos de comunicação com a IA:

### **Plano de Migração**
1. **Atual**: `mark3labs/mcp-go` para prototipagem
2. **Fallback**: HTTP adapter como backup
3. **Futuro**: Migração para SDK oficial MCP quando disponível
4. **Flexibilidade**: Troca de protocolo via configuração

---

## 🛡️ Segurança e Confiabilidade

### **Validação de Entrada**
- Sanitização de Issues do GitHub
- Validação de JSON-RPC estruturado
- Rate limiting para APIs

### **Tratamento de Erros**
- Circuit breaker para APIs externas
- Retry com backoff exponencial
- Fallback para HTTP em caso de falha MCP

### **Monitoramento**
- Métricas de performance
- Logs estruturados
- Health checks

---

## 🎯 Benefícios Esperados

### **Técnicos**
- **Performance**: Execução paralela de agentes
- **Escalabilidade**: Arquitetura modular e desacoplada
- **Manutenibilidade**: Código limpo e testável
- **Flexibilidade**: Múltiplos protocolos de comunicação

### **Organizacionais**
- **Introdução do Go**: Showcase da linguagem
- **Automação**: Redução de trabalho manual
- **Padronização**: Código consistente
- **Inovação**: Uso de tecnologias emergentes (MCP)

---

## 📊 Regras de Negócio

### **Execução de Agentes**
- Agentes executam em paralelo quando não há dependências
- Sistema de prioridades define ordem de execução
- Dependências entre agentes são respeitadas
- Falha de um agente não interrompe outros agentes independentes

### **Interpretação de Issues**
- Issues devem seguir templates específicos ou linguagem natural
- Sistema detecta automaticamente o tipo de solicitação
- Validação obrigatória da estrutura JSON-RPC gerada
- Feedback automático em caso de interpretação ambígua

### **Geração de Código**
- Sempre seguir padrões arquiteturais definidos
- Validação automática de código gerado
- Testes obrigatórios para todo código produzido
- Documentação automática de APIs e funcionalidades

### **Pull Request**
- PR criado automaticamente após consolidação
- Inclui descrição detalhada das mudanças
- Testes executados automaticamente
- Revisão humana obrigatória antes do merge

### **Configuração**
- Configuração por repositório através de arquivos
- Override de configurações via Issues
- Versionamento de configurações
- Rollback automático em caso de falha
