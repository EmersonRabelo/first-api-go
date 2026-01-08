# API Reference - First API Go

Documentação completa de todos os endpoints com exemplos de uso.

---

## 📑 Índice

1. [Health Check](#health-check)
2. [Users](#users)
3. [Posts](#posts)
4. [Likes](#likes)
5. [Replies](#replies)
6. [Reports](#reports)
7. [Códigos de Status](#códigos-de-status)
8. [Tratamento de Erros](#tratamento-de-erros)

---

## Health Check

### GET /health

Verifica se a aplicação está operacional.

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/health
```

**Response (200 OK):**
```json
{
  "status": "running",
  "time": "2026-01-08T10:30:45.123456Z"
}
```

---

## Users

### GET /users

Lista todos os usuários com paginação.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/users?limit=10&offset=0' \
  -H 'Content-Type: application/json'
```

**Query Parameters:**
| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| limit | integer | Não | Quantidade de registros (padrão: 10) |
| offset | integer | Não | Posição inicial (padrão: 0) |

**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "João Silva",
    "email": "joao@example.com",
    "created_at": "2026-01-08T10:00:00Z",
    "updated_at": "2026-01-08T10:00:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Maria Santos",
    "email": "maria@example.com",
    "created_at": "2026-01-08T10:05:00Z",
    "updated_at": "2026-01-08T10:05:00Z"
  }
]
```

---

### GET /users/{id}

Obtém um usuário específico pelo ID.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000' \
  -H 'Content-Type: application/json'
```

**Path Parameters:**
| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| id | UUID | ID do usuário |

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva",
  "email": "joao@example.com",
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

**Response (404 Not Found):**
```json
{
  "error": "usuário não encontrado"
}
```

---

### POST /users

Cria um novo usuário.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com"
  }'
```

**Request Body:**
```json
{
  "name": "João Silva",
  "email": "joao@example.com"
}
```

**Body Validation:**
| Campo | Tipo | Validação | Descrição |
|-------|------|-----------|-----------|
| name | string | required, max=255 | Nome do usuário |
| email | string | required, email, unique | Email único |

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva",
  "email": "joao@example.com",
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "email inválido"
}
```

**Response (409 Conflict):**
```json
{
  "error": "email já cadastrado"
}
```

---

### PUT /users/{id}

Atualiza um usuário existente.

**Request:**
```bash
curl -X PUT 'http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "João Silva Atualizado",
    "email": "joao.novo@example.com"
  }'
```

**Path Parameters:**
| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| id | UUID | ID do usuário |

**Request Body:**
```json
{
  "name": "João Silva Atualizado",
  "email": "joao.novo@example.com"
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva Atualizado",
  "email": "joao.novo@example.com",
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:30:00Z"
}
```

---

### DELETE /users/{id}

Deleta um usuário (soft delete).

**Request:**
```bash
curl -X DELETE 'http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000'
```

**Response (204 No Content):**
```
(vazio)
```

**Response (404 Not Found):**
```json
{
  "error": "usuário não encontrado"
}
```

---

## Posts

### GET /posts

Lista todos os posts visíveis com paginação.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/posts?limit=20&offset=0' \
  -H 'Content-Type: application/json'
```

**Query Parameters:**
| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| limit | integer | Não | Quantidade de registros (padrão: 20) |
| offset | integer | Não | Posição inicial (padrão: 0) |

**Response (200 OK):**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "body": "Este é meu primeiro post!",
    "flag": "visible",
    "is_active": true,
    "created_at": "2026-01-08T10:00:00Z",
    "updated_at": "2026-01-08T10:00:00Z"
  }
]
```

---

### GET /posts/{id}

Obtém um post específico.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/posts/660e8400-e29b-41d4-a716-446655440000'
```

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "body": "Este é meu primeiro post!",
  "flag": "visible",
  "is_active": true,
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

---

### POST /posts

Cria um novo post.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "body": "Este é meu primeiro post!"
  }'
```

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "body": "Este é meu primeiro post!"
}
```

**Body Validation:**
| Campo | Tipo | Validação | Descrição |
|-------|------|-----------|-----------|
| user_id | UUID | required, exists | ID do criador |
| body | string | required, max=280 | Conteúdo do post |

**Response (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "body": "Este é meu primeiro post!",
  "flag": "visible",
  "is_active": true,
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

---

### PUT /posts/{id}

Atualiza um post existente.

**Request:**
```bash
curl -X PUT 'http://localhost:8080/api/v1/posts/660e8400-e29b-41d4-a716-446655440000' \
  -H 'Content-Type: application/json' \
  -d '{
    "body": "Conteúdo atualizado do post!"
  }'
```

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "body": "Conteúdo atualizado do post!",
  "flag": "visible",
  "is_active": true,
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:30:00Z"
}
```

---

### DELETE /posts/{id}

Deleta um post (soft delete).

**Request:**
```bash
curl -X DELETE 'http://localhost:8080/api/v1/posts/660e8400-e29b-41d4-a716-446655440000'
```

**Response (204 No Content):**
```
(vazio)
```

---

## Likes

### GET /likes

Lista todas as curtidas com paginação.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/likes?limit=50&offset=0'
```

**Response (200 OK):**
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2026-01-08T10:00:00Z"
  }
]
```

---

### GET /likes/{id}

Obtém uma curtida específica.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/likes/770e8400-e29b-41d4-a716-446655440000'
```

**Response (200 OK):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2026-01-08T10:00:00Z"
}
```

---

### POST /likes

Cria uma curtida em um post.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/likes \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000"
  }'
```

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Body Validation:**
| Campo | Tipo | Validação | Descrição |
|-------|------|-----------|-----------|
| user_id | UUID | required, exists | ID de quem curte |
| post_id | UUID | required, exists | ID do post |

**Constraints:**
- Um usuário não pode curtir o mesmo post duas vezes
- Retorna 409 Conflict se já existe like

**Response (201 Created):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2026-01-08T10:00:00Z"
}
```

**Response (409 Conflict):**
```json
{
  "error": "usuário já curtiu este post"
}
```

---

### DELETE /likes/{id}

Remove uma curtida.

**Request:**
```bash
curl -X DELETE 'http://localhost:8080/api/v1/likes/770e8400-e29b-41d4-a716-446655440000'
```

**Response (204 No Content):**
```
(vazio)
```

---

## Replies

### GET /replies

Lista todos os comentários com paginação.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/replies?limit=50&offset=0'
```

**Response (200 OK):**
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "body": "Ótimo post!",
    "is_active": true,
    "created_at": "2026-01-08T10:00:00Z",
    "updated_at": "2026-01-08T10:00:00Z"
  }
]
```

---

### GET /replies/{id}

Obtém um comentário específico.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/replies/880e8400-e29b-41d4-a716-446655440000'
```

**Response (200 OK):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "body": "Ótimo post!",
  "is_active": true,
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

---

### POST /replies

Cria um comentário em um post.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/replies \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "body": "Ótimo post!"
  }'
```

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "body": "Ótimo post!"
}
```

**Body Validation:**
| Campo | Tipo | Validação | Descrição |
|-------|------|-----------|-----------|
| user_id | UUID | required, exists | ID de quem comenta |
| post_id | UUID | required, exists | ID do post |
| body | string | required | Conteúdo do comentário |

**Response (201 Created):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "body": "Ótimo post!",
  "is_active": true,
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

---

### PUT /replies/{id}

Atualiza um comentário.

**Request:**
```bash
curl -X PUT 'http://localhost:8080/api/v1/replies/880e8400-e29b-41d4-a716-446655440000' \
  -H 'Content-Type: application/json' \
  -d '{
    "body": "Comentário atualizado!"
  }'
