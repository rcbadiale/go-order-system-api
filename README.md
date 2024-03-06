# Go Expert Order System API

Desafio de aplicação dos conhecimentos sobre REST, gRPC e GraphQL aplicando os
conceitos de arquitetura limpa.

## Executando a aplicação

### Docker Compose

A aplicação pode ser executada diretamente utilizando containers, para isso
basta executar o Docker Compose.

```shell
docker compose up
```

Ao executar dessa maneira as dependências abaixo serão criadas nas portas em
conjunto com a aplicação:

- MySQL: 3306
- RabbitMQ: 5672 e 15672

**Obs.**: O MySQL tem seu volume montado na pasta `.docker/mysql` na raíz do
projeto.

**Obs.**: A aplicação faz a migração do banco de dados de forma automática caso
seja necessário.

## Utilização das interfaces

### REST

Porta definida no `docker-compose.yml`: 8000

Requisições definidas no arquivo [order.http](api/order.http).

### gRPC

Porta definida no `docker-compose.yml`: 50051

Acessando o serviço com o [evans](https://github.com/ktr0731/evans):

```shell
evans -r repl
> package pb
> service OrderService
```

Criando uma nova ordem:

```
> call CreateOrder
id => grpc1
price => 10
tax => 10
```

Listando ordens existentes:

```
> call ListOrders
```

### GraphQL

Porta definida no `docker-compose.yml`: 8080

Acessar a [web interface](http://localhost:8080/).

Criando uma nova ordem:

```graphql
mutation createOrder {
  createOrder(input: { id: "graphql1", Price: 20.0, Tax: 20.0 }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

Listando ordens existentes:

```graphql
query getOrders {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

## Geradores de código

Nessa seção são tratadas as ferramentas geradoras de código automatizado para
cada tipo de interface desenvolvida na aplicação.

**Obs.**: Os comandos apresentados são apenas referências e devem ser usados
apenas em tempo de desenvolvimento em casos de alterações nas
interfaces/dependências.

### Gerenciador de dependências

Para gerar os códigos de gerenciamento de dependencias é utilizada a ferramenta
[wire](https://github.com/google/wire).

Comando de referência para gerar os códigos de gerenciamento de dependencias.

```bash
cd cmd/order_system
wire gen
```

### gRPC

Para gerar os códigos é utilizada a ferramenta
[protoc](https://github.com/protocolbuffers/protobuf).

Comando de referência para gerar os código relacionados ao gRPC.

```bash
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
```

### GraphQL

Para gerar os códigos é utilizada a ferramenta
[gqlgen](https://github.com/99designs/gqlgen).

Comando de referência para gerar os código relacionados ao GraphQL.

```bash
go run github.com/99designs/gqlgen generate
```
