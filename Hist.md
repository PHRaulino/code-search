### Feature: Gestão de Acessos para Fundos

#### Contexto Atual

Atualmente, o sistema possui uma página de gestão de acessos onde é possível atribuir permissões de funcionalidades para os usuários. Dentro dessa página, os administradores podem selecionar as funcionalidades que cada grupo de acesso pode utilizar. As funcionalidades definem o que os usuários podem fazer dentro do sistema.

Na versão atual, quando um usuário acessa uma funcionalidade, a lista de fundos disponível para ele é determinada pelos fundos aos quais a empresa do usuário possui acesso. Isso significa que todos os fundos disponíveis para a empresa do usuário são automaticamente acessíveis a todos os usuários com permissão para essa funcionalidade.

#### Objetivo da Nova Feature

A nova feature de Gestão de Acessos para Fundos visa adicionar um nível adicional de controle sobre os acessos, permitindo que os administradores não apenas escolham as funcionalidades para cada grupo de acesso, mas também possam especificar quais fundos, classes e subclasses estão disponíveis para esses grupos. Isso permitirá uma gestão mais granular e segura dos acessos, garantindo que os usuários vejam apenas os fundos explicitamente permitidos para eles.

### Histórias e Tarefas

#### História 1: Melhoria no Micro Front-End (MFE)

**DOOR:** Ajustar o MFE e adaptá-lo para ter uma melhor implantação da página de gestão de acessos, modularizando componentes e melhorando a estrutura e a organização do front-end.

**Tarefas:**

1. **Renomear Páginas do Front-End:**
   - Renomear pastas de páginas.
   - Assegurar consistência nos nomes.

2. **Modularização do Front-End:**
   - Refatorar e modularizar componentes.
   - Criar novos componentes reutilizáveis:
     - Container de título e subtítulo.
     - Modal de tabela.
     - Modal de formulário.
     - Modal de mensagem.
     - Componente de tabela.
     - Card de tabela com busca.
     - Select e Multiselect.
     - Componente de busca.
     - Componente de card.
     - Input.
     - Botão.

3. **Ajuste de Layout:**
   - Implementar quebra de linhas.
   - Adicionar largura mínima nas tabelas.

4. **Handlers de Erro e Logs no Console:**
   - Adicionar handlers de erro e logs para diferentes situações.

5. **Gerenciamento de Estado com NgRx:**
   - Implementar NgRx para gerenciamento de estado.

6. **Adicionar ID Único para Componentes:**
   - Atribuir IDs únicos para componentes.

7. **Atualização da Versão do Voxel:**
   - Atualizar o Voxel para a versão mais recente.

8. **Adicionar Testes Unitários:**
   - Desenvolver e garantir cobertura de testes unitários.

9. **Responsividade do Componente de Tabela:**
   - Tornar a tabela responsiva.

10. **Usar Shimmer nos Componentes para Carregamento:**
    - Implementar shimmer effects nos componentes.

11. **Ajustes no Layout em Geral:**
    - Melhorar a disposição dos elementos na interface.

#### História 2: Adaptação do Back-End e Melhorias na Rota RBAC Access Manipulation

**DOOR:** Alterar o payload que será enviado para o back-end para fazer as manipulações no grupo de acesso. A principal alteração é a inclusão de classes e subclasses no JSON dentro das chaves `remove` e `add` de `policies`.

**Estrutura Atual e Nova do Payload**

**Estrutura Atual do Payload:**
```json
{
  "id": int,
  "name": "string | null",
  "description": "string | null",
  "changes": ["users", "policies"],
  "payload": {
    "users": {
      "remove": [int],
      "add": [int]
    },
    "policies": {
      "remove": [
        {
          "policy_id": int,
          "funds": [int]
        }
      ],
      "add": [
        {
          "policy_id": int,
          "funds": [int]
        }
      ]
    }
  }
}
```

**Nova Estrutura do Payload:**
```json
{
  "id": int,
  "name": "string | null",
  "description": "string | null",
  "changes": ["users", "policies"],
  "payload": {
    "users": {
      "remove": [int],
      "add": [int]
    },
    "policies": {
      "remove": [
        {
          "policy_id": int,
          "funds": [int],
          "classes": [int],
          "subclasses": [int]
        }
      ],
      "add": [
        {
          "policy_id": int,
          "funds": [int],
          "classes": [int],
          "subclasses": [int]
        }
      ]
    }
  }
}
```

