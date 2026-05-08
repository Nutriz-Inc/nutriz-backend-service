# Nutriz Backend Service

Serviço backend do Nutriz.

---

## Pré-requisitos

Antes de iniciar, você precisa ter instalado na máquina:

- Go
- Docker
- Docker Compose
- Make

### Instalação do Go

https://go.dev/doc/install

---

## Compatibilidade

Este serviço roda apenas em:

- Linux
- macOS

Caso esteja utilizando Windows, utilize o WSL:

```bash
wsl -d ubuntu
```

Depois execute todos os os comandos dentro do Ubuntu (WSL).

---

## Instalação do projeto

### 1. Clone o projeto

```bash
git clone <url-do-repositorio>
cd <nome-do-projeto>
```

### 2. Instale as dependências do Go

```bash
go mod tidy
```

---

## Subindo os serviços

Inicie os containers necessários:

```bash
docker compose up -d
```

---

## Executando as migrations

```bash
make migrations
```

---

## Rodando o serviço

```bash
make run
```

---

## Fluxo completo

```bash
go mod tidy

docker compose up -d

make migrations

make run
```

---

## Executando testes

Após executar testes que criam novos dados no banco, pode ocorrer conflito entre execuções.

Para evitar problemas, derrube os containers removendo os volumes:

```bash
docker compose down -v
```

Depois suba novamente os serviços:

```bash
docker compose up -d
```

---

## Observações

- Certifique-se de que as portas utilizadas pelos containers estejam livres.
- Caso altere variáveis de ambiente ou banco de dados, reinicie os containers.
- O comando abaixo remove completamente os dados locais do ambiente Docker:

```bash
docker compose down -v
```
