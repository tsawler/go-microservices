# Working with Microservices in Go

From the root level of the project, execute this command:

~~~
make up
~~~

Then start the front end:

~~~
make start
~~~


Hit the front end with your web browser at `http://locahost:80`

To stop everything:

~~~
make stop
make down
~~~

All make commands:

~~~
tcs@Grendel go-microservices % make help
 Choose a command:
  up                 Build all projects and start docker compose
  down               Stop docker compose
  start              starts the front end
  stop               stop the front end
  restart_broker     rebuilds and restarts broker-service
  restart_auth       rebuilds and restarts authentication-service
  restart_listener   rebuilds and restarts queue-listener-service
  restart_logger     rebuilds and restarts logger-service
  test               runs all tests
  clean              runs go clean and deletes binaries
  help               displays help
~~~