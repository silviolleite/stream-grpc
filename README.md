Examples of microservices using Apache Kafka and gRPC with Golang
===================================================

It is an example of the microservices created using Event Sourcing and CQRS. The services are written in Golang. It is using [Apache Kafka](https://kafka.apache.org/) and [gRPC](https://grpc.io/) endpoints.

**The server and client demonstrate:**
 - Connect services with gRPC
 - gRPC Services definitions are in `transactions/transactions.proto`.
 
**The procuder and consumer demonstrate:**
 - Asynchronous communications with kafka flow


## Services architecture

![](/assets/services.png)

--------

### Environment variables

We use [dotenv](https://github.com/joho/godotenv) to configure the application. All documentation for the project environment variables are located in `.env.example`.

Create a copy of `.env.example` renaming it to `.env` and fill the environment variables according to documentation. Check this file for detailed instructions.

# Dependencies

### gRPC

Package [grpc](https://godoc.org/google.golang.org/grpc) implements an RPC system called gRPC.

### Kafka-go

It provides both low and high level APIs for interacting with Kafka, mirroring concepts and implementing interfaces of the Go standard library to make it easy to use and integrate with existing software.
[kafka-go](github.com/segmentio/kafka-go)

### GORM

[GORM](https://gorm.io/) is fantastic ORM library for Golang.


### Proto-lens

API for protocol buffers using modern Haskell language and library patterns.

**Installing protoc**

In order to build Haskell packages with proto-lens, the Google protobuf compiler (which is a standalone binary named protoc) needs to be installed. 

You can see how to install it in [http://google.github.io/proto-lens/installing-protoc.html](http://google.github.io/proto-lens/installing-protoc.html)

**Generating the Classes**

Once the software is installed, there are the steps to using it. First you must compile the protocol buffer definitions and then import them, with the support library, into your program.

To compile the protocol buffer definition, run protoc with the `--go_out` parameter set to the directory you want to output the Go code to.

In our case is:
```sh
protoc -I transactions --go_out=plugins=grpc:transactions transactions/transactions.proto
```
The generated files will be suffixed `.pb.go` on `transactions` folder.

### SQlite

You must create a new empty file `sqlite.db` in the root of the project or create a file with another name and set it on environment variable `SQLITE_PATH` in `.env`. 

### Kafka

Install `landoop/fast-data-dev`

It is a Kafka distribution with Apache Kafka, Kafka Connect, Zookeeper, Confluent Schema Registry and REST Proxy

**Via docker command**
```shell script
docker pull landoop/fast-data-dev
```

Run
```shell script
docker run –rm -it -p 2181:2181 -p 3030:3030 -p 8081:8081 -p 8082:8082 -p 8083:8083 -p 9092:9092 -e ZK_PORT=2181 -e WEB_PORT=3030 -e REGISTRY_PORT=8081 -e REST_PORT=8082 -e CONNECT_PORT=8083 -e BROKER_PORT=9092 -e ADV_HOST=127.0.0.1 landoop/fast-data-dev
```

**Via Kitematic**

Go to `settings > General` and set the environment variables

|Variable name | Value|
|:------|:-----|
| ZK_PORT | 2181 |
| WEB_PORT | 3030 |
| REGISTRY_PORT | 8081 |
| REST_PORT | 8082 |
| CONNECT_PORT | 8083 |
| BROKER_PORT | 9092 |
| ADV_HOST | 127.0.0.1 |

Visit [http://127.0.0.1:3030](http://127.0.0.1:3030) to get into the `fast-data-dev` environment.

See more about [landoop/fast-data-dev](https://hub.docker.com/r/landoop/fast-data-dev) on this link.

# Run the sample code
To compile and run we are assuming you are in the root of the `stream-grpc`
folder, simply:

Run the Consumer

```sh
$ go run cmd/consumer/main.go
```

Run the Producer

```sh
$ go run cmd/producer/main.go
```

Run the Server
```sh
$ go run cmd/server/main.go
```

Run the client:

```sh
$ go run cmd/client/main.go
```


