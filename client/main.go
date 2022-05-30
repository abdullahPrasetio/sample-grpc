package main

import (
	"context"
	"fmt"
	config "golang-grpc/common/configs"
	"golang-grpc/common/models"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func serviceUser() models.UsersClient {
	port := config.SERVICE_USER_PORT
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return models.NewUsersClient(conn)
}

func main() {
	// user1 := models.User{
	// 	Age:       "20",
	// 	Name:     "Noval Agung",
	// 	Password: "kw8d hl12/3m,a",
	// 	Gender:   models.UserGender(models.UserGender_value["MALE"]),
	// }


	user := serviceUser()

	fmt.Println("\n", "===========> user test")

	// register user1
	// user.Register(context.Background(), &user1)
	users,err:=user.List(context.Background(),new(empty.Empty))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(users)

}