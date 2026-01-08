# Arquitetura de Software - First API Go

Documentação detalhada sobre a arquitetura, padrões e decisões tecnológicas.

## 📐 Visão Geral da Arquitetura

### Estilo Arquitetural: Layered Architecture + Event-Driven

A aplicação segue um modelo em camadas tradicional com componentes orientados a eventos para operações assíncronas.

```
┌─────────────────────────────────────────────────────┐
│             Presentation Layer (HTTP)               │
│  Controllers (Gin) - Request parsing & validation   │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│           Business Logic Layer (Services)           │
│  Orquestra, regras de negócio, chamadas de dados   │
└──────────────────────┬──────────────────────────────┘
                       │
         ┌─────────────┼──────────────┐
         │             │              │
    ┌────▼────┐  ┌─────▼────┐  ┌─────▼────────┐
    │ Data    │  │ Message  │  │ Cache        │
    │ Access  │  │ Queue    │  │ Layer        │
    │ Layer   │  │ Layer    │  │              │
    │(Repos)  │  │(Broker)  │  │(Redis)       │
    └────┬────┘  └─────┬────┘  └─────┬────────┘
         │             │              │
    ┌────┴─────────────┼──────────────┴─────┐
    │                  │                    │
┌───▼────┐   ┌────────▼──────┐   ┌────────▼──┐
│PostgreSQL  │   RabbitMQ    │   │  Redis   │
│   (GORM)   │   (Messages)  │   │ (Cache)  │
└────────────┴───────────────┴───────────────┘
```

---

## 🔄 Fluxo de Requisição

### Exemplo: POST /api/v1/posts

```
1. HTTP Request chegada
   └─> Gin Router intercepta

2. Router matches route
   └─> Chama PostHandler.Create()

3. PostHandler
   ├─> Parse JSON → DTO
   ├─> Validação (struct tags)
   ├─> Chama PostService.CreatePost()
   └─> Retorna JSON response

4. PostService
   ├─> Validações de negócio
   ├─> Chama PostRepository.Create()
   ├─> Atualiza cache (Redis)
   ├─> Publica evento (opcional)
   └─> Retorna entity

5. PostRepository
   ├─> Monta query GORM
   ├─> Executa INSERT PostgreSQL
   └─> Retorna resultado

6. HTTP Response
   └─> 201 Created com dados
```

**Timeline:** ~50-200ms (depende de I/O)

---

## 🏗️ Componentes Principais

### 1. HTTP API Layer (Gin)

**Responsabilidades:**
- Roteamento de requisições
- Parsing de JSON
- Validação de headers
- Formatação de respostas
- Tratamento de erros HTTP

