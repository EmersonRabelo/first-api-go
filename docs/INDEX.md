# Índice de Documentação

Guia completo para navegar na documentação do projeto First API Go.

---

## 📚 Documentação Disponível

### 🚀 Para Começar

**[README.md](../README.md)** - Visão geral do projeto
- Descrição do projeto
- Tecnologias utilizadas
- Setup rápido
- Estrutura de diretórios
- Endpoints básicos

**[ARCHITECTURE.md](ARCHITECTURE.md)** - Arquitetura do software
- Estilo arquitetural
- Fluxo de requisições
- Componentes principais
- Padrões de design
- Diagrama E-R

---

### 📖 Referência Técnica

**[API_REFERENCE.md](API_REFERENCE.md)** - Documentação completa de endpoints
- Todos os endpoints com exemplos
- Request/Response formats
- Codes de status HTTP
- Tratamento de erros
- Guia de integração (cURL, JS, Python)

**[DEPLOYMENT.md](DEPLOYMENT.md)** - Deploy e operações
- Setup de environments (dev, staging, prod)
- Docker e Docker Compose
- Deployment em várias plataformas
- Monitoring e logs
- Troubleshooting
- Backup e recovery

---

### 🛠️ Desenvolvimento

**[CONTRIBUTING.md](CONTRIBUTING.md)** - Guia de contribuição
- Code of conduct
- Como contribuir
- Padrões de código Go
- Convenções de naming
- Estrutura de commits
- Pull requests
- Testes e cobertura
- Documentação de código

**[FAQ_TROUBLESHOOTING.md](FAQ_TROUBLESHOOTING.md)** - Perguntas e problemas comuns
- Setup e instalação
- Erros em runtime
- Performance
- Banco de dados
- RabbitMQ
- Redis
- API
- Moderação

---

### 📋 Funcionalidades Específicas

**[POST_MODERATION_RULES.md](POST_MODERATION_RULES.md)** - Regras de moderação
- Política de moderação
- Scores da Perspective API
- Estados de posts
- Regras de decisão
- Priorização de revisão

---

## 🗺️ Mapa de Documentação por Tópico

