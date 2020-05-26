package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	grpcClient "github.com/PedPet/user/pkg/grpc"
	"github.com/PedPet/user/pkg/service"
	"google.golang.org/grpc"
)

func main() {
	grpcAddr := os.Getenv("PORT")
	ctx := context.Background()

	conn, err := grpc.Dial(":"+grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("gRPC dial:", err)
	}
	defer conn.Close()
	userService := grpcClient.NewClient(conn)

	flag.Parse()
	args := flag.Args()
	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "createUser":
		var username, email, password string

		username = args[0]
		email = args[1]
		password = args[2]
		createUser(ctx, userService, username, email, password)

	default:
		log.Fatalln("unknown command")
	}
}

func createUser(ctx context.Context, service service.User, username, email, password string) {
	user, err := service.CreateUser(ctx, username, email, password)
	if err != nil {
		log.Fatalln("err:", err)
	}

	fmt.Println("user:", user)
}
