# K8S

minikube start --nodes=2
minikube status
docker ps
kubectl get nodes
kubectl get pods -A
kubectl get pods
kubectl apply -f deployment.yml
kubectl get pods
kubectl get svc
kubectl get deployments

# stop service
kubectl delete deployments <deployment>

# stop minikube
minikube stop
