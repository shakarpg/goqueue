# 🚀 GoQueue: Sistema de Fila de Jobs e Processamento Assíncrono em Go

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="GoLang">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/GitHub_Actions-2671E5?style=for-the-badge&logo=githubactions&logoColor=white" alt="GitHub Actions">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License MIT">
</p>

Este projeto implementa um sistema de **fila de jobs** e **processamento assíncrono** em **Go (Golang)**, demonstrando o uso eficiente de **goroutines**, **channels**, **worker pools** e **graceful shutdown**. É uma solução robusta para lidar com tarefas que consomem tempo, como envio de e-mails, geração de relatórios PDF ou processamento de imagens, sem bloquear a thread principal da aplicação.

## ✨ Funcionalidades

*   **Fila de Jobs Concorrente:** Utiliza `channels` para gerenciar a fila de jobs de forma segura e eficiente.
*   **Worker Pool:** Um conjunto de `goroutines` processa os jobs em paralelo, otimizando o uso de recursos.
*   **Processamento Assíncrono:** Permite que a aplicação responda rapidamente enquanto tarefas demoradas são executadas em segundo plano.
*   **Graceful Shutdown:** Garante que os jobs em andamento sejam concluídos antes do desligamento da aplicação, evitando perda de dados.
*   **API RESTful:** Interface HTTP para enfileirar novos jobs, consultar o status de jobs existentes e obter estatísticas da fila.
*   **Logging Estruturado:** Integração com `zap` para logs de alta performance e fácil análise.

## 🏗️ Arquitetura

A arquitetura do `GoQueue` é modular e baseada em componentes que se comunicam de forma assíncrona. O diagrama abaixo ilustra o fluxo de um job desde a sua criação até o processamento:

```mermaid
graph TD
    Client[Cliente HTTP] -- POST /jobs --> Router[Router/Handlers]
    Router -- Enqueue --> Queue[Job Queue - Channel]
    Queue -- Dequeue --> WP[Worker Pool]
    
    subgraph "Worker Pool (Goroutines)"
        W1[Worker 1]
        W2[Worker 2]
        W3[Worker N]
    end
    
    WP --> W1
    WP --> W2
    WP --> W3
    
    W1 -- Process --> Job1[Email/PDF/Image]
    W2 -- Process --> Job2[Email/PDF/Image]
    W3 -- Process --> Job3[Email/PDF/Image]
    
    Job1 -- Update Status --> Storage[In-Memory Storage - Map + RWMutex]
    Job2 -- Update Status --> Storage
    Job3 -- Update Status --> Storage
    
    Router -- GET /jobs/:id --> Storage
    Router -- GET /stats --> Storage
```

### Componentes Principais

*   **`JobQueue` (`internal/queue`):** Gerencia a fila de jobs usando um `channel` para comunicação entre o produtor (API) e os consumidores (workers). Utiliza um `map` e `sync.RWMutex` para armazenar o estado dos jobs em memória de forma segura para concorrência.
*   **`WorkerPool` (`internal/worker`):** Responsável por criar e gerenciar um pool de `goroutines` (workers) que consomem jobs da `JobQueue` e os processam. Cada worker simula o processamento de diferentes tipos de jobs (e-mail, PDF, imagem).
*   **`Router` (`internal/router`):** Define as rotas da API HTTP usando o pacote `net/http` padrão do Go, encaminhando as requisições para os `handlers` apropriados.
*   **`Handlers` (`internal/handlers`):** Contém a lógica para receber requisições HTTP, criar jobs e interagir com a `JobQueue`.
*   **`Job` (`internal/models`):** Estrutura que representa um job, incluindo seu ID, tipo, payload, status e timestamps.

## 🚀 Como Rodar o Projeto

### Pré-requisitos

