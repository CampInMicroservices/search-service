# Search service

## Installation (locally)


1. Install make (`choco install make` for Win)
2. Install golang-migrate package (https://github.com/golang-migrate/migrate)
3. ...

### Run the following commands

```
# Create docker network
make network

# Start database container and create main db
make postgres
make createdb

# Run db migrations
make migrateup

# Start the service
make server
```

Service is avalilable at http://localhost:8080

```
GET  localhost:8080/v1/listings/:id
GET  localhost:8080/v1/listings?offset=0&limit=10
POST localhost:8080/v1/listings
```
### Database migrations

To update schema of the database, run the following command:

```
migrate create -ext sql -dir db/migration -seq <your-schema-name>
```

New file should be created in db/migrations.

## Docker

To run the service in docker, use the following command:
```
# Start
docker-compose up -d

# Stop
docker-compose down
```

## Kubernetes

### Local kubernetes cluster

Install `kubectl` and `minikube` on your local machine (+ `Docker`).

```
# Minikube initialization
minikube start
minikube dashboard

# Create deployment
kubectl apply -f k8s/deploymeny.yml
kubectl get services
kubectl get deployments
kubectl get pods

# Logs
kubectl logs search-service-deployment-b757866d9-dd7px

# Expose the app through a tunnel
minikube service search-service-service
```

### Azure kubernetes cluster

Install `az` CLI. See instructions on https://learn.microsoft.com/en-us/azure/aks/learn/quick-kubernetes-deploy-cli.

```
az login -u <username> -p <password>

# az account list
# az account set --subscription <id>

# az provider register --namespace Microsoft.OperationsManagement
# az provider register --namespace Microsoft.OperationalInsights

# Create a cluster inside resource group
az aks create -g RSO -n RSO-cluster --enable-managed-identity --node-count 1 --generate-ssh-keys

# Create a record in ~/.kube/config to access the cluster
az aks get-credentials --resource-group RSO --name RSO-cluster

# Apply deployment
kubectl apply -f k8s/deploymeny.yml
# kubectl get service search-service-service --watch      # Press ^c after the EXTERNAL-IP shows up
# (service type was changed to ClusterIP since we are using ingress as reverse proxy)

# Apply ingress deployment and ingress nginx controller (only in search-service)
kubectl apply -f k8s/ingress.yml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.5.1/deploy/static/provider/cloud/deploy.yaml

# Logs
kubectl logs search-service-deployment-577c88bcdb-5qjdq

# DELETE cluster
az aks delete --name RSO-cluster --resource-group RSO
```

## How the service was created

```
git clone repo-name
cd repo-name

go mod init service-name
go get -u github.com/gin-gonic/gin

go get github.com/spf13/viper
go get github.com/jmoiron/sqlx
```
