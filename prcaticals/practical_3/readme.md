
---

# Practical 3  – Full-Stack Microservices with gRPC, Databases, and Service Discovery


## 1. Introduction

In this practical, I built a **complete microservices ecosystem** where two independent services (Users and Products) communicate using **gRPC**. Each service has its own **PostgreSQL database** for persistence and uses **Consul** for service discovery.
An **API Gateway** was created as the single entry point. It receives HTTP requests, discovers the services dynamically through Consul, and translates them into gRPC calls. A special **composite endpoint** was also implemented to aggregate data from both services in one response.

---

## 2. System Architecture

* **API Gateway** → Acts as the public-facing entry point. Handles HTTP requests and converts them to gRPC calls.
* **Consul** → Provides service discovery so services don’t need to know each other’s addresses.
* **Users-Service** → Manages user data (name, email). Uses PostgreSQL for storage.
* **Products-Service** → Manages product data (name, price). Uses PostgreSQL for storage.
* **PostgreSQL Databases** → Each service has its own dedicated database.

---

## 3. Steps Completed

### Step 1 – Project Setup

* Installed gRPC and Protobuf tools.
* Defined service contracts in `.proto` files for both Users and Products.
* Generated Go gRPC code from the proto files.
* Created project structure with separate folders for `users-service`, `products-service`, `api-gateway`, and `proto`.

### Step 2 – Orchestration with Docker

* Wrote a `docker-compose.yml` file to run:

  * **Consul** (service discovery)
  * **users-db** and **products-db** (Postgres databases)
  * **users-service** and **products-service** (gRPC services)
  * **api-gateway** (HTTP entry point)

### Step 3 – Users-Service

* Implemented a gRPC server in Go with two methods:

  * `CreateUser`
  * `GetUser`
* Connected to `users-db` using GORM.
* Auto-migrated the User model into PostgreSQL.
* Registered service with Consul.
* Created Dockerfile and containerized the service.

### Step 4 – Products-Service

* Implemented a gRPC server in Go with two methods:

  * `CreateProduct`
  * `GetProduct`
* Connected to `products-db` using GORM.
* Auto-migrated the Product model into PostgreSQL.
* Registered service with Consul.
* Created Dockerfile and containerized the service.

### Step 5 – API Gateway

* Built an HTTP server with Gorilla Mux.
* Connected to both Users and Products services via gRPC.
* Implemented REST endpoints:

  * `POST /api/users` → Create user
  * `GET /api/users/{id}` → Get user
  * `POST /api/products` → Create product
  * `GET /api/products/{id}` → Get product
* Added a **composite endpoint**:

  * `GET /api/purchases/user/{userId}/product/{productId}`
  * Fetches data from both services and returns a combined JSON response.
* Fixed issue where API Gateway was directly calling services by port → Now uses **Consul service discovery**.

---

## 4. Testing

* Started all containers with `docker-compose up --build`.
* Verified in **Consul UI** (`http://localhost:8500`) that services were registered.
* Tested with **cURL/Postman**:

  * Created a user → worked correctly.
  * Retrieved user → returned correct details.
  * Created a product → worked correctly.
  * Retrieved product → returned correct details.
  * Called composite endpoint `/api/purchases/user/1/product/1` → received combined User + Product response.
  ![alt text](<Screenshot 2025-09-07 221736.png>)
![alt text](<Screenshot 2025-09-07 221747.png>)


  ![alt text](<Screenshot 2025-09-07 222010.png>)

---

## 5. Conclusion

This practical helped me understand how to build a **full microservices ecosystem** using Go, gRPC, PostgreSQL, and Consul. I learned how to:

* Implement independent services with their own databases.
* Register and discover services dynamically using Consul.
* Use an API Gateway to translate HTTP requests into gRPC calls.
* Build a composite endpoint that aggregates data from multiple services.

This architecture shows how microservices can remain decoupled but still work together efficiently. 

