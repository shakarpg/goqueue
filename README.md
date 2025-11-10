# 🔄 GoQueue

![Go Version](https://img.shields.io/badge/Go-1.22-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build](https://github.com/shakarpg/goqueue/actions/workflows/go.yml/badge.svg)

Sistema de **fila de jobs** e **processamento assíncrono** escrito em **Golang**, demonstrando uso de **goroutines**, **channels**, **worker pools**, e **graceful shutdown**.

---

## 🚀 Tecnologias

- **Go 1.22**
- **Chi Router** (rotas HTTP)
- **Goroutines & Channels** (concorrência)
- **Worker Pool Pattern**
- **Zap** (logs estruturados)
- **Context** (graceful shutdown)
- **Docker**
- **GitHub Actions** (CI/CD)

---

## 🎯 Funcionalidades

✅ API REST para criar e gerenciar jobs  
✅ Fila de jobs em memória com channels  
✅ Worker pool com 5 workers concorrentes  
✅ Suporte a 3 tipos de jobs: `email`, `pdf`, `image`  
✅ Endpoint de métricas (`/api/metrics`)  
✅ Graceful shutdown (SIGINT/SIGTERM)  
✅ Logs estruturados com Zap  
✅ Testes automatizados  

---

## 🧰 Como rodar o projeto

### 1️⃣ Clone o repositório
```bash
git clone https://github.com/shakarpg/goqueue.git
cd goqueue
```

### 2️⃣ Instale as dependências
```bash
go mod tidy
```

### 3️⃣ Rode a aplicação
```bash
make run
```

Acesse: [http://localhost:8080/health](http://localhost:8080/health)

---

## 🧪 Rodar os testes

```bash
make test
```

---

## 📡 Endpoints da API

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Health check |
| POST | `/api/jobs` | Criar novo job |
| GET | `/api/jobs` | Listar todos os jobs |
| GET | `/api/jobs/{id}` | Obter job específico |
| GET | `/api/metrics` | Estatísticas da fila |

---

## 🧾 Exemplo de uso

### 1. Criar job de email
```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {
      "to": "user@example.com",
      "subject": "Hello",
      "body": "Test email"
    }
  }'
```

**Resposta:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "email",
  "status": "pending",
  "payload": {
    "to": "user@example.com",
    "subject": "Hello",
    "body": "Test email"
  },
  "created_at": "2025-11-10T10:00:00Z"
}
```

### 2. Criar job de PDF
```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "pdf",
    "payload": {
      "filename": "report.pdf",
      "content": "Report data"
    }
  }'
```

### 3. Verificar status do job
```bash
curl http://localhost:8080/api/jobs/550e8400-e29b-41d4-a716-446655440000
```

### 4. Ver métricas
```bash
curl http://localhost:8080/api/metrics
```

**Resposta:**
```json
{
  "total": 10,
  "pending": 2,
  "running": 1,
  "completed": 6,
  "failed": 1
}
```

---

## 🐳 Docker

### Build
```bash
make docker-build
```

### Run
```bash
make docker-run
```

---

## 📂 Estrutura do Projeto

```
goqueue/
├── .github/
│   └── workflows/
│       └── go.yml           # GitHub Actions CI/CD
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── models/
│   │   └── job.go           # Modelo de Job
│   ├── queue/
│   │   └── queue.go         # Fila de jobs
│   ├── worker/
│   │   └── worker.go        # Worker pool
│   ├── handlers/
│   │   └── jobs.go          # Handlers HTTP
│   └── router/
│       └── router.go        # Configuração de rotas
├── tests/
│   └── queue_test.go        # Testes
├── .env
├── .gitignore
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

---

## 🧠 Conceitos demonstrados

### 🔹 Goroutines & Channels
- Workers rodando concorrentemente
- Comunicação via channels
- Select statement para cancelamento

### 🔹 Worker Pool Pattern
- Pool de 5 workers processando jobs
- Distribuição automática de carga
- Processamento assíncrono

### 🔹 Context & Graceful Shutdown
- Context para cancelamento de goroutines
- Captura de sinais SIGINT/SIGTERM
- Shutdown gracioso do servidor HTTP

### 🔹 Concurrency-Safe Storage
- Mutex para acesso seguro ao storage
- RWMutex para otimizar leituras

---

## 🧠 Próximos passos (melhorias)

- [ ] Persistência em Redis ou banco de dados
- [ ] Retry automático para jobs falhados
- [ ] Priorização de jobs
- [ ] Rate limiting por tipo de job
- [ ] Dashboard web para visualização
- [ ] Webhooks para notificação de conclusão
- [ ] Suporte a jobs agendados (cron)

---

## 📄 Licença

MIT License - sinta-se livre para usar e modificar!

---

## 👤 Autor

**Shakarpg**  
GitHub: [@shakarpg](https://github.com/shakarpg)
