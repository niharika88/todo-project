package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	g "github.com/todo-project/server/grpc"

	"github.com/gin-gonic/gin"
	"github.com/todo-project/pb"
	"github.com/todo-project/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client

	// Creating Todo Variables
	todoService    services.TodoService
	todoCollection *mongo.Collection
)

func init() {
	config, err := LoadConfig("cmd")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	// Try to find if mongodb uri is supplied from the env.
	mongoDBURL := os.Getenv("MONGODB")
	if len(mongoDBURL) == 0 {
		fmt.Println("Mongodb url not found in env, picking value from config file:", config.DBUri)
		mongoDBURL = config.DBUri
	} else {
		fmt.Println("got mongodb url from env:", mongoDBURL)
	}
	mongoconn := options.Client().ApplyURI(mongoDBURL)
	mongoClient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	//  Instantiate the Constructors
	todoCollection = mongoClient.Database("golang_mongodb").Collection("todos")
	todoService = services.NewTodoService(todoCollection, ctx)

	server = gin.Default()
}

func main() {
	config, err := LoadConfig("cmd")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	startGrpcServer(config)
}

func startGrpcServer(config Config) {

	todoServer, err := g.NewGrpcTodoServer(todoCollection, todoService)
	if err != nil {
		log.Fatal("cannot create grpc todoServer: ", err)
	}

	grpcServer := grpc.NewServer()

	// 👇 Register the Todo gRPC service
	pb.RegisterToDoServiceServer(grpcServer, todoServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}