**Arquivos:**
- [router/router.go](router/router.go): Definição de rotas
- [internal/controller/*.go](internal/controller/): Handlers

**Exemplo:**
```go
// Rota
posts := v1.Group("/posts")
posts.POST("", postHandler.Create)

// Handler
func (h *PostHandler) Create(c *gin.Context) {
    var req dtos.CreatePostRequest
    c.ShouldBindJSON(&req)  // Parse + validação
    
    post, err := h.service.CreatePost(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, post)
}
```

### 2. Business Logic Layer (Services)

**Responsabilidades:**
- Validações de regras de negócio
- Orquestração de operações
- Chamadas a repositórios e serviços externos
- Transações de negócio

**Arquivos:**
- [internal/service/*.go](internal/service/): Serviços de domínio

**Exemplo - PostService:**
```go
type PostService struct {
    repository repository.PostRepository
    userService *UserService
}

func (s *PostService) CreatePost(ctx context.Context, req *dtos.CreatePostRequest) (*entity.Post, error) {
    // 1. Validar usuário existe
    user, err := s.userService.FindById(ctx, req.UserId)
    if err != nil {
        return nil, fmt.Errorf("usuário não encontrado")
    }
    
    // 2. Validar conteúdo
    if len(req.Body) > 280 {
        return nil, fmt.Errorf("post muito longo")
    }
    
    // 3. Persistir
    post := &entity.Post{
        Id:       uuid.New(),
        UserId:   user.Id,
        Body:     req.Body,
        Flag:     entity.Visible,
        IsActive: true,
    }
    
    if err := s.repository.Create(ctx, post); err != nil {
        return nil, err
    }
    
    return post, nil
}
```

### 3. Data Access Layer (Repositories)

**Responsabilidades:**
- Executar queries (GORM)
- Mapear dados (Entity ↔ Database)
- Tratamento de transações
- Erros de banco

**Arquivos:**
- [internal/repository/*.go](internal/repository/): DAOs

**Exemplo - PostRepository:**
```go
type PostRepository interface {
    Create(ctx context.Context, post *entity.Post) error
    FindById(ctx context.Context, id uuid.UUID) (*entity.Post, error)
    FindAll(ctx context.Context, limit, offset int) ([]*entity.Post, error)
    Update(ctx context.Context, post *entity.Post) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type postRepository struct {
    db *gorm.DB
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) error {
    return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) FindById(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
    var post entity.Post
    err := r.db.WithContext(ctx).
        Where("post_id = ?", id).
        First(&post).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("post não encontrado")
    }
    return &post, err
}
```

### 4. Message Queue Layer (RabbitMQ)

**Responsabilidades:**
- Publicar eventos assíncronos
- Consumir mensagens em background
- Garantir entrega (acknowledgments)
- Roteamento por tópicos

**Arquivos:**
- [internal/queue/producer.go](internal/queue/producer.go): Publicador
- [internal/queue/consumer.go](internal/queue/consumer.go): Consumidor
- [internal/handler/report_handler.go](internal/handler/report_handler.go): Processador

**Fluxo Exemplo - Report:**
```go
// Producer (síncrono)
func (s *ReportService) CreateReport(...) error {
    report := &entity.Report{ ... }
    s.repository.Create(ctx, report)
    
    // Publica em RabbitMQ
    s.producer.Publish(reportMessage)
}

// Consumer (assíncrono em goroutine)
go reportConsumer.Start() // Em main.go

// Handler (processa mensagem)
func (h *ReportHandler) Handle(message []byte) error {
    // Análise com Perspective API
    scores := perspectiveAPI.Analyze(post.Body)
    
    // Atualiza post
    if scores.Threat > 0.7 {
        post.Flag = entity.HiddenPendingReview
    }
    
    repository.Update(post)
}
```

### 5. Cache Layer (Redis)

**Responsabilidades:**
- Armazenar contadores (likes, replies)
- Melhorar performance de reads frequentes
- Cache de sessão (futuro)

**Arquivos:**
- [internal/redis/client.go](internal/redis/client.go): Conexão
- [internal/redis/counter.go](internal/redis/counter.go): Operações

**Exemplo - Like Counter:**
```go
// Ao curtir
func (s *LikeService) CreateLike(ctx context.Context, req *dtos.CreateLikeRequest) error {
    like := &entity.Like{...}
    s.repository.Create(ctx, like)
    
    // Incrementa cache
    s.redisClient.IncrementLikeCount(like.PostId) // post:{postId}:likes +1
}

// Trigger PostgreSQL atualiza tabela post_likes_count
// Ambos permanecem sincronizados
```

---

## 🏛️ Padrões de Design Utilizados

### 1. Repository Pattern

Abstrai acesso a dados, permitindo trocar implementação sem afetar camadas superiores.

```
Service ──── Interface ──── PostgreSQL (via GORM)
         PostRepository    Implementation
```

**Benefícios:**
- Fácil testar com mocks
- Trocar banco sem alterar service
- Independência de framework ORM

### 2. Dependency Injection

Dependências injetadas no main, não criadas dentro de classes.

```go
// ✅ Bom - DI
userRepository := repository.NewUserRepository(db)
userService := service.NewUserService(userRepository)
userHandler := controller.NewUserHandler(userService)

// ❌ Ruim - Sem DI
type UserHandler struct {}
func (h *UserHandler) Create(...) {
    service := NewUserService() // Cria dentro!
}
```

### 3. DTO Pattern (Data Transfer Objects)

Estruturas específicas para validação de entrada HTTP, desacopladas de entities.

```go
// DTO - para requisição HTTP
type CreatePostRequest struct {
    UserId string `json:"user_id" validate:"required,uuid"`
    Body   string `json:"body" validate:"required,max=280"`
}

// Entity - para banco de dados
type Post struct {
    Id     uuid.UUID
    UserId uuid.UUID
    Body   string
    Flag   ProcessFlag
}

// Converter DTO → Entity
post := &Post{
    Id:     uuid.New(),
    UserId: uuid.MustParse(req.UserId),
    Body:   req.Body,
    Flag:   Visible,
}
```

### 4. Producer-Consumer Pattern

Desacopla produção de eventos do processamento, permite escalabilidade.

```
POST /report
    ↓
ReportService
    ↓
RabbitMQ (fila)
    ↓
Consumer (background)
    ↓
PerspectiveAPI
    ↓
PostRepository (update)
```

### 5. Template Method (via Gin Middleware - potencial)

Middleware para cross-cutting concerns:
```go
r.Use(func(c *gin.Context) {
    start := time.Now()
    c.Next()
    duration := time.Since(start)
    log.Printf("[%s] %s %d %dms", 
        c.Request.Method, c.Request.URL, c.Writer.Status(), duration.Milliseconds())
})
```

---

## 🗄️ Decisões Arquiteturais

### D1: PostgreSQL + GORM

**Decidido por:**
- Dados estruturados com relacionamentos
- ACID compliance
- Suporte a triggers (para contadores)
- GORM abstrai SQL, fácil testar com mocks

**Alternativas consideradas:**
- MongoDB: Menos apropriado para dados relacionais
- Redis: Apenas cache, não banco principal

---

### D2: RabbitMQ para Moderação

**Decidido por:**
- Operação de moderação é intensiva (chamada API externa)
- Não bloqueia requisição de denúncia
- Permite retry automático
- Escalável: múltiplos consumers

**Alternativas consideradas:**
- Síncrono: Bloqueava requisição, experiência ruim
- Redis pubsub: Sem persistência, mensagens perdidas
- Kafka: Overkill, complexo para volume atual

---

### D3: Redis para Contadores

**Decidido por:**
- Leitura/escrita rápida (em-memória)
- Suporta operações atômicas (INCR, DECR)
- Triggers PostgreSQL sincronizam dados
- Fallback: pode reconstruir de tabelas

**Alternativas consideradas:**
- Apenas PostgreSQL: Mais lento para contadores
- Cache bidirecional: Complexo, risco desincronização

---

### D4: Gin Framework

**Decidido por:**
- Rápido (compilado)
- Simples API
- Middleware support
- Ativo na comunidade Go

**Alternativas consideradas:**
- Echo: Similar, escolha entre preferências
- Fiber: Mais novo, menos maduro à época
- standard library: Muito verbose

---

## 📊 Diagrama de Entidades Expandido

```
┌─────────────┐
│   USERS     │
├─────────────┤
│ id (PK)     │
│ name        │
│ email       │
│ created_at  │
│ updated_at  │
└───┬───────┬─┘
    │       │
    │       └──────────────────┐
    │                          │
┌───▼──────┐          ┌────────▼─────┐
│  POSTS   │◄─┐       │   REPLIES    │
├──────────┤  │       ├──────────────┤
│ id (PK)  │  │   ┌───│ id (PK)      │
│ user_id  │──┤   │   │ user_id (FK) │
│ body     │  │   │   │ post_id (FK) │
│ flag     │  │   │   │ body         │
│ is_active│  │   │   │ is_active    │
│ created  │  │   │   │ created      │
└────┬─────┘  │   │   └──────────────┘
     │        │   │
┌────▼────┐   │   │
│  LIKES  │───┤   │
├─────────┤   │   │
│ id (PK) │   │   │
│ user_id │   │   │
│ post_id │◄──┘   │
│ created │       │
└────┬────┘       │
     │            │
┌────▼──────────────┘
│
│   ┌──────────────┐
│   │ POST_LIKES   │  ← Tabela desnormalizada para performance
│   │ _COUNT       │     (triggers atualizam automaticamente)
│   ├──────────────┤
│   │ post_id (FK) │
│   │ count        │
│   └──────────────┘

┌──────────────────────────────────┐
│       REPORTS                    │  ← Denúncias p/ moderação
├──────────────────────────────────┤
│ id (PK)                          │
│ user_id (FK)  ┐                  │
│ post_id (FK)  ├─ Quem denuncia o quê
│ reason        ┘                  │
│ perspective_* │                  │
│               ├─ Scores de análise Perspective API
│ status        ┘                  │
│ created_at                       │
└──────────────────────────────────┘
```

---

## 🔐 Segurança (Checklist)

- [ ] SQL Injection: GORM previne com prepared statements
- [ ] XSS: JSON responses (não HTML templates)
- [ ] CSRF: Stateless REST API
- [ ] Auth: Implementar JWT (não feito ainda)
- [ ] Rate Limiting: Implementar no Gin
- [ ] Validação: DTOs com struct tags
- [ ] Erros: Não expor detalhes internos (customizar responses)

**Próximos passos de segurança:**
```go
// Middleware de autenticação
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sem token"})
            return
        }
        // Verificar JWT
        c.Next()
    }
}

// Rate limiting
limiter := rate.NewLimiter(rate.Every(time.Second), 100) // 100 req/s
if !limiter.Allow() {
    c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{})
}
```

---

## ⚡ Performance

### Otimizações Implementadas

1. **Índices PostgreSQL**
   - PK automático
   - FK em relacionamentos
   - Unique: `users(email)`, `likes(user_id, post_id)`

2. **Redis Cache para Contadores**
   - Evita queries frequentes
   - O(1) incremento/decremento

3. **GORM Eager Loading** (quando necessário)
   ```go
   db.Preload("User").Find(&posts)
   ```

4. **Paginação**
   - Limit/offset em FindAll
   - Evita retornar datasets grandes

### Benchmarks (Estimados)

| Operação | Tempo |
|----------|-------|
| POST /posts | 50-100ms |
| GET /posts/:id | 10-20ms |
| POST /likes | 30-50ms (com Redis) |
| POST /reports | 50-100ms (assíncrono) |

### Scalability

**Horizontal:**
- API stateless → múltiplas instâncias
- PostgreSQL: read replicas (futuro)
- RabbitMQ: múltiplos consumers
- Redis: sentinel/cluster (futuro)

**Vertical:**
- Connection pooling em GORM
- Worker pool para goroutines

---

## 🧪 Testabilidade

### Injeção de Dependência

```go
// Fácil mockar em testes
type PostServiceTest struct {
    mockRepo *MockPostRepository
    service  *PostService
}

func (t *PostServiceTest) TestCreate() {
    t.mockRepo.On("Create").Return(nil)
    
    post, err := t.service.CreatePost(ctx, req)
    
    assert.NoError(t, err)
    assert.NotNil(t, post)
}
```

### DTOs com Validação

```go
// Validação automática com struct tags
type CreatePostRequest struct {
    UserId string `json:"user_id" validate:"required,uuid"`
    Body   string `json:"body" validate:"required,max=280"`
}

// No handler
if err := validate.Struct(req); err != nil {
    // retorna 400 Bad Request
}
```

---

## 📈 Evolução Futura

### Curto Prazo (v1.1)
- [ ] Autenticação JWT
- [ ] Rate limiting
- [ ] Testes unitários
- [ ] Documentação Swagger

### Médio Prazo (v1.2)
- [ ] Autenticação OAuth2
- [ ] Filtros avançados (busca, ordenação)
- [ ] Notificações (WebSocket/SSE)
- [ ] CI/CD (GitHub Actions)

### Longo Prazo (v2.0)
- [ ] Microserviços (moderação separada)
- [ ] Event sourcing
- [ ] CQRS
- [ ] GraphQL alternative

---

## 🎯 Métricas de Saúde

```go
type HealthMetrics struct {
    Status          string
    Uptime          time.Duration
    DBConnections   int
    RedisPing       time.Duration
    RabbitMQStatus  string
    TotalRequests   int64
    ErrorRate       float64
}

// GET /health/metrics (admin only)
```

---

**Versão:** 1.0.0  
**Data:** 8 de janeiro de 2026
