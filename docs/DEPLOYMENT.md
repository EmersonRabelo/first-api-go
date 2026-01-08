# Deployment e Operações - First API Go

Guia completo para deploy, monitoramento e manutenção da aplicação.

---

## 📑 Índice

1. [Environments](#environments)
2. [Build e Packaging](#build-e-packaging)
3. [Docker](#docker)
4. [Deployment](#deployment)
5. [Monitoring](#monitoring)
6. [Troubleshooting](#troubleshooting)
7. [Backup e Recovery](#backup-e-recovery)
8. [Scaling](#scaling)

---

## Environments

### Development

**Propósito:** Desenvolvimento local

**Arquivo:** `.env.development` (ou `.env`)

```env
ENVIRONMENT=development
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=first_api_go_dev
DB_SSL_MODE=disable

# Message Broker
BROKER_HOST=localhost
BROKER_PORT=5672
BROKER_USER=guest
BROKER_PASSWORD=guest

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

**Iniciar Serviços:**

```bash
# PostgreSQL
docker run -d --name postgres-dev \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=first_api_go_dev \
  -p 5432:5432 \
  postgres:15-alpine

# Redis
docker run -d --name redis-dev \
  -p 6379:6379 \
  redis:7-alpine

# RabbitMQ
docker run -d --name rabbitmq-dev \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3.12-management-alpine
```

**Iniciar Aplicação:**

```bash
go run cmd/api/main.go
```

---

### Staging

**Propósito:** Ambiente pré-produção para testes

**Arquivo:** `.env.staging`

```env
ENVIRONMENT=staging
SERVER_PORT=8080

# Database
DB_HOST=staging-db.internal
DB_PORT=5432
DB_USER=api_staging
DB_PASSWORD=${DB_PASSWORD_STAGING}
DB_NAME=first_api_go_staging
DB_SSL_MODE=require

# Message Broker
BROKER_HOST=staging-rabbitmq.internal
BROKER_PORT=5672
BROKER_USER=api_staging
BROKER_PASSWORD=${BROKER_PASSWORD_STAGING}

# Redis
REDIS_HOST=staging-redis.internal
REDIS_PORT=6379
```

**Deploy:**

```bash
# Build binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o build/api cmd/api/main.go

# Deploy via Docker/Kubernetes
docker push myregistry.azurecr.io/first-api-go:staging
kubectl set image deployment/first-api-go-staging \
  first-api-go=myregistry.azurecr.io/first-api-go:staging
```

---

### Production

**Propósito:** Ambiente de produção

**Arquivo:** `.env.production` (variáveis de ambiente no servidor)

```env
ENVIRONMENT=production
SERVER_PORT=8080

# Database
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_USER=api_prod
DB_PASSWORD=${DB_PASSWORD_PROD}
DB_NAME=first_api_go
DB_SSL_MODE=require

# Message Broker
BROKER_HOST=prod-rabbitmq.example.com
BROKER_PORT=5672
BROKER_USER=api_prod
BROKER_PASSWORD=${BROKER_PASSWORD_PROD}

# Redis
REDIS_HOST=prod-redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD_PROD}

# Segurança
JWT_SECRET=${JWT_SECRET}
API_KEY=${API_KEY}
```

**Variáveis Sensíveis:**
- Usar Azure Key Vault, AWS Secrets Manager, ou similares
- Nunca commitar arquivos `.env` em produção
- Rotacionar secrets regularmente

---

## Build e Packaging

### Compilação Local

```bash
# Build padrão
go build -o bin/api cmd/api/main.go

# Build otimizado para produção
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build \
    -ldflags="-s -w" \
    -o bin/api cmd/api/main.go

# Verificar tamanho
ls -lh bin/api

# Testar binary
./bin/api
```

### Versionamento

```bash
# Adicionar versão ao binary
VERSION=$(git describe --tags --always)
go build \
  -ldflags="-X main.Version=${VERSION}" \
  -o bin/api cmd/api/main.go
```

**No código:**
```go
var Version = "dev"

func init() {
    fmt.Printf("API Version: %s\n", Version)
}
```

### Cross-compilation

```bash
# Para Windows
GOOS=windows GOARCH=amd64 go build -o api.exe cmd/api/main.go

# Para macOS
GOOS=darwin GOARCH=amd64 go build -o api-mac cmd/api/main.go

# Para Linux ARM (Raspberry Pi)
GOOS=linux GOARCH=arm GOARM=7 go build -o api-rpi cmd/api/main.go
```

---

## Docker

### Dockerfile

```dockerfile
# Multi-stage build para otimizar tamanho
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Dependências de build
RUN apk add --no-cache git

# Copy go.mod e go.sum
COPY go.mod go.sum ./

# Download dependências
RUN go mod download

# Copy código fonte
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o api cmd/api/main.go

# Stage final - runtime
FROM alpine:3.18

# Instalar CA certificates para HTTPS
RUN apk add --no-cache ca-certificates

# Criar usuário não-root
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copy binary do builder
COPY --from=builder /app/api .

# Mudar proprietário
RUN chown -R appuser:appuser /home/appuser

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

CMD ["./api"]
```

### Docker Compose (Local Development)

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: api-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: first_api_go_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: api-redis
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    container_name: api-rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - "5672:5672"
      - "15672:15672"
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: api
    ports:
      - "8080:8080"
    environment:
      ENVIRONMENT: development
      SERVER_PORT: 8080
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: first_api_go_dev
      DB_SSL_MODE: disable
      BROKER_HOST: rabbitmq
      BROKER_PORT: 5672
      BROKER_USER: guest
      BROKER_PASSWORD: guest
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

volumes:
  postgres_data:
```

**Usar Docker Compose:**

```bash
# Iniciar todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Parar
docker-compose down

# Remover dados
docker-compose down -v
```

### Build e Push para Registry

```bash
# Build com tag
docker build -t myregistry.azurecr.io/first-api-go:1.0.0 .

# Login no registry
docker login myregistry.azurecr.io

# Push
docker push myregistry.azurecr.io/first-api-go:1.0.0

# Pull em outro lugar
docker pull myregistry.azurecr.io/first-api-go:1.0.0
```

---

## Deployment

### Linux (systemd)

**1. Criar arquivo de serviço:**

```bash
sudo nano /etc/systemd/system/first-api-go.service
```

**Conteúdo:**

```ini
[Unit]
Description=First API Go
After=network.target
After=postgresql.service
After=redis.service
After=rabbitmq-server.service

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/first-api-go
EnvironmentFile=/opt/first-api-go/.env
ExecStart=/opt/first-api-go/api
Restart=on-failure
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true

[Install]
WantedBy=multi-user.target
```

**2. Criar usuário:**

```bash
sudo useradd -r -s /bin/false appuser
```

**3. Copiar arquivos:**

```bash
sudo mkdir -p /opt/first-api-go
sudo cp bin/api /opt/first-api-go/
sudo cp .env.production /opt/first-api-go/.env
sudo chown -R appuser:appuser /opt/first-api-go
sudo chmod 600 /opt/first-api-go/.env
```

**4. Iniciar serviço:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable first-api-go
sudo systemctl start first-api-go

# Ver status
sudo systemctl status first-api-go

# Ver logs
sudo journalctl -u first-api-go -f
```

---

### Kubernetes

**deployment.yaml:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: first-api-go
  labels:
    app: first-api-go
spec:
  replicas: 3
  selector:
    matchLabels:
      app: first-api-go
  template:
    metadata:
      labels:
        app: first-api-go
    spec:
      containers:
      - name: api
        image: myregistry.azurecr.io/first-api-go:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: ENVIRONMENT
          value: production
        - name: SERVER_PORT
          value: "8080"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: api-config
              key: db_host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: db_password
        # ... mais variáveis
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: first-api-go
spec:
  selector:
    app: first-api-go
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
```

**Deploy:**

```bash
kubectl apply -f deployment.yaml

# Ver deployments
kubectl get deployments

# Ver pods
kubectl get pods

# Ver logs
kubectl logs -f deployment/first-api-go

# Scale
kubectl scale deployment first-api-go --replicas=5

# Update image
kubectl set image deployment/first-api-go \
  api=myregistry.azurecr.io/first-api-go:2.0.0

# Rollback
kubectl rollout undo deployment/first-api-go
```

---

### Cloud Providers

#### Azure App Service

```bash
# Login
az login

# Criar grupo de recursos
az group create -n first-api-go-rg -l eastus

# Criar App Service Plan
az appservice plan create \
  -n first-api-go-plan \
  -g first-api-go-rg \
  --sku B2

# Criar App Service
az webapp create \
  -n first-api-go \
  -g first-api-go-rg \
  -p first-api-go-plan \
  --runtime "Go|1.25"

# Deploy via ZIP
az webapp deployment source config-zip \
  -n first-api-go \
  -g first-api-go-rg \
  --src build.zip

# Configure variáveis
az webapp config appsettings set \
  -n first-api-go \
  -g first-api-go-rg \
  --settings \
    ENVIRONMENT=production \
    DB_HOST=prod-db.example.com \
    DB_PASSWORD=@Microsoft.KeyVault(...) \
```

#### AWS EC2

```bash
# Conectar via SSH
ssh -i key.pem ubuntu@ec2-instance

# Instalar Go
sudo wget https://go.dev/dl/go1.25.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Instalar dependências
sudo apt-get update
sudo apt-get install -y postgresql redis-server rabbitmq-server

# Deploy
git clone https://github.com/EmersonRabelo/first-api-go.git
cd first-api-go
go build -o api cmd/api/main.go

# Executar com supervisord
sudo apt-get install supervisor
# ... configurar /etc/supervisor/conf.d/api.conf
sudo service supervisor restart
```

#### Google Cloud Run

```dockerfile
# Dockerfile otimizado para Cloud Run
FROM golang:1.25.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api cmd/api/main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/api .
ENV PORT 8080
EXPOSE 8080
CMD ["./api"]
```

```bash
# Deploy para Cloud Run
gcloud run deploy first-api-go \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars=ENVIRONMENT=production,DB_HOST=...
```

---

## Monitoring

### Health Checks

```bash
# Manual
curl http://localhost:8080/api/v1/health

# Contínuo
watch -n 5 'curl -s http://localhost:8080/api/v1/health | jq'
```

### Logs

**Ver logs da aplicação:**

```bash
# Via systemd
sudo journalctl -u first-api-go -f

# Via Docker
docker logs -f first-api-go

# Via Kubernetes
kubectl logs -f pod/first-api-go-xxx
```

**Rotação de logs:**

```bash
# /etc/logrotate.d/first-api-go
/var/log/first-api-go/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 appuser appuser
    sharedscripts
    postrotate
        systemctl reload first-api-go > /dev/null 2>&1 || true
    endscript
}
```

### Métricas

**Implementar Prometheus:**

```go
// internal/metrics/metrics.go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "path"},
    )
)

func RegisterMetrics() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}
```

**Middleware:**

```go
r.Use(func(c *gin.Context) {
    start := time.Now()
    
    c.Next()
    
    duration := time.Since(start).Seconds()
    httpRequestsTotal.WithLabelValues(
        c.Request.Method,
        c.Request.URL.Path,
        fmt.Sprintf("%d", c.Writer.Status()),
    ).Inc()
    
    httpRequestDuration.WithLabelValues(
        c.Request.Method,
        c.Request.URL.Path,
    ).Observe(duration)
})

// Endpoint para Prometheus
r.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))
```

---

## Troubleshooting

### Problema: "Connection refused" do PostgreSQL

```bash
# Verificar se PostgreSQL está rodando
docker ps | grep postgres

# Testar conexão
psql -h localhost -U postgres -d first_api_go_dev -c "SELECT 1"

# Ver logs
docker logs postgres-dev

# Verificar conectividade
docker exec api ping postgres
```

---

### Problema: RabbitMQ não conecta

```bash
# Verificar se RabbitMQ está rodando
docker ps | grep rabbitmq

# Ver logs
docker logs api-rabbitmq

# Testar conexão
docker exec -it api-rabbitmq rabbitmq-diagnostics ping

# Management UI
curl -u guest:guest http://localhost:15672/api/overview
```

---

### Problema: Redis desincronizado

```bash
# Conectar ao Redis
redis-cli

# Ver todas as chaves
keys *

# Ver contadores
get post:660e8400-e29b-41d4-a716-446655440000:likes

# Limpar cache (reconstruir)
FLUSHDB

# Depois: restart aplicação para reconstruir
```

---

### Problema: Performance lenta

```bash
# Analisar queries PostgreSQL
docker exec -it api-postgres psql -U postgres -d first_api_go_dev

-- Enable query logging
ALTER DATABASE first_api_go_dev SET log_statement = 'all';
ALTER DATABASE first_api_go_dev SET log_duration = on;
ALTER DATABASE first_api_go_dev SET log_min_duration_statement = 1000;

-- Ver slow queries
SELECT * FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 10;
```

---

### Problema: Alto uso de memória

```bash
# Docker
docker stats api

# Kubernetes
kubectl top pod first-api-go-xxx

# Aumentar limite
# deployment.yaml
resources:
  limits:
    memory: "1Gi"
```

---

## Backup e Recovery

### PostgreSQL Backup

```bash
# Backup completo
docker exec api-postgres pg_dump -U postgres first_api_go_dev > backup.sql

# Backup comprimido
docker exec api-postgres pg_dump -U postgres first_api_go_dev | gzip > backup.sql.gz

# Backup via schedule
0 2 * * * /usr/bin/docker exec api-postgres pg_dump -U postgres first_api_go_dev | gzip > /backups/db_$(date +\%Y\%m\%d_\%H\%M\%S).sql.gz
```

### PostgreSQL Restore

```bash
# Restore completo
docker exec -i api-postgres psql -U postgres first_api_go_dev < backup.sql

# Restore comprimido
gunzip < backup.sql.gz | docker exec -i api-postgres psql -U postgres first_api_go_dev
```

### Estratégia de Backup

```yaml
# backup-policy.yaml (Kubernetes)
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
spec:
  schedule: "0 2 * * *"  # 2 AM diário
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:15
            command:
            - /bin/sh
            - -c
            - pg_dump -h postgres -U postgres first_api_go | gzip > /backups/db_$(date +%Y%m%d_%H%M%S).sql.gz
            volumeMounts:
            - name: backups
              mountPath: /backups
          volumes:
          - name: backups
            persistentVolumeClaim:
              claimName: backup-pvc
          restartPolicy: OnFailure
```

---

## Scaling

### Horizontal Scaling

```bash
# Kubernetes
kubectl scale deployment first-api-go --replicas=10

# Docker Swarm
docker service scale first-api-go=10

# Verificar distribuição
kubectl get pods -o wide
```

### Vertical Scaling

```yaml
# Aumentar recursos por pod
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

### Database Scaling

```yaml
# Read Replica no PostgreSQL
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: first-api-go-db
spec:
  instances: 3
  postgresql:
    parameters:
      max_connections: "500"
  storage:
    size: 100Gi
```

### Redis Scaling

```bash
# Redis Cluster (3+ nodes)
redis-cli --cluster create \
  127.0.0.1:6379 \
  127.0.0.1:6380 \
  127.0.0.1:6381
```

---

**Versão:** 1.0.0  
**Data:** 8 de janeiro de 2026