```

**Response (200 OK):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "body": "Comentário atualizado!",
  "is_active": true,
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:30:00Z"
}
```

---

### DELETE /replies/{id}

Deleta um comentário (soft delete).

**Request:**
```bash
curl -X DELETE 'http://localhost:8080/api/v1/replies/880e8400-e29b-41d4-a716-446655440000'
```

**Response (204 No Content):**
```
(vazio)
```

---

## Reports

### POST /posts/{id}/report

Denuncia um post para moderação.

**Request:**
```bash
curl -X POST 'http://localhost:8080/api/v1/posts/660e8400-e29b-41d4-a716-446655440000/report' \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "reason": "Conteúdo ofensivo"
  }'
```

**Path Parameters:**
| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| id | UUID | ID do post denunciado |

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "reason": "Conteúdo ofensivo"
}
```

**Body Validation:**
| Campo | Tipo | Validação | Descrição |
|-------|------|-----------|-----------|
| user_id | UUID | required, exists | ID de quem denuncia |
| reason | string | required | Motivo da denúncia |

**Response (201 Created):**
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "reason": "Conteúdo ofensivo",
  "perspective_toxicity": 0.0,
  "perspective_insult": 0.0,
  "perspective_profanity": 0.0,
  "perspective_threat": 0.0,
  "perspective_identity_attack": 0.0,
  "perspective_severe_toxicity": 0.0,
  "status": "pending",
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:00:00Z"
}
```

**Fluxo Assíncrono:**
1. Denúncia criada com `status: pending`
2. Mensagem enviada para RabbitMQ
3. Consumer processa em background (Perspective API)
4. Scores calculados, post flag atualizado
5. Report atualizado com status final

---

### GET /reports

