# Docker swarm


## Build images:
```bash
docker build -f front-end.dockerfile -t front-end:tag1 .
docker push tsawler/front-end:tag1
```

## Manage

```bash
docker swarm init
docker stack deploy -c <stack>.yml <name>
docker stack rm <name>
docker service ls
docker service scale <name>=<instances>
```

## Updating (scale first)
```bash
docker service update --image tsawler/listener:1.0.1 myapp_listener-service
 ```