# Working with Microservices in Go

This is the source code for the Udemy course **Working with Microservices and Go**. This project
consists of a number of loosely coupled microservices, all written in Go:

- broker-service: an optional single entry point to connect to all services from one place (accepts JSON)
- authentication-service: authenticates users against a Postgres database (accepts JSON)
- logger-service: logs important events to a MongoDB database (accepts RPC and JSON)
- queue-listener-service: consumes messages from amqp (RabbitMQ) and initiates actions based on payload
- mail-service: sends email (accepts JSON)

In addition to those services, the included `docker-compose.yml` at the root level of the project
starts the following services:

- Postgresql 14
- etcd
- mailhog
- MongoDB

## Running the project
From the root level of the project, execute this command (this assumes that you have 
[GNU make](https://www.gnu.org/software/make/) and a recent version
of [Docker](https://www.docker.com/products/docker-desktop) installed on your machine):

~~~
make up
~~~

Then start the front end:

~~~
make start
~~~


Hit the front end with your web browser at `http://localhost:80`

To stop everything:

~~~
make stop
make down
~~~

All make commands:

~~~
tcs@Grendel go-microservices % make help
 Choose a command:
  up                 starts all containers in the background without forcing build
  up_build           build all projects and start docker compose
  down               stop docker compose
  start              starts the front end
  stop               stop the front end
  restart_broker     rebuilds and restarts broker-service
  restart_auth       rebuilds and restarts authentication-service
  restart_listener   rebuilds and restarts queue-listener-service
  restart_logger     rebuilds and restarts logger-service
  restart_mail       rebuilds and restarts mail-service
  test               runs all tests
  clean              runs go clean and deletes binaries
  help               displays help
~~~