
# WEB303 – Practical 4: Kubernetes Microservices with Kong Gateway & Resilience Patterns

## Overview

This practical demonstrates the deployment of a **microservices-based cafe order management system** using **Go**, **React**, **Kong Gateway**, **Consul**, and **Kubernetes**.  
It focuses on **production-grade orchestration**, **service discovery**, **API gateway management**, and **resilience patterns** to build a reliable distributed system.

---

## Objectives

- Containerize microservices and deploy them on Kubernetes.
- Integrate **Kong API Gateway** for intelligent request routing.
- Implement **Consul** for service discovery and inter-service communication.
- Debug and fix issues with microservice communication and order submission.
- Apply **resilience patterns** (Timeout, Retry, Circuit Breaker) to improve reliability.

---

## System Architecture

The **Student Cafe App** consists of three main components:

1. **Frontend (React.js)** – Displays menu items and allows students to place orders.
2. **Food Catalog Service (Go + Chi)** – Provides a list of available food items.
3. **Order Service (Go + Chi)** – Handles order creation and communicates with the catalog service.
4. **Consul** – Manages service discovery between Go microservices.
5. **Kong Gateway** – Acts as an API gateway for frontend-to-backend communication.
6. **Kubernetes (Minikube)** – Hosts and orchestrates all microservices and infrastructure components.

### Data Flow

1. The **React frontend** sends API requests through the **Kong Gateway**.
2. Kong routes requests to the appropriate backend service.
3. The **Order Service** queries the **Food Catalog Service** via Consul for data validation.
4. Responses are returned to the frontend through Kong.

---


## Project Structure
````
student-cafe/
├── food-catalog-service/
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── order-service/
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── cafe-ui/
│   ├── src/
│   ├── public/
│   ├── Dockerfile
│   └── package.json
├── app-deployment.yaml
├── kong-ingress.yaml
└── README.md

````

---

## Tools & Technologies

| Component | Technology Used |
|------------|----------------|
| Frontend | React.js |
| Backend | Go (Chi Router) |
| Service Discovery | Consul |
| API Gateway | Kong Gateway |
| Orchestration | Kubernetes (Minikube) |
| Containerization | Docker |
| Package Manager | Helm |
| Language Runtime | Go v1.23, Node.js v18 |

---

## Deployment & Execution Steps

### 1. Start Kubernetes Cluster

```bash
minikube start --cpus 4 --memory 4096
eval $(minikube -p minikube docker-env)
````

### 2. Create Namespace

```bash
kubectl create namespace student-cafe
```

### 3. Deploy Consul

```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install consul hashicorp/consul --namespace student-cafe --set server.replicas=1
```

### 4. Deploy Kong Gateway

```bash
helm repo add kong https://charts.konghq.com
helm install kong kong/kong --namespace student-cafe
```

### 5. Build and Deploy Microservices

```bash
# Build Docker images (ensure you're in Minikube's Docker environment)
docker build -t food-catalog-service:v1 ./food-catalog-service
docker build -t order-service:v1 ./order-service
docker build -t cafe-ui:v1 ./cafe-ui

# Apply Kubernetes manifests
kubectl apply -f app-deployment.yaml -n student-cafe
kubectl apply -f kong-ingress.yaml -n student-cafe
```

### 6. Access the Application

```bash
minikube service -n student-cafe kong-kong-proxy --url
```

Open the provided URL in browser to view the **Student Cafe UI**.

---

## Debugging & Troubleshooting

### Common Commands

```bash
kubectl get pods -n student-cafe
kubectl get services -n student-cafe
kubectl describe ingress cafe-ingress -n student-cafe
kubectl logs -f deployment/order-deployment -n student-cafe
kubectl logs -f deployment/food-catalog-deployment -n student-cafe
kubectl logs -f deployment/kong-kong -n student-cafe
```

### Test API Endpoints

```bash
curl -X POST http://$(minikube ip):32147/api/orders/orders \
  -H "Content-Type: application/json" \
  -d '{"item_ids": ["1", "2"]}'
```

### Common Fixes

* **Go Version Error:** Update Dockerfile to `golang:1.23-alpine`
* **Images not found:** Run `eval $(minikube -p minikube docker-env)` before building.
* **Kong not accessible:** Check ingress setup and Kong pod logs.
* **Pods pending:** Ensure Minikube has enough resources (`minikube start --cpus 2 --memory 4096`).

---

## Resilience Patterns (Part 2)

### Implemented Patterns

| Pattern             | Description                                                       |
| ------------------- | ----------------------------------------------------------------- |
| **Timeout**         | Prevents service calls from hanging indefinitely.                 |
| **Retry**           | Automatically retries failed requests with backoff.               |
| **Circuit Breaker** | Temporarily halts requests to failing services to allow recovery. |

These patterns were implemented using middleware logic within the Go microservices to ensure improved fault tolerance and reliability.

---

## Screenshots

1. **Frontend Screenshot:** Food menu and successful order placement.
![alt text](<assets/Screenshot 2025-10-03 222422.png>)
2. **Kubernetes Pods:** Output of `kubectl get pods -n student-cafe`.
![alt text](<assets/Screenshot 2025-10-03 222820.png>)
3. **Kong Gateway:** Output of `minikube service -n student-cafe kong-kong-proxy --url`.
![alt text](<assets/Screenshot 2025-10-03 222422.png>)
![alt text](<assets/Screenshot 2025-10-03 222732.png>)
4. **Services Overview:** Output of `kubectl get services -n student-cafe`.
![alt text](<assets/Screenshot 2025-10-03 222956.png>)

---

## Learning Outcomes

* Deployed and managed a **multi-service application** on Kubernetes.
* Integrated **Kong Gateway** for API management and routing.
* Configured **Consul** for dynamic service discovery.
* Identified and resolved inter-service communication issues.
* Implemented **resilience patterns** to improve reliability and fault tolerance.


---

## Conclusion

Through this practical, I successfully deployed and tested a **Kubernetes-based microservices architecture** integrated with **Kong Gateway** and **Consul**.
The exercise strengthened my understanding of **cloud-native deployment**, **microservice orchestration**, and **resilience design patterns** essential for modern distributed applications.

---

