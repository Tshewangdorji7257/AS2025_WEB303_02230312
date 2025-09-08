// api-gateway/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "practical-three/proto/gen/proto"
)

var consulClient *api.Client

// A struct to hold the aggregated data
type UserPurchaseData struct {
	User    *pb.User    `json:"user"`
	Product *pb.Product `json:"product"`
}

func main() {
	// Initialize Consul client
	config := api.DefaultConfig()
	config.Address = "consul:8500"  // Use consul service name from docker-compose
	var err error
	consulClient, err = api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to Consul: %v", err)
	}

	r := mux.NewRouter()
	// User routes
	r.HandleFunc("/api/users", createUserHandler).Methods("POST")
	r.HandleFunc("/api/users/{id}", getUserHandler).Methods("GET")
	// Product routes
	r.HandleFunc("/api/products", createProductHandler).Methods("POST")
	r.HandleFunc("/api/products/{id}", getProductHandler).Methods("GET")

	// The combined endpoint to get aggregated data
	r.HandleFunc("/api/purchases/user/{userId}/product/{productId}", getPurchaseDataHandler).Methods("GET")

	log.Println("API Gateway listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

// Function to discover service address from Consul
func discoverService(serviceName string) (string, error) {
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances of %s found", serviceName)
	}

	// Use the first healthy instance
	service := services[0]
	address := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	return address, nil
}

// Function to get gRPC connection to users service
func getUsersServiceConnection() (pb.UserServiceClient, *grpc.ClientConn, error) {
	address, err := discoverService("users-service")
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return client, conn, nil
}

// Function to get gRPC connection to products service
func getProductsServiceConnection() (pb.ProductServiceClient, *grpc.ClientConn, error) {
	address, err := discoverService("products-service")
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewProductServiceClient(conn)
	return client, conn, nil
}

// User Handlers
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := getUsersServiceConnection()
	if err != nil {
		http.Error(w, fmt.Sprintf("Service discovery failed: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer conn.Close()

	var req pb.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)
	
	res, err := client.CreateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.User)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := getUsersServiceConnection()
	if err != nil {
		http.Error(w, fmt.Sprintf("Service discovery failed: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer conn.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	
	res, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.User)
}

// Product Handlers
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := getProductsServiceConnection()
	if err != nil {
		http.Error(w, fmt.Sprintf("Service discovery failed: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer conn.Close()

	var req pb.CreateProductRequest
	json.NewDecoder(r.Body).Decode(&req)
	
	res, err := client.CreateProduct(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Product)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := getProductsServiceConnection()
	if err != nil {
		http.Error(w, fmt.Sprintf("Service discovery failed: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer conn.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	
	res, err := client.GetProduct(context.Background(), &pb.GetProductRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Product)
}

// Fixed handler for combined data - properly aggregating from both services
func getPurchaseDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	productId := vars["productId"]

	var wg sync.WaitGroup
	var user *pb.User
	var product *pb.Product
	var userErr, productErr error

	wg.Add(2)

	// Concurrently fetch user data
	go func() {
		defer wg.Done()
		client, conn, err := getUsersServiceConnection()
		if err != nil {
			userErr = err
			return
		}
		defer conn.Close()

		res, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: userId})
		if err != nil {
			userErr = err
			return
		}
		user = res.User
	}()

	// Concurrently fetch product data
	go func() {
		defer wg.Done()
		client, conn, err := getProductsServiceConnection()
		if err != nil {
			productErr = err
			return
		}
		defer conn.Close()

		res, err := client.GetProduct(context.Background(), &pb.GetProductRequest{Id: productId})
		if err != nil {
			productErr = err
			return
		}
		product = res.Product
	}()

	wg.Wait()

	if userErr != nil {
		http.Error(w, fmt.Sprintf("Failed to get user data: %v", userErr), http.StatusNotFound)
		return
	}
	
	if productErr != nil {
		http.Error(w, fmt.Sprintf("Failed to get product data: %v", productErr), http.StatusNotFound)
		return
	}

	purchaseData := UserPurchaseData{
		User:    user,
		Product: product,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchaseData)
}