*   [Go](https://golang.org/doc/install) (versão 1.18 ou superior)
*   [Docker](https://docs.docker.com/get-docker/) (opcional, para rodar via Docker)

### 1. Clonar o Repositório

```bash
git clone https://github.com/shakarpg/goqueue.git
cd goqueue
```

### 2. Instalar as dependências

```bash
go mod tidy
```

### 3. Rodar a aplicação

```bash
go run cmd/main.go
```

O servidor será iniciado na porta `8080` (ou na porta definida pela variável de ambiente `PORT`).

### 4. Rodar com Docker

```bash
docker build -t goqueue .
docker run -p 8080:8080 goqueue
```

## 🔌 Endpoints da API

O `GoQueue` expõe os seguintes endpoints:

*   **`POST /jobs`**
    *   **Descrição:** Enfileira um novo job para processamento assíncrono.
    *   **Corpo da Requisição (JSON):**
        ```json
        {
            "type": "email", // ou "pdf", "image"
            "payload": {
                "to": "teste@example.com",
                "subject": "Assunto do Email",
                "body": "Corpo do email"
            }
        }
        ```
    *   **Resposta (JSON):**
        ```json
        {
            "id": "uuid-do-job",
            "status": "pending"
        }
        ```

*   **`GET /jobs/{id}`**
    *   **Descrição:** Retorna o status e detalhes de um job específico.
    *   **Exemplo:** `GET /jobs/a1b2c3d4-e5f6-7890-1234-567890abcdef`
    *   **Resposta (JSON):**
        ```json
        {
            "id": "uuid-do-job",
            "type": "email",
            "status": "completed", // ou "pending", "running", "failed"
            "createdAt": "2023-10-27T10:00:00Z",
            "startedAt": "2023-10-27T10:00:05Z",
            "endedAt": "2023-10-27T10:00:10Z",
            "result": "Email enviado para teste@example.com",
            "error": ""
        }
        ```

*   **`GET /stats`**
    *   **Descrição:** Retorna estatísticas gerais da fila de jobs.
    *   **Resposta (JSON):**
        ```json
        {
            "total": 10,
            "pending": 2,
            "running": 3,
            "completed": 5,
            "failed": 0
        }
        ```

## 🛑 Graceful Shutdown

O `GoQueue` implementa um mecanismo de *graceful shutdown*. Isso significa que, ao receber um sinal de interrupção (`SIGINT` ou `SIGTERM`), o servidor HTTP para de aceitar novas requisições, mas aguarda um tempo configurável (10 segundos) para que as requisições e jobs em processamento sejam concluídos. Os workers também são sinalizados para finalizar suas tarefas atuais antes de encerrar, garantindo que nenhum job seja perdido durante o desligamento.

## 🧪 Rodar os Testes

```bash
make test
```

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

## 🧠 Conceitos Demonstrados

### 🔹 Goroutines & Channels
*   Workers rodando concorrentemente.
*   Comunicação eficiente via channels.
*   Uso de `select` statement para cancelamento e multiplexação.

### 🔹 Worker Pool Pattern
*   Pool de workers processando jobs em paralelo.
*   Distribuição automática de carga.
*   Processamento assíncrono para tarefas de longa duração.

### 🔹 Context & Graceful Shutdown
*   Uso de `context.Context` para cancelamento de goroutines.
*   Captura de sinais `SIGINT`/`SIGTERM` para desligamento controlado.
*   Shutdown gracioso do servidor HTTP para evitar interrupções abruptas.

### 🔹 Concurrency-Safe Storage
*   Utilização de `sync.Mutex` e `sync.RWMutex` para acesso seguro a dados compartilhados em ambientes concorrentes.

## 📈 Próximos Passos (Melhorias Potenciais)

*   [ ] Persistência em Redis ou banco de dados para jobs (atualmente em memória).
*   [ ] Retry automático para jobs falhados.
*   [ ] Priorização de jobs na fila.
*   [ ] Rate limiting por tipo de job ou cliente.
*   [ ] Dashboard web para visualização do status da fila e workers.
*   [ ] Webhooks para notificação de conclusão de jobs.
*   [ ] Suporte a jobs agendados (cron).

## 🤝 Contribuições

Contribuições são muito bem-vindas! Se você tiver ideias para melhorias, correções de bugs ou novas funcionalidades, sinta-se à vontade para:

1.  Fazer um **fork** deste repositório.
2.  Criar uma nova **branch** (`git checkout -b feature/minha-feature`).
3.  Realizar suas alterações e **commit** (`git commit -am 'feat: adiciona nova funcionalidade'`).
4.  Enviar suas alterações para o seu fork (`git push origin feature/minha-feature`).
5.  Abrir um **Pull Request** detalhando suas mudanças.

## 📄 Licença

Este projeto está licenciado sob a **MIT License**. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👤 Autor

**Rafael Galhardo**  
GitHub: [@shakarpg](https://github.com/shakarpg)
LinkedIn: [linkedin.com/in/rpg2011](https://linkedin.com/in/rpg2011)
