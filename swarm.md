# Docker swarm


## Build images:
```bash
docker build -f front-end.dockerfile -t tsawler/front-end:tag1 .
docker push tsawler/front-end:tag1
```

## Manage

```bash
docker swarm init
docker swarm join-token worker
docker swarm join-token manager
docker stack deploy -c <stack>.yml <name>
docker service ls
watch docker service ls
docker service scale <name>=<instances>
```

## Updating (pull image and scale first)
```bash
docker service update --image tsawler/listener:1.0.1 myapp_listener-service
```

## Bringing swarm down
Easy method:
```bash
docker stack rm myapp
```
To stop them, scale all services to 0, or just type
```bash
docker swarm leave
```