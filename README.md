# Cloud-Native Observability Project
![CI/CD](https://github.com/Amir23156/cloud-native-observability/actions/workflows/ci-cd-eks.yml/badge.svg)
[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-Automated%20Build-blue)](https://hub.docker.com/r/<your-dockerhub-username>/go-orders-api)
![Go Version](https://img.shields.io/github/go-mod/go-version/Amir23156/cloud-native-observability)
![Kubernetes](https://img.shields.io/badge/Kubernetes-EKS-blue)

## Overview
This project demonstrates a **cloud-native observability stack** running on **Amazon EKS (Elastic Kubernetes Service)**.  
It includes:
- A **Go microservice (`go-orders-api`)** instrumented with:
  - **Prometheus metrics** for monitoring
  - **OpenTelemetry tracing** for distributed tracing
- **Kubernetes manifests** to deploy the application and observability stack
- **GitHub Actions CI/CD pipelines**:
  - Continuous Integration (CI):
    - Builds and scans Docker images
    - Pushes them to **Docker Hub**
  - Continuous Deployment (CD):
    - Automatically deploys the latest version to **Amazon EKS**
    - Uses GitHub OIDC for **secretless AWS authentication**
- **Grafana**, **Loki**, and **Tempo** for a full observability experience:
  - Metrics → **Prometheus**
  - Logs → **Loki**
  - Traces → **Tempo**

---

## Architecture

```text
Developer -> GitHub -> GitHub Actions CI/CD
         -> Docker Hub -> EKS
                           -> Grafana
                           -> Loki
                           -> Tempo
                           -> Prometheus
```
### Flow:
1. Developer pushes code to GitHub.
2. GitHub Actions CI pipeline builds and scans the Docker image.
3. The image is pushed to Docker Hub with tags:
   - `latest`
   - `sha-<commit>`
4. GitHub Actions CD pipeline deploys the image to EKS via `kubectl`.
5. The application exposes:
   - `/health` endpoint for health checks
   - `/metrics` endpoint for Prometheus scraping
   - Traces are sent to the **OpenTelemetry Collector**, then to **Tempo**.

## Repository structure 

```bash
.
├── app/                  # Go microservice
│   ├── cmd/server        # Entry point for the API
│   ├── internal          # Internal packages
│   │   ├── api           # HTTP routes & handlers
│   │   ├── metrics       # Prometheus metrics
│   │   └── observability # OpenTelemetry tracing
│   ├── go.mod
│   └── go.sum
│
├── k8s/                  # Kubernetes manifests
│   ├── namespace.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   └── service-monitor.yaml
│
├── .github/workflows/
│   └── dockerhub.yml     # GitHub Actions workflow
│
└── helm-values/          # Helm chart values
    ├── kube-prometheus-stack.yaml
    ├── loki-stack.yaml
    ├── tempo.yaml
    └── otel-collector.yaml
```
## Local Development
### Clone the repository
```bash
git clone https://github.com/Amir23156/cloud-native-observability.git
cd cloud-native-observability
```
### Run the Go app locally
```bash
cd app
go run ./cmd/server
```
The app will start on port 5000. You should now test endpoints :
```bash
curl http://localhost:5000/health
curl http://localhost:5000/orders
curl http://localhost:5000/metrics
```
### Build and Push Docker Image
```bash
docker build -t go-orders-api:latest ./app #Build The Image
docker run -p 5000:5000 go-orders-api:latest #Run it locally

docker tag go-orders-api:latest <DOCKERHUB_USERNAME>/go-orders-api:latest
docker push <DOCKERHUB_USERNAME>/go-orders-api:latest #Push the Image to DockerHub
```
## Kubernetes Depployment

### Create namespace

Creates a separate Kubernetes namespace called observability to isolate your application's resources, This Keeps your app's resources (Pods, Services, etc.) organized and Avoids conflicts with other workloads in the cluster.
```bash
kubectl apply -f k8s/namespace.yaml
```
### Deploy The app

Deploy all manifests (Deployments, Services, ConfigMaps, etc.) inside the k8s/ directory, Kubernetes reads all .yaml files inside the folder and Creates Pods, Services, and other components defined for your app.
```bash
kubectl apply -f k8s/
```
### Verify

Confirms that the deployment was successful.
```bash
kubectl get pods -n observability #Shows the Pods running in the observability namespace.
kubectl get svc -n observability  #Lists Services, showing how your application is exposed inside or outside the cluster.
```
### Port-forward to test
Allows you to locally access your Kubernetes service without exposing it externally.
kubectl port-forward connects your local port 8080 to the service port 80 so that any request to http://localhost:8080 is tunneled directly to the go-orders-api app inside the cluster
```bash
kubectl -n observability port-forward svc/go-orders-api 8080:80    
curl http://localhost:8080/health
```

## GitHub Actions CI/CD

**Workflow file:** `.github/workflows/dockerhub.yml`

The pipeline performs:

### 1. CI Stage:
- **Build** Docker image
- **Scan** with Trivy
- **Push** to Docker Hub

### 2. CD Stage:
- **Assume** AWS role via OIDC
- **Deploy** latest image to EKS
- **Rollout update** with zero downtime


