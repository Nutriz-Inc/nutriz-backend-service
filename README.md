# Nutriz Backend Service

Serviço backend do Nutriz: API REST para gestão de doações de leite materno (usuárias, endereços, doações, etapas do processo, pontos de doação, agenda interna e dashboard de métricas).

---

## Pré-requisitos

- Docker
- Docker Compose

Para a opção de execução local (sem Docker para a aplicação em si), também é necessário:

- Go (https://go.dev/doc/install)
- Make

---

## Compatibilidade

Este serviço roda apenas em:

- Linux
- macOS

Caso esteja utilizando Windows, utilize o WSL:

```bash
wsl -d ubuntu
```

Depois execute todos os comandos dentro do Ubuntu (WSL).

---

## Opção A — Executar tudo via Docker (recomendado)

Este fluxo não exige Go instalado na máquina: build da aplicação, banco de dados, migrations e execução acontecem inteiramente em containers.

### 1. Clone o projeto

```bash
git clone <url-do-repositorio>
cd nutriz-backend-service
```

### 2. Construa a imagem da aplicação

```bash
docker compose build
```

### 3. Suba o banco de dados e o Redis

```bash
docker compose up -d database redis
```

### 4. Rode as migrations (cria as tabelas e aplica os seeds)

```bash
docker compose run --rm app ./migrate
```

### 5. Inicie o container da aplicação

```bash
docker compose up -d app
```

A API estará disponível em `http://localhost:3333`.

### 6. Acesse e teste a API

- Documentação interativa (Swagger UI): `http://localhost:3333/docs`
- Especificação OpenAPI (JSON): `http://localhost:3333/openapi.json`
- Health check: `http://localhost:3333/health`

Veja a seção [Testando a API](#testando-a-api) para exemplos de requisição.

### 7. Interromper e remover os containers

```bash
# Para os containers, mantendo os dados do banco
docker compose down

# Para os containers e remove também os volumes (reset completo dos dados)
docker compose down -v
```

---

## Opção B — Aplicação local com Go, banco via Docker

Use esta opção se preferir rodar a aplicação diretamente com `go run` (hot reload via `make run`), mantendo apenas o banco e o Redis em containers.

### 1. Clone o projeto e instale as dependências

```bash
git clone <url-do-repositorio>
cd nutriz-backend-service
go mod tidy
```

### 2. Suba o banco de dados e o Redis

```bash
docker compose up -d database redis
```

### 3. Rode as migrations

```bash
make migrations
```

### 4. Rode o serviço

```bash
make run
```

A API estará disponível em `http://localhost:3333` (Swagger em `http://localhost:3333/docs`).

### Fluxo completo (resumo)

```bash
go mod tidy
docker compose up -d database redis
make migrations
make run
```

### Executando testes

```bash
make test
```

Após executar testes que criam novos dados no banco, pode ocorrer conflito entre execuções. Para evitar problemas, derrube os containers removendo os volumes e suba novamente:

```bash
docker compose down -v
docker compose up -d database redis
```

---

## Testando a API

### Versionamento e documentação

As rotas são agrupadas por nível de acesso (`/public` para rotas sem autenticação, `/internal` para rotas autenticadas). A documentação completa de todos os endpoints — parâmetros, corpo de requisição, respostas e códigos de erro — está disponível via Swagger/OpenAPI em `http://localhost:3333/docs` assim que a aplicação estiver no ar.

### Autenticação

Rotas `/internal/*` exigem um JWT no cabeçalho `Authorization: Bearer <token>`, obtido no login. O cabeçalho `action-by`, usado para autorização por tipo de usuária, é preenchido automaticamente pelo middleware a partir do token — não precisa ser enviado manualmente.

### Exemplo: cadastro de usuária (autocadastro, rota pública)

```bash
curl -X POST http://localhost:3333/public/user \
  -H "Content-Type: application/json" \
  -d '{
    "type": "common",
    "name": "Joana Souza",
    "cpf": "57694673800",
    "email": "joana.souza@email.com",
    "password": "12345678",
    "phone_number": "+5511991111111",
    "birth_date": "1990-01-01",
    "address": { "zip_code": "01001000" },
    "consent_log": {
      "terms_version": "1.0",
      "ip_address": "127.0.0.1",
      "user_agent": "curl/8.0"
    }
  }'
```

### Exemplo: login

```bash
curl -X POST http://localhost:3333/public/auth/login \
  -H "Content-Type: application/json" \
  -d '{ "email": "joana.souza@email.com", "password": "12345678" }'
```

A resposta traz um `token` JWT. Use-o nas chamadas às rotas `/internal/*`:

### Exemplo: buscar a própria usuária (rota autenticada)

```bash
curl http://localhost:3333/internal/user/<id_user> \
  -H "Authorization: Bearer <token>"
```

### Exemplo: listar pontos de doação (rota pública)

```bash
curl "http://localhost:3333/public/donation/point?page=1&page_size=25"
```

Para os demais endpoints (doações, etapas de doação, jobs, dashboard), consulte a lista completa com exemplos de corpo de requisição, parâmetros e respostas em `http://localhost:3333/docs`.

---

## Observações

- Certifique-se de que as portas utilizadas pelos containers estejam livres (`3333` para a API, `5435` para o Postgres, `6398` para o Redis).
- Caso altere variáveis de ambiente ou banco de dados, reinicie os containers.
- O `.env.development` já vem com valores padrão para desenvolvimento; não é necessário criá-lo manualmente.
- O comando abaixo remove completamente os dados locais do ambiente Docker:

```bash
docker compose down -v
```
