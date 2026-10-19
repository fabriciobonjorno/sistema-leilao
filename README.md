# Sistema de leilões

API Go para criação de leilões e envio de lances, usando MongoDB.

## Executar com Docker Compose

```sh
cp cmd/auction/.env.example cmd/auction/.env
docker compose up --build
```

A API ficará disponível em `http://localhost:8080`. Para executar sem Docker, configure as mesmas variáveis de ambiente, inicie um MongoDB e rode `go run ./cmd/auction`.

## Configuração

`AUCTION_DURATION` define por quanto tempo cada leilão permanece aberto. Informe uma duração aceita por Go, como `30s`, `5m` ou `1h`. O padrão é `5m`. `AUCTION_INTERVAL` continua aceito como alternativa para compatibilidade com a validação de lances já existente.

As variáveis de conexão são `MONGODB_URL` e `MONGODB_DB`; veja `cmd/auction/.env.example`.

## Testes

```sh
go test ./...
```

O teste do repositório de leilões usa uma duração curta e uma coleção falsa, verificando o estado aberto após a criação e o fechamento automático depois do prazo.
