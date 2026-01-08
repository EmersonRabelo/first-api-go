# First API Go - Documentação Completa

Uma API RESTful construída em Go com moderação de conteúdo inteligente, sistema de fila de mensagens e cache distribuído.

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura](#arquitetura)
3. [Tecnologias](#tecnologias)
4. [Configuração e Instalação](#configuração-e-instalação)
5. [Estrutura do Projeto](#estrutura-do-projeto)
6. [API Endpoints](#api-endpoints)
7. [Entidades do Negócio](#entidades-do-negócio)
8. [Fluxo de Moderação](#fluxo-de-moderação)
9. [Banco de Dados](#banco-de-dados)
10. [Sistema de Fila](#sistema-de-fila)
11. [Cache e Redis](#cache-e-redis)
12. [Guia de Desenvolvimento](#guia-de-desenvolvimento)

---

## 🎯 Visão Geral

**First API Go** é uma plataforma social que permite usuários criar posts, commentários (replies) e interagir através de curtidas, com um sistema robusto de moderação de conteúdo baseado em análise automática e denúncias.

### Principais Funcionalidades

- ✅ Gerenciamento de usuários
- ✅ Criação, edição e exclusão de posts
- ✅ Sistema de curtidas com contagem em tempo real
- ✅ Respostas (replies) com contagem
- ✅ Sistema de denúncias com moderação automática
- ✅ Análise de conteúdo tóxico (Perspective API)
- ✅ Estados de processamento de posts (visível, limitado, oculto, removido)
- ✅ Cache distribuído com Redis
- ✅ Fila de mensagens com RabbitMQ

---

## 🏛️ Arquitetura

A aplicação segue uma arquitetura em camadas com separação clara de responsabilidades:

```
┌─────────────────────────────────────────────────────┐
│                   HTTP API (Gin)                    │
│              /api/v1/{resource}                     │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│              Controllers/Handlers                   │
│  (Request parsing, validation, response formatting)│
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                  Services                           │
│    (Business logic, domain rules, orchestration)   │
└────────────────────┬────────────────────────────────┘
                     │
         ┌───────────┴───────────┬──────────────┐
         │                       │              │
┌────────▼────────┐  ┌───────────▼────┐   ┌────▼─────────┐
│  Repositories   │  │  Queue Service │   │  Redis Cache │
│  (Data Access)  │  │  (Messaging)   │   │  (Counters)  │
└────────┬────────┘  └───────────┬────┘   └────┬─────────┘
         │                       │              │
         └───────────┬───────────┴──────────────┘
                     │
         ┌───────────▼──────────────┐
         │     PostgreSQL (GORM)    │
         │     + RabbitMQ + Redis   │
         └──────────────────────────┘
```

### Padrões de Design Utilizados

- **Repository Pattern**: Abstração de acesso a dados
- **Service Layer**: Lógica de negócio centralizada
- **Dependency Injection**: Injeção de dependências no main.go
- **DTO Pattern**: Data Transfer Objects para validação
- **Producer-Consumer Pattern**: Processamento assíncrono com RabbitMQ

---

## 🛠️ Tecnologias

### Backend Framework
- **Gin**: Framework web rápido e minimalista
- **GORM**: ORM para Go com suporte a PostgreSQL

### Banco de Dados
- **PostgreSQL**: Banco relacional principal
- **golang-migrate**: Versionamento de schema

### Cache e Message Queue
- **Redis**: Cache distribuído e contadores
- **RabbitMQ**: Sistema de fila para moderação

### Validação e Utilities
- **Validator/v10**: Validação de estruturas
- **UUID**: Identificadores universalmente únicos
- **godotenv**: Gerenciamento de variáveis de ambiente

### Versão Go
- **Go 1.25.4**

---

## 🚀 Configuração e Instalação

### Pré-requisitos

```bash
# Verificar versão do Go
go version  # Deve ser >= 1.25.4

# Ter Docker instalado (para serviços auxiliares)
docker --version

# Ter PostgreSQL, Redis e RabbitMQ disponíveis
```

### Setup Local

1. **Clonar repositório**
```bash
git clone https://github.com/EmersonRabelo/first-api-go.git
cd first-api-go
```

2. **Instalar dependências**
```bash
go mod download
go mod tidy
```

3. **Configurar variáveis de ambiente**
```bash
cp .env.example .env
# Editar .env com suas credenciais
```

**Variáveis obrigatórias:**
```env
# Server
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=first_api_go
DB_SSL_MODE=disable

# Message Broker (RabbitMQ)
BROKER_HOST=localhost
BROKER_PORT=5672
BROKER_USER=guest
BROKER_PASSWORD=guest

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Environment
ENVIRONMENT=development  # development|staging|production
```

4. **Inicializar banco de dados**
```bash
# PostgreSQL
createdb first_api_go

# Executar migrations (automático ao iniciar aplicação)
go run cmd/api/main.go
```

5. **Iniciar serviços com Docker** (opcional)
```bash
# PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=first_api_go \
  -p 5432:5432 \
  postgres:latest

# Redis
docker run -d --name redis \
  -p 6379:6379 \
  redis:latest

# RabbitMQ
docker run -d --name rabbitmq \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:management
```

6. **Executar a aplicação**
```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

---

## 📁 Estrutura do Projeto

### Organização de Diretórios

```
firstApiGo/
├── cmd/
│   └── api/
│       └── main.go                    # Ponto de entrada da aplicação
├── db/
│   └── migrations/                    # Migrações de schema (golang-migrate)
│       ├── 000001_create_users_table.up.sql
│       ├── 000002_create_posts_table.up.sql
│       ├── 000003_create_likes_table.up.sql
│       ├── 000004_create_replies_table.up.sql
│       ├── 000005_add_partial_unique_constraint_to_likes_table.up.sql
│       ├── 000006_create_post_likes_count_table.up.sql
│       ├── 000007_create_increment_function_to_post_like_count.up.sql
│       ├── 000008_create_trigger_to_call_increment_function.up.sql
│       ├── 000009_create_decrement_function_to_post_like_count.up.sql
│       ├── 000010_create_trigger_to_call_decrement_function.up.sql
│       ├── 000011_alter_likes_remove_quantity_column.up.sql
│       ├── 000012_create_post_replies_count_table.up.sql
│       ├── 000013_create_increment_function_to_post_reply_count.up.sql
│       ├── 000014_create_trigger_to_call_replies_count_increment_function.up.sql
│       ├── 000015_create_decrement_function_to_post_replies_count.up.sql
│       ├── 000016_create_trigger_to_call_replies_count_decrement_function.up.sql
│       ├── 000017_alter_replies_remove_quantity_column.up.sql
│       ├── 000018_create_reports_table_and_indexes.up.sql
│       ├── 000019_alter_reports_table_perspective_identity_hate_column_name.up.sql
│       ├── 000020_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql
│       └── 000021_alter_table_posts_add_new_flag_column.up.sql
├── docs/
│   └── POST_MODERATION_RULES.md       # Regras de moderação de posts
├── internal/
│   ├── config/                        # Configuração e inicialização
│   │   ├── broker.go                  # Configuração RabbitMQ
│   │   ├── config.go                  # Variáveis de ambiente
│   │   └── database.go                # Configuração PostgreSQL e GORM
│   ├── controller/                    # Handlers HTTP (Gin)
│   │   ├── like_controller.go
│   │   ├── post_controller.go
│   │   ├── reply_controller.go
│   │   ├── report_controller.go
│   │   └── user_controller.go
│   ├── database/                      # Gerenciamento de migrations
│   │   └── migration.go
│   ├── dtos/                          # Data Transfer Objects
│   │   ├── like/
│   │   ├── post/
│   │   ├── reply/
│   │   ├── report/
│   │   ├── shared/
│   │   └── user/
│   ├── entity/                        # Modelos de domínio (GORM)
│   │   ├── like.go
│   │   ├── post.go
│   │   ├── post_like_count.go
│   │   ├── reply.go
│   │   ├── report.go
│   │   └── user.go
│   ├── handler/                       # Handlers de fila de mensagens
│   │   └── report_handler.go
│   ├── queue/                         # Produtor/Consumidor RabbitMQ
│   │   ├── consumer.go
│   │   └── producer.go
│   ├── redis/                         # Cliente e serviços Redis
│   │   ├── client.go
│   │   └── counter.go
│   ├── repository/                    # Data Access Objects
│   │   ├── like_repository.go
│   │   ├── post_repository.go
│   │   ├── reply_repository.go
│   │   ├── report_repository.go
│   │   └── user_repository.go
│   └── service/                       # Lógica de negócio
│       ├── consumer/                  # Consumidor de relatórios
│       │   └── report_consumer.go
│       ├── report/                    # Serviço de moderação
│       │   └── report_service.go
│       ├── like_service.go
│       ├── post_service.go
│       ├── reply_service.go
│       └── user_service.go
├── postman/                           # Coleções Postman para testes
│   ├── collections/
│   ├── environments/
│   └── globals/
├── router/                            # Configuração de rotas (Gin)
│   └── router.go
├── go.mod                             # Definição de módulo Go
├── go.sum                             # Checksums de dependências
└── README.md                          # Este arquivo
```

### Responsabilidades por Camada

#### `cmd/api/main.go` - Inicialização
- Carrega configurações
- Instancia conexões (DB, Redis, RabbitMQ)
- Cria injeção de dependências
- Inicia servidor HTTP e consumidores assíncronos

#### `config/` - Configuração
- **config.go**: Leitura de variáveis de ambiente
- **database.go**: Conexão PostgreSQL via GORM
- **broker.go**: Conexão RabbitMQ

#### `controller/` - API HTTP
Handlers que recebem requisições HTTP e chamam services:
- Parsing de JSON
- Validação de entrada
- Tratamento de erros
- Formatação de respostas

#### `service/` - Lógica de Negócio
Orquestra operações de negócio:
- Validações de regras de negócio
- Chamadas a múltiplos repositórios
- Chamadas a fila de mensagens
- Chamadas a cache

#### `repository/` - Acesso a Dados
Interface com o banco de dados:
- Queries ao PostgreSQL via GORM
- Transações quando necessário
- Tratamento de erros de banco

#### `queue/` - Processamento Assíncrono
Sistema produtor-consumidor:
- **producer.go**: Publica mensagens em RabbitMQ
- **consumer.go**: Consome mensagens de forma assíncrona

#### `redis/` - Cache Distribuído
Cache em memória para performance:
- **client.go**: Conexão e operações básicas
- **counter.go**: Incrementos/decrementos de contadores

#### `entity/` - Modelos de Domínio
Estruturas Go com tags GORM, representam tabelas do banco

#### `dtos/` - Validação de Entrada
Estruturas Go com tags `validate` para validar dados de entrada das requisições

---

## 📡 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```
GET /health
Resposta: { "status": "running", "time": "2026-01-08T10:30:00Z" }
```

### Users (Usuários)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/users` | Lista todos os usuários |
| GET | `/users/:id` | Obtém um usuário por ID |
| POST | `/users` | Cria um novo usuário |
| PUT | `/users/:id` | Atualiza um usuário |
| DELETE | `/users/:id` | Deleta um usuário |

**Exemplo de criação:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "João Silva", "email": "joao@example.com"}'
```

### Posts (Publicações)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/posts` | Lista todos os posts (públicos) |
| GET | `/posts/:id` | Obtém um post específico |
| POST | `/posts` | Cria um novo post |
| PUT | `/posts/:id` | Atualiza um post |
| DELETE | `/posts/:id` | Deleta um post |
| POST | `/posts/:id/report` | Denuncia um post |

**Exemplo de criação:**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "body": "Este é meu primeiro post!"
  }'
```

**Estados do Post (flag):**
- `visible`: Post está visível normalmente
- `limited`: Post visível com restrições
- `hidden_pending_review`: Oculto aguardando revisão humana
- `removed`: Permanentemente removido

### Likes (Curtidas)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/likes` | Lista todas as curtidas |
| GET | `/likes/:id` | Obtém uma curtida específica |
| POST | `/likes` | Cria uma curtida em um post |
| DELETE | `/likes/:id` | Remove uma curtida |

**Exemplo:**
```bash
curl -X POST http://localhost:8080/api/v1/likes \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000"
  }'
```

### Replies (Respostas/Comentários)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/replies` | Lista todos os comentários |
| GET | `/replies/:id` | Obtém um comentário específico |
| POST | `/replies` | Cria um comentário em um post |
| PUT | `/replies/:id` | Atualiza um comentário |
| DELETE | `/replies/:id` | Deleta um comentário |

**Exemplo:**
```bash
curl -X POST http://localhost:8080/api/v1/replies \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "body": "Ótimo post!"
  }'
```

### Reports (Denúncias/Moderação)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/posts/:id/report` | Denuncia um post para moderação |
| GET | `/reports` | Lista relatórios (admin) |
| GET | `/reports/:id` | Obtém detalhes de um relatório |

**Exemplo de denúncia:**
```bash
curl -X POST http://localhost:8080/api/v1/posts/660e8400-e29b-41d4-a716-446655440000/report \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "reason": "Conteúdo ofensivo"
  }'
```

---

## 👥 Entidades do Negócio

### User (Usuário)

```go
type User struct {
    Id        uuid.UUID     // Identificador único
    Name      string        // Nome completo
    Email     string        // Email único
    CreatedAt time.Time     // Data de criação
    UpdatedAt time.Time     // Data da última atualização
    DeletedAt gorm.DeletedAt // Soft delete
}
```

**Operações:**
- Criar usuário
- Listar usuários
- Buscar por ID
- Atualizar perfil
- Deletar conta (soft delete)

---

### Post (Publicação)

```go
type Post struct {
    Id        uuid.UUID       // Identificador único
    UserId    uuid.UUID       // ID do criador
    User      User            // Relacionamento com usuário
    Body      string          // Conteúdo (max 280 caracteres)
    Flag      ProcessFlag     // Estado de moderação
    IsActive  bool            // Se está ativo
    CreatedAt time.Time       // Data de criação
    UpdatedAt time.Time       // Data da última atualização
    DeletedAt gorm.DeletedAt  // Soft delete
}
```

**Estados Possíveis (ProcessFlag):**

| Estado | Descrição |
|--------|-----------|
| `visible` | Post visível normalmente a todos |
| `limited` | Post visível com restrições (reduz alcance, remove de recomendações) |
| `hidden_pending_review` | Oculto até revisão humana |
| `removed` | Removido permanentemente |

**Operações:**
- Criar post
- Listar posts
- Buscar por ID
- Atualizar conteúdo
- Deletar post
- Denuncia (abre fluxo de moderação)

---

### Like (Curtida)

```go
type Like struct {
    Id        uuid.UUID      // Identificador único
    UserId    uuid.UUID      // ID de quem curtiu
    PostId    uuid.UUID      // ID do post curtido
    CreatedAt time.Time      // Data da curtida
    DeletedAt gorm.DeletedAt // Soft delete
}
```

**Constraints:**
- Um usuário não pode curtir o mesmo post duas vezes
- Não é possível curtir seu próprio post (verificado no serviço)

**Operações:**
- Criar curtida
- Remover curtida
- Listar curtidas
- Contagem é mantida em Redis (performance)

---

### Reply (Resposta/Comentário)

```go
type Reply struct {
    Id        uuid.UUID      // Identificador único
    UserId    uuid.UUID      // ID de quem respondeu
    PostId    uuid.UUID      // ID do post sendo comentado
    Body      string         // Conteúdo do comentário
    IsActive  bool           // Se está ativo
    CreatedAt time.Time      // Data da resposta
    UpdatedAt time.Time      // Data da última atualização
    DeletedAt gorm.DeletedAt // Soft delete
}
```

**Operações:**
- Criar comentário
- Listar comentários por post
- Buscar comentário por ID
- Atualizar comentário
- Deletar comentário
- Contagem é mantida em Redis (performance)

---

### Report (Denúncia)

```go
type Report struct {
    Id                    uuid.UUID     // Identificador único
    UserId                uuid.UUID     // ID de quem denunciou
    PostId                uuid.UUID     // ID do post denunciado
    Reason                string        // Motivo da denúncia
    PerspectiveToxicity   float64       // Score Perspective API
    PerspectiveInsult     float64       // Score Perspective API
    PerspectiveProfanity  float64       // Score Perspective API
    PerspectiveThreat     float64       // Score Perspective API
    PerspectiveIdentityAttack float64   // Score Perspective API
    PerspectiveSevereToxicity float64   // Score Perspective API
    Status                string        // Estado do relatório
    CreatedAt             time.Time     // Data da denúncia
    UpdatedAt             time.Time     // Data da última atualização
}
```

**Fluxo:**
1. Usuário denuncia post via API
2. Serviço de relatório envia mensagem para fila (RabbitMQ)
3. Consumidor processa denúncia (análise Perspective API)
4. Resultado atualiza status do post

---

## 🚨 Fluxo de Moderação

### Visão Geral

A moderação é pós-publicação: posts são publicados imediatamente e analisados apenas quando denunciados.

### Fluxo Detalhado

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. Usuário denuncia um post via POST /posts/:id/report          │
└────────────────┬────────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────────┐
│ 2. ReportController recebe denúncia                             │
│    - Valida dados (user_id, post_id, reason)                   │
│    - Chama ReportService.Create()                              │
└────────────────┬────────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────────┐
│ 3. ReportService orquestra operação                             │
│    - Valida se post existe e usuário tem permissão             │
│    - Cria documento Report no banco                            │
│    - Publica mensagem em RabbitMQ (topic_report, routing key)  │
└────────────────┬────────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────────┐
│ 4. RabbitMQ armazena mensagem                                  │
│    Exchange: topic_report                                      │
│    Routing Key: post.report.created                            │
└────────────────┬────────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────────┐
│ 5. Consumidor RabbitMQ processa mensagem (background)          │
│    - Consome fila q.report.response                            │
│    - Análise com Perspective API                               │
│    - Calcula scores de toxicidade                              │
└────────────────┬────────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────────┐
│ 6. ConsumerReportService aplica regras de moderação            │
│    Baseado em scores:                                          │
│    - THREAT > 0.7 ou IDENTITY_ATTACK > 0.8                    │
│      → hidden_pending_review                                  │
│    - SEVERE_TOXICITY > 0.9                                    │
│      → hidden_pending_review                                  │
│    - TOXICITY > 0.85 ou INSULT > 0.8                          │
│      → limited (restrições)                                   │
│    - Caso contrário → visible                                 │
└────────────────┬────────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────────┐
│ 7. Atualiza Post com novo flag                                 │
│    - Persiste no PostgreSQL                                    │
│    - Post fica visível/limitado/oculto conforme decisão      │
└─────────────────────────────────────────────────────────────────┘
```

### Regras de Decisão

Consultar [docs/POST_MODERATION_RULES.md](docs/POST_MODERATION_RULES.md) para detalhes completos.

**Resumo:**

| Condição | Ação | Justificativa |
|----------|------|---------------|
| THREAT > 0.7 | Hidden | Ameaça imediata |
| IDENTITY_ATTACK > 0.8 | Hidden | Ataque a identidade |
| SEVERE_TOXICITY > 0.9 | Hidden | Toxicidade severa |
| TOXICITY > 0.85 OU INSULT > 0.8 | Limited | Conteúdo tóxico |
| Demais casos | Visible | Baixo risco |

### Estados de Processamento

```
                    ┌──────────────┐
                    │   VISIBLE    │
                    │ (Padrão)     │
                    └──────┬───────┘
                           │ (denúncia)
                    ┌──────▼───────┐
                    │ Análise API  │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼────┐ ┌─────▼──────┐ ┌──▼──────────────┐
       │  LIMITED  │ │   VISIBLE  │ │ HIDDEN_PENDING  │
       │ (restr.)  │ │(mantém)    │ │   _REVIEW       │
       └───────────┘ └────────────┘ └────────────────┘
              │            │                 │
              │            │          (revisão manual)
              │            │                 │
              └────────────┼─────────────────┘
                           │
                    ┌──────▼────────┐
                    │   REMOVED     │
                    │(permanente)   │
                    └───────────────┘
```

---

## 🗄️ Banco de Dados

### Visão Geral

PostgreSQL com GORM ORM. Schema versionado com golang-migrate.

### Diagrama E-R Simplificado

```
┌─────────────┐         ┌──────────────┐
│   USERS     │         │    POSTS     │
├─────────────┤         ├──────────────┤
│ id (PK)     │◄────────│ id (PK)      │
│ name        │ user_id │ user_id (FK) │
│ email       │         │ body         │
│ created_at  │         │ flag         │
│ updated_at  │         │ is_active    │
│ deleted_at  │         │ created_at   │
└─────────────┘         │ updated_at   │
       ▲                │ deleted_at   │
       │                └──────┬───────┘
       │                       │
       │         ┌─────────────┼──────────────┐
       │         │             │              │
   ┌───┴──┐  ┌──▼────┐   ┌────▼────┐   ┌────▼──────┐
   │LIKES │  │REPLIES│   │ REPORTS │   │POST_LIKES  │
   │      │  │       │   │         │   │_COUNT      │
   │ id   │  │ id    │   │ id      │   │            │
   │user_ │  │user_  │   │ user_id │   │post_id(FK) │
   │id(FK)│  │id(FK) │   │post_id  │   │count       │
   │post_ │  │post_  │   │reason   │   │updated_at  │
   │id(FK)│  │id(FK) │   │analysis │   │            │
   │      │  │body   │   │scores   │   └────────────┘
   └──────┘  │       │   └─────────┘
             └───────┘   ┌──────────────────┐
                         │POST_REPLIES_COUNT│
                         │                  │
                         │post_id (FK)      │
                         │count             │
                         │updated_at        │
                         └──────────────────┘
```

### Principais Tabelas

#### `users`
```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### `posts`
```sql
CREATE TABLE posts (
    post_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    post_body VARCHAR(280),
    flag VARCHAR(48) DEFAULT 'visible',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### `likes`
```sql
CREATE TABLE likes (
    like_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    post_id UUID NOT NULL REFERENCES posts(post_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id, post_id) -- Um usuário, uma curtida por post
);
```

#### `replies`
```sql
CREATE TABLE replies (
    reply_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    post_id UUID NOT NULL REFERENCES posts(post_id),
    reply_body TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### `reports`
```sql
CREATE TABLE reports (
    report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    post_id UUID NOT NULL REFERENCES posts(post_id),
    reason TEXT,
    perspective_toxicity FLOAT8,
    perspective_insult FLOAT8,
    perspective_profanity FLOAT8,
    perspective_threat FLOAT8,
    perspective_identity_attack FLOAT8,
    perspective_severe_toxicity FLOAT8,
    status VARCHAR(48),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Índices e Performance

- **Primary Keys**: Índices automáticos em todas as PK
- **Foreign Keys**: Índices em relacionamentos
- **Unique Constraints**: `users(email)`, `likes(user_id, post_id)`
- **Soft Deletes**: Índices em `deleted_at` para queries eficientes

### Migrações

Executadas automaticamente na inicialização via `database.RunMigrations()`.

**Convenção de Nomes:**
```
000NNN_descriptive_name.up.sql
```

Exemplo:
- `000001_create_users_table.up.sql`
- `000018_create_reports_table_and_indexes.up.sql`

---

## 📨 Sistema de Fila

### Arquitetura

RabbitMQ com padrão **Topic Exchange** para roteamento flexível.

```
┌──────────────────────────────────────────────┐
│          ReportController                    │
│  (recebe denúncia HTTP)                      │
└────────────┬─────────────────────────────────┘
             │ POST /reports
             │
┌────────────▼─────────────────────────────────┐
│          ReportService                       │
│  (orquestra, publica msg)                    │
└────────────┬─────────────────────────────────┘
             │ reportProducer.Publish()
             │
┌────────────▼─────────────────────────────────┐
│     RabbitMQ Topic Exchange                  │
│  Exchange: topic_report                      │
│  Routing Keys:                               │
│    - post.report.created (producer)          │
│    - post.report.response (consumer)         │
└────────────┬─────────────────────────────────┘
             │
┌────────────▼─────────────────────────────────┐
│        RabbitMQ Queues                       │
│  q.report.response                           │
│  (armazena mensagens)                        │
└────────────┬─────────────────────────────────┘
             │
┌────────────▼─────────────────────────────────┐
│      ReportConsumer (goroutine)              │
│  (processa assincronamente)                  │
└────────────┬─────────────────────────────────┘
             │
┌────────────▼─────────────────────────────────┐
│   ConsumerReportService                      │
│  (aplica regras, atualiza BD)                │
└──────────────────────────────────────────────┘
```

### Configuração

**Arquivo:** [internal/config/broker.go](internal/config/broker.go)

```go
conn, channel := config.InitBroker()

exchange := "topic_report"
routingKeyProducer := "post.report.created"
routingKeyConsumer := "post.report.response"
queueName := "q.report.response"
```

### Producer

**Arquivo:** [internal/queue/producer.go](internal/queue/producer.go)

Publica denúncias quando criadas:
```go
reportProducer := queue.NewReportProducer(channel, exchange, routingKeyProducer)
reportService := reportService.NewReportService(reportRepository, reportProducer, ...)

// Internamente, ao criar report:
reportProducer.Publish(reportMessage) // Publica em RabbitMQ
```

### Consumer

**Arquivo:** [internal/queue/consumer.go](internal/queue/consumer.go)

Consome mensagens em goroutine separada:
```go
reportConsumer := queue.NewReportConsumer(
    channel,
    exchange,
    routingKeyConsumer,
    queueName,
    handler,
)

go reportConsumer.Start() // Executa em background
```

### Handler

**Arquivo:** [internal/handler/report_handler.go](internal/handler/report_handler.go)

Processa cada mensagem:
- Análise com Perspective API
- Calcula scores de toxicidade
- Atualiza estado do post

---

## 💾 Cache e Redis

### Propósito

Cache distribuído para contadores (likes, replies) evita queries frequentes.

### Arquitetura

```
┌─────────────────────────────────┐
│   Application Layer             │
│  (like_service, reply_service)  │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│  redis.NewClient()              │
│  redis.IncrementLikeCount()     │
│  redis.DecrementLikeCount()     │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│       Redis Server              │
│  Keys:                          │
│  - post:{postId}:likes          │
│  - post:{postId}:replies        │
│  (contadores em tempo real)     │
└─────────────────────────────────┘
```

### Cliente Redis

**Arquivo:** [internal/redis/client.go](internal/redis/client.go)

```go
client := redis.NewClient()
```

Configuração padrão:
- Host: `localhost` (via variável `REDIS_HOST`)
- Port: `6379` (via variável `REDIS_PORT`)

### Operações de Contador

**Arquivo:** [internal/redis/counter.go](internal/redis/counter.go)

```go
// Incrementar
client.IncrementLikeCount(postId)        // post:{postId}:likes +1
client.IncrementReplyCount(postId)       // post:{postId}:replies +1

// Decrementar
client.DecrementLikeCount(postId)        // post:{postId}:likes -1
client.DecrementReplyCount(postId)       // post:{postId}:replies -1

// Obter
count := client.GetLikeCount(postId)     // Retorna int
```

### Sincronização com Banco

O PostgreSQL mantém verdade:
- Tabelas `post_likes_count` e `post_replies_count`
- Triggers automáticos sincronizam com Redis
- Em caso de desincronização, Redis é reconstruído do banco

### Triggers PostgreSQL

**Likes:**
```sql
-- 000006: CREATE TABLE post_likes_count
-- 000007: CREATE FUNCTION increment_post_like_count()
-- 000008: CREATE TRIGGER trg_increment_like_count
-- 000009: CREATE FUNCTION decrement_post_like_count()
-- 000010: CREATE TRIGGER trg_decrement_like_count
```

**Replies:**
```sql
-- 000012: CREATE TABLE post_replies_count
-- 000013: CREATE FUNCTION increment_post_reply_count()
-- 000014: CREATE TRIGGER trg_increment_reply_count
-- 000015: CREATE FUNCTION decrement_post_reply_count()
-- 000016: CREATE TRIGGER trg_decrement_reply_count
```

---

## 🔨 Guia de Desenvolvimento

### Ambiente de Desenvolvimento

**Requisitos:**
- Go 1.25.4+
- PostgreSQL 12+
- Redis 6+
- RabbitMQ 3.8+

**Ferramentas Recomendadas:**
```bash
# Editor/IDE
VS Code + Go extension
GoLand / IntelliJ IDEA

# CLI tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest  # Hot reload

# Database tools
psql (PostgreSQL CLI)
redis-cli (Redis CLI)
```

### Estrutura de um Novo Endpoint

Para adicionar um novo endpoint, siga este padrão:

#### 1. Criar DTO (se necessário)
```go
// internal/dtos/feature/create_request.go
package feature

type CreateRequest struct {
    Name string `json:"name" validate:"required"`
}
```

#### 2. Criar Entity (modelo)
```go
// internal/entity/feature.go
package entity

type Feature struct {
    Id        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
    Name      string         `json:"name"`
    CreatedAt time.Time      `json:"created_at"`
}
```

#### 3. Criar Repository
```go
// internal/repository/feature_repository.go
package repository

type FeatureRepository interface {
    Create(ctx context.Context, feature *entity.Feature) error
    FindById(ctx context.Context, id uuid.UUID) (*entity.Feature, error)
}

type featureRepository struct {
    db *gorm.DB
}

func (r *featureRepository) Create(ctx context.Context, f *entity.Feature) error {
    return r.db.WithContext(ctx).Create(f).Error
}
```

#### 4. Criar Service
```go
// internal/service/feature_service.go
package service

type FeatureService struct {
    repository repository.FeatureRepository
}

func (s *FeatureService) CreateFeature(ctx context.Context, req *dtos.CreateRequest) (*entity.Feature, error) {
    feature := &entity.Feature{
        Id:   uuid.New(),
        Name: req.Name,
    }
    
    if err := s.repository.Create(ctx, feature); err != nil {
        return nil, err
    }
    
    return feature, nil
}
```

#### 5. Criar Controller
```go
// internal/controller/feature_controller.go
package controller

type FeatureHandler struct {
    service *service.FeatureService
}

func (h *FeatureHandler) Create(c *gin.Context) {
    var req dtos.CreateRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    feature, err := h.service.CreateFeature(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, feature)
}
```

#### 6. Registrar Routes
```go
// router/router.go
func SetupRouter(...) *gin.Engine {
    r := gin.Default()
    v1 := r.Group("/api/v1")
    
    features := v1.Group("/features")
    {
        features.POST("", featureHandler.Create)
        features.GET("/:id", featureHandler.FindById)
    }
    
    return r
}
```

#### 7. Atualizar main.go
```go
// cmd/api/main.go
func main() {
    // ... existing code ...
    
    featureRepository := repository.NewFeatureRepository(db)
    featureService := service.NewFeatureService(featureRepository)
    featureHandler := controller.NewFeatureHandler(featureService)
    
    r := router.SetupRouter(..., featureHandler)
}
```

### Boas Práticas

#### 1. **Separação de Responsabilidades**
- Controllers: apenas HTTP
- Services: lógica de negócio
- Repositories: acesso a dados
- DTOs: validação de entrada

#### 2. **Tratamento de Erros**
```go
// ❌ Ruim
if err != nil {
    panic(err)
}

// ✅ Bom
if err != nil {
    log.Error("operação falhou", err)
    return nil, fmt.Errorf("falha ao criar feature: %w", err)
}
```

#### 3. **Validação de Entrada**
```go
// ❌ Ruim
name := c.PostForm("name")
// Usa direto sem validar

// ✅ Bom
var req dtos.CreateRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

#### 4. **Logging**
```go
import "log"

// Use log.Println, log.Printf para debug
log.Printf("Criando feature: %v", feature)
```

#### 5. **Context em Operações Assíncronas**
```go
// ✅ Bom - respeita timeout/cancellação
func (s *FeatureService) CreateFeature(ctx context.Context, req *dtos.CreateRequest) error {
    return s.repository.Create(ctx, feature)
}

// ❌ Ruim - ignora contexto
func (s *FeatureService) CreateFeature(req *dtos.CreateRequest) error {
    return s.repository.Create(context.Background(), feature)
}
```

### Testes

#### Testes Unitários
```go
// feature_service_test.go
package service

import "testing"

func TestCreateFeature(t *testing.T) {
    // Arrange
    mockRepo := NewMockRepository()
    service := NewFeatureService(mockRepo)
    
    // Act
    feature, err := service.CreateFeature(context.Background(), &dtos.CreateRequest{Name: "Test"})
    
    // Assert
    if err != nil {
        t.Fatalf("erro inesperado: %v", err)
    }
    if feature.Name != "Test" {
        t.Errorf("esperado Test, recebido %s", feature.Name)
    }
}
```

Executar testes:
```bash
go test ./...
go test -v ./internal/service/...
go test -cover ./...  # Com cobertura
```

### Debugging

#### Com VS Code
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to Delve",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "dlvToolPath": "${workspaceFolder}/../dlv",
            "port": 38697
        }
    ]
}
```

Executar com debugger:
```bash
dlv debug ./cmd/api
```

#### Logs
```go
import "log"

log.Println("Iniciando...")
log.Printf("User ID: %s", userId)
```

### Migrations

Criar nova migration:
```bash
migrate create -ext sql -dir db/migrations -seq create_features_table
```

Editar arquivo `.up.sql`:
```sql
CREATE TABLE features (
    feature_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Aplicar migrations:
```bash
migrate -path db/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME" up
```

### Deploy

#### Variáveis de Ambiente (Produção)

```env
ENVIRONMENT=production

# Database
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=secure_password
DB_NAME=first_api_go_prod
DB_SSL_MODE=require

# Server
SERVER_PORT=8080

# Broker
BROKER_HOST=prod-rabbitmq.example.com
BROKER_PORT=5672
BROKER_USER=prod_broker_user
BROKER_PASSWORD=secure_broker_password

# Redis
REDIS_HOST=prod-redis.example.com
REDIS_PORT=6379
```

#### Docker

```dockerfile
# Dockerfile
FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
```

Construir e executar:
```bash
docker build -t first-api-go:latest .
docker run -p 8080:8080 --env-file .env first-api-go:latest
```

### CI/CD (GitHub Actions - Exemplo)

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.25.4
      
      - name: Run tests
        run: go test ./...
      
      - name: Build
        run: go build -o api cmd/api/main.go
      
      - name: Deploy
        run: ./deploy.sh
        env:
          DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
```

---

## 📊 Monitoramento e Logs

### Estrutura de Logs

```
[TIMESTAMP] [LEVEL] [MODULE] message
```

Exemplo:
```
2026-01-08 10:30:45 INFO config Aplicação inicializada
2026-01-08 10:30:46 INFO database Conexão PostgreSQL estabelecida
2026-01-08 10:30:47 INFO broker RabbitMQ conectado
```

### Métricas Importantes

- **Response Time**: Tempo médio de requisição
- **Error Rate**: % de requisições com erro
- **Database Queries**: Quantidade e tempo de queries
- **Cache Hit Rate**: % de hits em Redis
- **Queue Length**: Tamanho da fila RabbitMQ

### Heatlh Check

```bash
curl http://localhost:8080/api/v1/health
```

Resposta esperada:
```json
{
  "status": "running",
  "time": "2026-01-08T10:30:00Z"
}
```

---

## 📚 Referências Adicionais

- [Documentação POST_MODERATION_RULES](docs/POST_MODERATION_RULES.md)
- [Postman Collections](postman/collections/)
- [Go Official Docs](https://golang.org/doc/)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/docs/)
- [RabbitMQ](https://www.rabbitmq.com/documentation.html)
- [Redis](https://redis.io/documentation)

---

## 👥 Contribuindo

### Workflow Git

```bash
# 1. Criar branch feature
git checkout -b feature/minha-feature

# 2. Fazer commits
git add .
git commit -m "feat: descrição da mudança"

# 3. Push e criar PR
git push origin feature/minha-feature
```

### Convenção de Commits

```
feat:    Nova funcionalidade
fix:     Correção de bug
docs:    Alterações em documentação
style:   Formatação, sem mudança de lógica
refactor: Refatoração de código
test:    Adição de testes
chore:   Tarefas de build, dependências
```

Exemplo:
```
feat: adicionar endpoint GET /features/:id
fix: corrigir validação de email em CreateUserRequest
docs: atualizar README com instruções de deploy
```

---

## 📝 Licença

Projeto pessoal.

---

## 📞 Contato

**Desenvolvedor:** Emerson Rabelo

Para dúvidas ou sugestões, abra uma issue no repositório.

---

**Última atualização:** 8 de janeiro de 2026
**Versão da Documentação:** 1.0.0