### Setup e Instalação
1. [README.md - Configuração e Instalação](../README.md#-configuração-e-instalação)
2. [DEPLOYMENT.md - Environments](DEPLOYMENT.md#environments)
3. [FAQ - Instalação e Setup](FAQ_TROUBLESHOOTING.md#instalação-e-setup)

### API e Endpoints
1. [README.md - API Endpoints](../README.md#-api-endpoints)
2. [API_REFERENCE.md - Documentação Completa](API_REFERENCE.md)
3. [FAQ - API e Requisições](FAQ_TROUBLESHOOTING.md#api-e-requisições)

### Arquitetura
1. [README.md - Arquitetura](../README.md#-arquitetura)
2. [ARCHITECTURE.md - Documentação Detalhada](ARCHITECTURE.md)
3. [README.md - Estrutura do Projeto](../README.md#-estrutura-do-projeto)

### Banco de Dados
1. [README.md - Banco de Dados](../README.md#-banco-de-dados)
2. [ARCHITECTURE.md - Diagrama E-R](ARCHITECTURE.md#-diagrama-de-entidades-expandido)
3. [DEPLOYMENT.md - Backup](DEPLOYMENT.md#backup-e-recovery)
4. [FAQ - Banco de Dados](FAQ_TROUBLESHOOTING.md#banco-de-dados)

### Sistema de Fila (RabbitMQ)
1. [README.md - Sistema de Fila](../README.md#-sistema-de-fila)
2. [ARCHITECTURE.md - Message Queue Layer](ARCHITECTURE.md#4-message-queue-layer-rabbitmq)
3. [FAQ - RabbitMQ](FAQ_TROUBLESHOOTING.md#fila-de-mensagens)

### Cache (Redis)
1. [README.md - Cache e Redis](../README.md#-cache-e-redis)
2. [ARCHITECTURE.md - Cache Layer](ARCHITECTURE.md#5-cache-layer-redis)
3. [FAQ - Redis](FAQ_TROUBLESHOOTING.md#cache-redis)

### Moderação
1. [README.md - Fluxo de Moderação](../README.md#-fluxo-de-moderação)
2. [POST_MODERATION_RULES.md - Regras Detalhadas](POST_MODERATION_RULES.md)
3. [FAQ - Moderação](FAQ_TROUBLESHOOTING.md#moderação-e-relatórios)

### Desenvolvimento
1. [README.md - Guia de Desenvolvimento](../README.md#-guia-de-desenvolvimento)
2. [CONTRIBUTING.md - Padrões de Código](CONTRIBUTING.md)
3. [FAQ - Troubleshooting](FAQ_TROUBLESHOOTING.md)

### Deployment e Operações
1. [DEPLOYMENT.md - Guia Completo](DEPLOYMENT.md)
2. [ARCHITECTURE.md - Performance](ARCHITECTURE.md#-performance)
3. [README.md - Monitoring](../README.md#-monitoramento-e-logs)

---

## 👥 Documentação por Perfil

### Para Product Manager

**Entender o sistema:**
1. [README.md - Visão Geral](../README.md#-visão-geral)
2. [README.md - Entidades do Negócio](../README.md#-entidades-do-negócio)
3. [POST_MODERATION_RULES.md](POST_MODERATION_RULES.md)

**Acompanhar progresso:**
1. [ARCHITECTURE.md - Decisões Arquiteturais](ARCHITECTURE.md#-decisões-arquiteturais)
2. [DEPLOYMENT.md - Monitoramento](DEPLOYMENT.md#monitoring)

---

### Para Desenvolvedor Backend

**Começar:**
1. [README.md - Setup](../README.md#-configuração-e-instalação)
2. [README.md - Estrutura do Projeto](../README.md#-estrutura-do-projeto)
3. [CONTRIBUTING.md - Padrões de Código](CONTRIBUTING.md#padrões-de-código)

**Implementar features:**
1. [README.md - Guia de Desenvolvimento](../README.md#-guia-de-desenvolvimento)
2. [CONTRIBUTING.md - Estrutura de um Novo Endpoint](CONTRIBUTING.md#estrutura-de-um-novo-endpoint)
3. [ARCHITECTURE.md - Padrões de Design](ARCHITECTURE.md#-padrões-de-design-utilizados)

**Debugar:**
1. [FAQ_TROUBLESHOOTING.md](FAQ_TROUBLESHOOTING.md)
2. [DEPLOYMENT.md - Troubleshooting](DEPLOYMENT.md#troubleshooting)

---

### Para Desenvolvedor Frontend

**Entender API:**
1. [API_REFERENCE.md - Documentação de Endpoints](API_REFERENCE.md)
2. [README.md - API Endpoints](../README.md#-api-endpoints)
3. [API_REFERENCE.md - Guia de Integração](API_REFERENCE.md#guia-de-integração)

**Testar:**
1. [API_REFERENCE.md - Exemplos cURL](API_REFERENCE.md)
2. [README.md - Postman Collections](../README.md)

---

### Para DevOps/SRE

**Deploy:**
1. [DEPLOYMENT.md - Deployment](DEPLOYMENT.md#deployment)
2. [DEPLOYMENT.md - Docker](DEPLOYMENT.md#docker)
3. [DEPLOYMENT.md - Kubernetes](DEPLOYMENT.md#kubernetes)

**Operações:**
1. [DEPLOYMENT.md - Monitoring](DEPLOYMENT.md#monitoring)
2. [DEPLOYMENT.md - Troubleshooting](DEPLOYMENT.md#troubleshooting)
3. [DEPLOYMENT.md - Scaling](DEPLOYMENT.md#scaling)

**Backup:**
1. [DEPLOYMENT.md - Backup e Recovery](DEPLOYMENT.md#backup-e-recovery)

---

### Para Reviewer/Maintainer

**Code Review:**
1. [CONTRIBUTING.md - Padrões de Código](CONTRIBUTING.md#padrões-de-código)
2. [CONTRIBUTING.md - Pull Requests](CONTRIBUTING.md#pull-requests)
3. [ARCHITECTURE.md - Padrões de Design](ARCHITECTURE.md#-padrões-de-design-utilizados)

**Testes:**
1. [CONTRIBUTING.md - Testes](CONTRIBUTING.md#testes)
2. [CONTRIBUTING.md - Cobertura](CONTRIBUTING.md#cobertura-de-testes)

---

## 📋 Checklist de Documentação

### Para Novo Desenvolvedor

- [ ] Ler [README.md](../README.md)
- [ ] Fazer setup seguindo [Configuração e Instalação](../README.md#-configuração-e-instalação)
- [ ] Entender [Arquitetura](ARCHITECTURE.md)
- [ ] Testar alguns [Endpoints](API_REFERENCE.md)
- [ ] Ler [Padrões de Código](CONTRIBUTING.md#padrões-de-código)
- [ ] Implementar primeira feature seguindo [Guia de Desenvolvimento](../README.md#-guia-de-desenvolvimento)

### Para Deploy em Produção

- [ ] Revisar [DEPLOYMENT.md - Environments](DEPLOYMENT.md#environments)
- [ ] Preparar [variáveis de ambiente](DEPLOYMENT.md#production)
- [ ] Escolher [estratégia de deploy](DEPLOYMENT.md#deployment)
- [ ] Configurar [monitoring](DEPLOYMENT.md#monitoring)
- [ ] Testar [health checks](DEPLOYMENT.md#health-checks)
- [ ] Preparar [backup strategy](DEPLOYMENT.md#backup-e-recovery)

### Para Troubleshooting

- [ ] Checar [FAQ_TROUBLESHOOTING.md](FAQ_TROUBLESHOOTING.md)
- [ ] Revisar logs relevantes
- [ ] Consultar [Troubleshooting no DEPLOYMENT.md](DEPLOYMENT.md#troubleshooting)
- [ ] Abrir issue se necessário

---

## 🔍 Busca Rápida

### "Como fazer..."

| Pergunta | Resposta |
|----------|----------|
| Como instalar? | [README.md - Setup](../README.md#-configuração-e-instalação) |
| Como criar um endpoint? | [CONTRIBUTING.md](CONTRIBUTING.md#estrutura-de-um-novo-endpoint) |
| Como testar? | [CONTRIBUTING.md - Testes](CONTRIBUTING.md#testes) |
| Como fazer deploy? | [DEPLOYMENT.md](DEPLOYMENT.md) |
| Como monitorar? | [DEPLOYMENT.md - Monitoring](DEPLOYMENT.md#monitoring) |
| Como debugar? | [FAQ_TROUBLESHOOTING.md](FAQ_TROUBLESHOOTING.md) |
| Como contribuir? | [CONTRIBUTING.md](CONTRIBUTING.md) |
| Como escalar? | [DEPLOYMENT.md - Scaling](DEPLOYMENT.md#scaling) |
| Como fazer backup? | [DEPLOYMENT.md - Backup](DEPLOYMENT.md#backup-e-recovery) |
| Como reportar bug? | [FAQ - Bugs](FAQ_TROUBLESHOOTING.md#onde-reportar-bugs) |

---

## 📞 Contato e Suporte

**Para dúvidas sobre documentação:**
- Abrir issue no repositório
- Sugerir melhorias
- Reportar erros ou informações desatualizadas

**Para reportar bugs:**
- Usar [FAQ - Reportar bugs](FAQ_TROUBLESHOOTING.md#onde-reportar-bugs)

**Para solicitar features:**
- Usar [FAQ - Solicitar features](FAQ_TROUBLESHOOTING.md#como-solicitar-nova-feature)

---

## 🔄 Atualização de Documentação

Todos são encorajados a:
- Corrigir erros de digitação
- Atualizar exemplos desatualizados
- Adicionar exemplos para casos confusos
- Melhorar clareza e estrutura

Veja [CONTRIBUTING.md - Melhorias de Documentação](CONTRIBUTING.md#melhorias-de-documentação)

---

## 📊 Estrutura de Documentos

```
docs/
├── INDEX.md (este arquivo)
├── README.md (visão geral)
├── ARCHITECTURE.md (design)
├── API_REFERENCE.md (endpoints)
├── DEPLOYMENT.md (deploy & ops)
├── CONTRIBUTING.md (padrões)
├── POST_MODERATION_RULES.md (regras)
└── FAQ_TROUBLESHOOTING.md (problemas)
```

**Total:** 8 documentos abrangentes
**Linhas:** ~8000+ linhas de documentação
**Cobertura:** 95%+ dos tópicos relevantes

---

## 🎓 Recursos de Aprendizado

### Padrões e Práticas

- [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Project Layout](https://github.com/golang-standards/project-layout)
- [Conventional Commits](https://www.conventionalcommits.org/)

### Frameworks e Bibliotecas

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Redis](https://redis.io/documentation)
- [RabbitMQ](https://www.rabbitmq.com/documentation.html)

### Desenvolvimento

- [Go Official](https://golang.org/)
- [Go Packages](https://pkg.go.dev/)
- [Docker Docs](https://docs.docker.com/)
- [Kubernetes Docs](https://kubernetes.io/docs/)

---

**Versão:** 1.0.0  
**Data:** 8 de janeiro de 2026

*Última atualização: 8 de janeiro de 2026*
