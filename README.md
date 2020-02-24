Example of use gRPC with golang
===================================================

The server and client demonstrate:
 - how to use grpc go libraries.
 - How to create a Kafka Producer and Consumer 

See the definition of the service in transactions/transactions.proto.


## Services architecture

![](/assets/services.png)

--------

# Dependencies

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