Lista relatórios de denúncias (admin).

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/reports?limit=50&offset=0'
```

**Response (200 OK):**
```json
[
  {
    "id": "990e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "reason": "Conteúdo ofensivo",
    "perspective_toxicity": 0.85,
    "perspective_insult": 0.80,
    "perspective_profanity": 0.438,
    "perspective_threat": 0.070,
    "perspective_identity_attack": 0.102,
    "perspective_severe_toxicity": 0.354,
    "status": "limited",
    "created_at": "2026-01-08T10:00:00Z",
    "updated_at": "2026-01-08T10:01:00Z"
  }
]
```

---

### GET /reports/{id}

Obtém detalhes de um relatório específico.

**Request:**
```bash
curl -X GET 'http://localhost:8080/api/v1/reports/990e8400-e29b-41d4-a716-446655440000'
```

**Response (200 OK):**
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "660e8400-e29b-41d4-a716-446655440000",
  "reason": "Conteúdo ofensivo",
  "perspective_toxicity": 0.85,
  "perspective_insult": 0.80,
  "perspective_profanity": 0.438,
  "perspective_threat": 0.070,
  "perspective_identity_attack": 0.102,
  "perspective_severe_toxicity": 0.354,
  "status": "limited",
  "created_at": "2026-01-08T10:00:00Z",
  "updated_at": "2026-01-08T10:01:00Z"
}
```

---

## Códigos de Status HTTP

| Código | Significado | Quando Ocorre |
|--------|-------------|---------------|
| 200 | OK | Requisição bem-sucedida (GET, PUT) |
| 201 | Created | Recurso criado com sucesso (POST) |
| 204 | No Content | Operação bem-sucedida sem resposta (DELETE) |
| 400 | Bad Request | Validação falhou, dados inválidos |
| 404 | Not Found | Recurso não encontrado |
| 409 | Conflict | Violação de constraint (ex: email duplicado) |
| 500 | Internal Server Error | Erro do servidor |

---

## Tratamento de Erros

### Formato Padrão de Erro

```json
{
  "error": "descrição do erro"
}
```

### Exemplos de Erros Comuns

**400 Bad Request - JSON Inválido**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{invalid json}'
```

Resposta:
```json
{
  "error": "invalid json"
}
```

---

**400 Bad Request - Validação**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "",
    "email": "invalid-email"
  }'
```

Resposta:
```json
{
  "error": "validação falhou: name é obrigatório, email inválido"
}
```

---

**404 Not Found**
```bash
curl -X GET 'http://localhost:8080/api/v1/users/invalid-uuid'
```

Resposta:
```json
{
  "error": "usuário não encontrado"
}
```

---

**409 Conflict - Email Duplicado**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "João",
    "email": "joao@example.com"
  }'
```

Resposta:
```json
{
  "error": "email já cadastrado"
}
```

---

**409 Conflict - Like Duplicado**
```bash
curl -X POST http://localhost:8080/api/v1/likes \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "post_id": "660e8400-e29b-41d4-a716-446655440000"
  }'
```

(Assuming like already exists)

Resposta:
```json
{
  "error": "usuário já curtiu este post"
}
```

---

## Guia de Integração

### cURL

```bash
# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name": "João", "email": "joao@example.com"}'

# Criar post
curl -X POST http://localhost:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -d '{"user_id": "...", "body": "Hello!"}'

# Curtir post
curl -X POST http://localhost:8080/api/v1/likes \
  -H 'Content-Type: application/json' \
  -d '{"user_id": "...", "post_id": "..."}'

# Comentar post
curl -X POST http://localhost:8080/api/v1/replies \
  -H 'Content-Type: application/json' \
  -d '{"user_id": "...", "post_id": "...", "body": "Great!"}'

# Denunciar post
curl -X POST http://localhost:8080/api/v1/posts/660e8400.../report \
  -H 'Content-Type: application/json' \
  -d '{"user_id": "...", "reason": "Offensive"}'
```

### JavaScript/Fetch

```javascript
// Criar usuário
const createUser = async (name, email) => {
  const res = await fetch('http://localhost:8080/api/v1/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email })
  });
  return res.json();
};

// Listar posts
const getPosts = async () => {
  const res = await fetch('http://localhost:8080/api/v1/posts');
  return res.json();
};

// Curtir post
const likePost = async (userId, postId) => {
  const res = await fetch('http://localhost:8080/api/v1/likes', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, post_id: postId })
  });
  return res.json();
};
```

### Python/Requests

```python
import requests

BASE_URL = 'http://localhost:8080/api/v1'

# Criar usuário
def create_user(name, email):
    res = requests.post(f'{BASE_URL}/users', json={
        'name': name,
        'email': email
    })
    return res.json()

# Listar posts
def get_posts():
    res = requests.get(f'{BASE_URL}/posts')
    return res.json()

# Curtir post
def like_post(user_id, post_id):
    res = requests.post(f'{BASE_URL}/likes', json={
        'user_id': user_id,
        'post_id': post_id
    })
    return res.json()
```

---

**Versão:** 1.0.0  
**Data:** 8 de janeiro de 2026