**Tarefas:**

1. **Atualizar a Rota RBAC Access Manipulation:**
   - Receber e processar o novo payload JSON.
   - Manipular as chaves `id`, `name`, `description` e `changes`.

2. **Manipulação de Usuários:**
   - Remover usuários (`users.remove`).
   - Adicionar usuários (`users.add`).

3. **Manipulação de Políticas:**
   - Remover políticas (`policies.remove`).
   - Adicionar políticas (`policies.add`).

4. **Validação de Dados:**
   - Validar dados do payload.
   - Garantir a existência e validade dos IDs.

5. **Log de Erros e Monitoramento:**
   - Adicionar logs detalhados.
   - Integrar monitoramento para performance e problemas.

6. **Atualização de Documentação:**
   - Atualizar documentação da API.
   - Incluir exemplos de requests e responses.

7. **Testes Unitários e de Integração:**
   - Desenvolver testes unitários.
   - Implementar testes de integração.

#### História 3: Implementação da Tela de Gestão de Acessos para Fundos

**DOOR:** Adaptar a tela de gestão de acessos existente para incorporar um sistema de abas que organiza os diferentes recursos (usuários, funcionalidades, classes) em uma interface mais intuitiva e fácil de usar.

**Desenho no Fingerman:**
[Link do Fingerman](#)

**Tarefas:**

1. **Criação do Card de Tabs:**
   - Implementar o card de tabs.
   - Criar abas para usuários, funcionalidades e classes.

2. **Migrar Tabelas Existentes para as Abas Correspondentes:**
   - Mover tabela de funcionalidades para a aba de funcionalidades.
   - Mover tabela de usuários para a aba de usuários.
   - Criar aba de classes com a tabela correspondente.

3. **Formulários para Adicionar e Remover Recursos:**
   - Implementar formulários em cada aba.
   - Garantir validação adequada dos campos.

4. **Responsividade e Layout:**
   - Garantir responsividade da interface.
   - Ajustar o layout para consistência visual.

5. **Integração com o Back-End:**
   - Integrar operações de adicionar e remover recursos com o back-end.
   - Atualizar chamadas de API.

6. **Teste e Validação:**
   - Testar a interface.
   - Realizar testes de usabilidade.

7. **Documentação:**
   - Atualizar documentação do projeto.
   - Incluir exemplos e instruções.

#### História 4: Desenvolvimento das Rotas de Consumo do RBAC

**DOOR:** Criar novos endpoints na API de autorização que receberão novos parâmetros e retornarão apenas os acessos concedidos especificamente para o usuário, de acordo com a nova lógica de fundos, classes e subclasses.

**Estrutura Atual e Nova das Rotas**

**Novo Endpoint:**
```
GET /api/rbac/access
```

**Parâmetros:**
- `usuario`: ID do usuário (int)
- `política`: ID da política (int)
- `tipo_vinculo`: Tipo de vínculo do usuário (string)

**Request Exemplo:**
```http
GET /api/rbac/access?usuario=123&política=456&tipo_vinculo=gerente
```

**Response Exemplo:**
```json
{
  "fundos": [1, 2, 3],
  "classes": [10, 20, 30],
  "subclasses": [100, 200, 300]
}
```

**Tarefas:**

1. **Criação de Novos Endpoints:**
   - Desenvolver endpoints para `usuário`, `política` e `tipo vínculo`.

2. **Integração com o Banco de Dados:**
   - Modificar consultas ao banco de dados.
   - Atualizar lógica de acesso com base nas políticas.

3. **Validação de Parâmetros:**
   - Implementar validações para os parâmetros.

4. **Teste e Validação dos Endpoints:**
   - Desenvolver testes unitários e de integração.
   - Realizar testes de carga.

5. **Atualização da Documentação:**
   - Atualizar a documentação da API.
   - Incluir exemplos de requests e responses.

6. **Monitoramento e Logs:**
   - Adicionar logs detalhados.
   - Integrar monitoramento de performance.

### Conclusão

A implementação da feature de Gestão de Acessos para Fundos proporcionará um controle mais granular e seguro dos acessos, aumentando a segurança e a personalização dos acessos para cada usuário. As histórias e tarefas detalhadas garantirão
