# K8S

```
minikube start --nodes=2

minikube status
docker ps
kubectl get nodes
kubectl get pods -A
kubectl get pods
kubectl apply -f deployment.yml # or kubectl apply -f <folder>
kubectl get pods
kubectl get svc
kubectl get deployments
minikube dashboard
```

# hit the app
```
kubectl expose deployment broker --type=LoadBalancer --port=8080 --target-port=80
minikube tunnel
```

# stop service/app
```
kubectl delete deployments <deployment>
kubectl delete services <service>
```

# stop minikube
```
minikube stop
```
