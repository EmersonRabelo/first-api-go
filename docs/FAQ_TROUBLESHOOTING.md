# FAQ e Troubleshooting

Respostas para perguntas frequentes e soluções para problemas comuns.

---

## 📑 Índice

1. [Instalação e Setup](#instalação-e-setup)
2. [Erros em Runtime](#erros-em-runtime)
3. [Problemas de Performance](#problemas-de-performance)
4. [Banco de Dados](#banco-de-dados)
5. [Fila de Mensagens](#fila-de-mensagens)
6. [Cache Redis](#cache-redis)
7. [API e Requisições](#api-e-requisições)
8. [Moderação e Relatórios](#moderação-e-relatórios)

---

## Instalação e Setup

### P: Como instalar as dependências?

**R:** Execute:
```bash
go mod download
go mod tidy
```

Isso baixará todas as dependências definidas no `go.mod`.

---

### P: Qual versão de Go é necessária?

**R:** Go 1.25.4 ou superior.

Verificar versão instalada:
```bash
go version
```

Instalar/atualizar Go:
```bash
# macOS (com Homebrew)
brew install go

# Linux
wget https://go.dev/dl/go1.25.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz
```

---

### P: Como configurar as variáveis de ambiente?

**R:** Criar arquivo `.env` na raiz do projeto:

```bash
# Copiar exemplo
cp .env.example .env

# Editar com seus valores
nano .env
```

**Variáveis obrigatórias:**
```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=first_api_go
DB_SSL_MODE=disable
BROKER_HOST=localhost
BROKER_PORT=5672
BROKER_USER=guest
BROKER_PASSWORD=guest
REDIS_HOST=localhost
REDIS_PORT=6379
```

---

### P: Erro "connection refused" ao iniciar a aplicação

**R:** Verificar se os serviços estão rodando:

```bash
# PostgreSQL
psql -h localhost -U postgres -c "SELECT 1"

# Redis
redis-cli ping

# RabbitMQ
docker logs rabbitmq-dev  # se usar Docker

# Ou iniciar com Docker Compose
docker-compose up -d
```

---

### P: "port already in use" ao iniciar aplicação

**R:** Mudar a porta ou liberar a porta ocupada:

```bash
# Mudar porta no .env
SERVER_PORT=8081

# Ou encontrar processo usando porta 8080
lsof -i :8080

# E matar o processo
kill -9 <PID>
```

---

## Erros em Runtime

### P: Error "failed to dial primary server"

**R:** PostgreSQL não está conectando. Verificar:

```bash
# 1. PostgreSQL está rodando?
docker ps | grep postgres

# 2. Credenciais estão corretas?
psql -h localhost -U postgres -d first_api_go

# 3. Logs do PostgreSQL
docker logs postgres-dev

# 4. Verificar firewall
# localhost:5432 está acessível?
nc -zv localhost 5432
```

---

### P: "panic: template is nil"

**R:** Erro interno do Gin. Verificar:

```bash
# Não usar templates HTML, use JSON
// ❌ Ruim
c.HTML(200, "template.html", data)

// ✅ Bom
c.JSON(200, data)
```

---

### P: "json: cannot unmarshal string into Go value of type UUID"

**R:** Formato de UUID inválido. Certificar que:

```bash
# UUID é formato válido (com hífens)
# ✅ Correto
550e8400-e29b-41d4-a716-446655440000

# ❌ Incorreto
550e8400e29b41d4a716446655440000

# Testar em Go
uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
```

---

### P: "runtime error: invalid memory address"

**R:** Tentando usar pointer nil. Verificar:

```go
// ❌ Ruim
var user *User
user.Name = "João"  // nil pointer dereference

// ✅ Bom
user := &User{}
user.Name = "João"

// Ou
user := new(User)
user.Name = "João"
```

---

## Problemas de Performance

### P: API respondendo lentamente

**R:** Investigar pontos de gargalo:

```bash
# 1. Checar CPU/Memória
docker stats api

# 2. Ver tempo de resposta
time curl http://localhost:8080/api/v1/health

# 3. Verificar queries lentas
# No PostgreSQL
SELECT * FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 10;

# 4. Ver query atual em execução
SELECT pid, query, query_start FROM pg_stat_activity;

# 5. Verificar índices
SELECT * FROM pg_stat_user_indexes;
```

---

### P: Curtidas (likes) estão lentas

**R:** Certificar que Redis está sendo usado:

```bash
# 1. Verificar se Redis está rodando
redis-cli ping

# 2. Checar se contadores estão em Redis
redis-cli get post:660e8400-e29b-41d4-a716-446655440000:likes

# 3. Ver se tabela post_likes_count existe
psql -d first_api_go -c "SELECT * FROM post_likes_count LIMIT 5;"

# 4. Lipar cache e refazer
redis-cli FLUSHDB
# Restart da aplicação para reconstruir
```

---

### P: Muitos logs gerando lentidão

**R:** Reduzir nível de log em produção:

```go
// Em config.go
if !IsProd() {
    gin.SetMode(gin.DebugMode)
} else {
    gin.SetMode(gin.ReleaseMode)
}
```

---

## Banco de Dados

### P: Migrations não executam automaticamente

**R:** Verificar logs da migração:

```bash
# Ver logs completos
go run cmd/api/main.go 2>&1 | grep -i migration

# Verificar migrations aplicadas
psql -d first_api_go -c "SELECT version, dirty FROM schema_migrations;"

# Limpar bancos
psql -c "DROP DATABASE first_api_go;"
psql -c "CREATE DATABASE first_api_go;"

# Restart da aplicação
go run cmd/api/main.go
```

---

### P: "UNIQUE constraint violated"

**R:** Tentando inserir dado duplicado:

```bash
# Para email duplicado
SELECT * FROM users WHERE email = 'joao@example.com';

# Limp limpar (apenas em dev!)
DELETE FROM users WHERE email = 'joao@example.com';

# Ou usar outro email
```

---

### P: "foreign key constraint violation"

**R:** Referência a um registro que não existe:

```bash
# Exemplo: criar like para post que não existe
# Verificar se post existe
SELECT * FROM posts WHERE post_id = 'uuid-aqui';

# Se não existir, criar post primeiro
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{"user_id": "...", "body": "..."}'
```

---

### P: Dados desincronizados entre tabelas

**R:** Triggers podem estar falhando:

```sql
-- Verificar se triggers existem
SELECT * FROM pg_trigger WHERE tgname LIKE '%like%';

-- Recriar trigger se necessário
-- Ver migrations 000007, 000008, etc.
```

---

## Fila de Mensagens

### P: RabbitMQ Management UI inacessível

**R:** Acessar em `http://localhost:15672`

```bash
# Credenciais padrão
Username: guest
Password: guest

# Se não funcionar, verificar
docker logs rabbitmq-dev

# Ou acessar diretamente
docker exec -it rabbitmq-dev rabbitmq-diagnostics ping
```

---

### P: Mensagens não são consumidas

**R:** Consumer pode estar travado. Verificar:

```bash
# 1. Ver se consumer está rodando
curl http://localhost:8080/api/v1/health

# Deve mostrar status "running"

# 2. Verificar fila em RabbitMQ
curl -u guest:guest http://localhost:15672/api/queues

# 3. Verificar se consumer é iniciado
# Em main.go, deve ter:
go func() {
    reportConsumer.Start()
}()

# 4. Ver logs do consumer
# Aumentar logging em queue/consumer.go
```

---

### P: Denúncias (reports) não estão sendo processadas

**R:** Problema no fluxo RabbitMQ-Consumer. Debugar:

```bash
# 1. Verificar se fila tem mensagens
curl -u guest:guest http://localhost:15672/api/queues

# 2. Ver logs do handler
# Em internal/handler/report_handler.go
log.Printf("Processando report: %v", message)

# 3. Simular call à Perspective API
# O handler pode estar falhando nessa chamada

# 4. Verificar se report foi criado no BD
psql -d first_api_go -c "SELECT * FROM reports ORDER BY created_at DESC LIMIT 5;"

# Flag do post foi atualizado?
psql -d first_api_go -c "SELECT post_id, flag FROM posts WHERE flag != 'visible';"
```

---

## Cache Redis

### P: Redis não está sendo usado

**R:** Verificar inicialização:

```bash
# 1. Redis está rodando?
docker ps | grep redis

# 2. Conectar e testar
redis-cli
> ping
PONG

# 3. Ver se chaves estão sendo criadas
redis-cli KEYS "*"

# 4. Testar manualmente
redis-cli INCR test:counter
redis-cli GET test:counter
```

---

### P: Contadores desincronizados entre Redis e PostgreSQL

**R:** Reconstruir cache:

```bash
# 1. Verificar estado em PostgreSQL
psql -d first_api_go -c "SELECT * FROM post_likes_count;"

# 2. Limpar Redis
redis-cli FLUSHDB

# 3. Restart aplicação
# Redis será reconstruído com dados do PostgreSQL

# Ou reconstruir manualmente
# redis-cli SET post:<id>:likes <count>
```

---

### P: "redis: connection refused"

**R:** Redis não está acessível:

```bash
# 1. Verificar se está rodando
docker ps | grep redis

# 2. Iniciá-lo
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 3. Testar conexão
redis-cli ping

# 4. Verificar logs
docker logs redis
```

---

## API e Requisições

### P: Erro 400 Bad Request em POST

**R:** Verificar validação:

```bash
# Exemplo - criar post sem user_id
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{"body": "Teste"}'

# Erro esperado:
# {"error": "validação falhou: user_id é obrigatório"}

# Solução: incluir user_id
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "body": "Teste"
  }'
```

---

### P: Erro 401 Unauthorized (quando implementar autenticação)

**R:** Token ausente ou inválido:

```bash
# ❌ Sem token
curl http://localhost:8080/api/v1/posts

# ✅ Com token
curl -H "Authorization: Bearer seu-token-jwt" \
  http://localhost:8080/api/v1/posts
```

---

### P: Erro 404 Not Found

**R:** Recurso não existe:

```bash
# Verificar se ID existe no BD
psql -d first_api_go -c "SELECT * FROM posts WHERE post_id = '660e8400-...';"

# Se não existir, retorna 404 corretamente
# Se existir, pode ser soft delete (deleted_at != null)
psql -d first_api_go -c "SELECT *, deleted_at FROM posts WHERE post_id = '660e8400-...';"
```

---

### P: Erro 409 Conflict

**R:** Violação de constraint (duplicado, etc):

```bash
# Email duplicado
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "João",
    "email": "existente@example.com"
  }'

# Resposta:
# {"error": "email já cadastrado"}

# Solução: usar email único
```

---

### P: CORS errors no frontend

**R:** Adicionar CORS middleware:

```go
// router/router.go
import "github.com/gin-contrib/cors"

r := gin.Default()

// Adicionar CORS
config := cors.DefaultConfig()
config.AllowOrigins = []string{"http://localhost:3000"}
r.Use(cors.New(config))

// Ou permitir tudo (dev only)
r.Use(cors.Default())
```

---

## Moderação e Relatórios

### P: Post ainda está visível após denúncia

**R:** Consumer pode não ter processado:

```bash
# 1. Verificar se report foi criado
psql -d first_api_go -c "SELECT * FROM reports ORDER BY created_at DESC LIMIT 1;"

# 2. Verificar fila RabbitMQ
curl -u guest:guest http://localhost:15672/api/queues

# 3. Ver logs do consumer
docker logs api 2>&1 | grep -i report

# 4. Consumer pode estar esperando Perspective API
# Verificar se chave da API está configurada

# 5. Força testa manualmente
# Atualizar post flag diretamente
psql -d first_api_go -c "UPDATE posts SET flag = 'limited' WHERE post_id = '...';"
```

---

### P: Scores da Perspective API zerados

**R:** API não está sendo chamada corretamente:

```bash
# 1. Verificar se API key está configurada
echo $PERSPECTIVE_API_KEY

# 2. Verificar URL e headers na chamada
# Em handler/report_handler.go

# 3. Testar a API diretamente
curl -X POST https://commentanalyzer.googleapis.com/v1/comments:analyze?key=YOUR_KEY \
  -H 'Content-Type: application/json' \
  -d '{
    "comment": {"text": "seu texto"},
    "requestedAttributes": {
      "TOXICITY": {},
      "THREAT": {}
    }
  }'

# 4. Se falhar, verificar
# - Chave válida?
# - API habilitada?
# - Rate limit atingido?
```

---

### P: Post marcado como "hidden_pending_review" incorretamente

**R:** Revisar regras de moderação:

```bash
# Regras de decisão em docs/POST_MODERATION_RULES.md

# Scores que devem ocultar:
# - THREAT > 0.7
# - IDENTITY_ATTACK > 0.8
# - SEVERE_TOXICITY > 0.9

# Verificar scores armazenados
psql -d first_api_go -c "
SELECT 
  report_id,
  perspective_threat,
  perspective_identity_attack,
  perspective_severe_toxicity
FROM reports
WHERE post_id = '...'
ORDER BY created_at DESC
LIMIT 1;"

# Ajustar scores ou regras conforme necessário
```

---

## Perguntas Gerais

### P: Como resetar banco de dados em desenvolvimento?

**R:** 
```bash
# PostgreSQL
psql -c "DROP DATABASE first_api_go_dev;"
psql -c "CREATE DATABASE first_api_go_dev;"

# Migrations serão executadas automaticamente na próxima startup
go run cmd/api/main.go

# Ou usar Docker
docker-compose down -v
docker-compose up -d
```

---

### P: Como fazer backup do banco?

**R:**
```bash
# Backup em SQL
pg_dump -U postgres first_api_go > backup.sql

# Backup comprimido
pg_dump -U postgres first_api_go | gzip > backup.sql.gz

# Restore
psql -U postgres first_api_go < backup.sql
```

---

### P: Como escalar a aplicação?

**R:** Ver [docs/DEPLOYMENT.md](DEPLOYMENT.md#scaling)

- Horizontal: múltiplas instâncias (Docker, Kubernetes)
- Vertical: aumentar recursos (CPU, RAM)
- Database: read replicas, sharding
- Cache: Redis cluster
- Queue: múltiplos consumers

---

### P: Como contribuir com código?

**R:** Ver [docs/CONTRIBUTING.md](CONTRIBUTING.md)

1. Fork repositório
2. Criar branch (`git checkout -b feature/...`)
3. Fazer commits com mensagens descritivas
4. Enviar PR com descrição clara
5. Responder feedback e mergear

---

### P: Onde reportar bugs?

**R:** 
1. Checar se já existe issue
2. Abrir nova issue com:
   - Descrição do problema
   - Steps para reproduzir
   - Comportamento esperado
   - Logs relevantes
3. Labels: `bug`, `priority: high/medium/low`

---

### P: Como solicitar nova feature?

**R:**
1. Abrir issue com `feature request` label
2. Descrever:
   - Caso de uso
   - Comportamento esperado
   - Possíveis implementações
3. Aguardar feedback e aprovação
4. Implementar e enviar PR

---

## Recursos Adicionais

- [README.md](../README.md) - Visão geral e setup
- [API_REFERENCE.md](API_REFERENCE.md) - Documentação de endpoints
- [ARCHITECTURE.md](ARCHITECTURE.md) - Decisões e padrões
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deploy e operações
- [CONTRIBUTING.md](CONTRIBUTING.md) - Padrões de código
- [POST_MODERATION_RULES.md](POST_MODERATION_RULES.md) - Regras de moderação

---

**Versão:** 1.0.0  
**Data:** 8 de janeiro de 2026

*Última atualização: 8 de janeiro de 2026*
