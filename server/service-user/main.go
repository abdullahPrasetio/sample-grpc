package main

import (
	"context"
	"database/sql"
	config "golang-grpc/common/configs"
	"golang-grpc/common/models"
	"golang-grpc/connections"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)



type UserServer struct {
	db *sql.DB
}

func NewServer(db *sql.DB) (*UserServer) {
	return &UserServer{db:db} 
}

func(s *UserServer)Register(ctx context.Context,param *models.User)(*models.UserWithoutPassword,error){
	password :=[]byte(param.Password)
	hashedPassword,err :=bcrypt.GenerateFromPassword(password,bcrypt.DefaultCost)
	if err != nil{
		panic (err.Error())
	}
	sql :="INSERT INTO users (name,age,password,gender) VALUES (?,?,?,?)"
	res,err:=s.db.ExecContext(ctx,sql,param.Name,param.Age,hashedPassword,param.Gender)
	
	if err!=nil {
		log.Println("Error inserting user",err.Error())
	}
	id, err := res.LastInsertId()
    if err != nil {
        panic (err.Error())
    }
    user:=models.UserWithoutPassword{
		Id:id,
		Name:param.Name,
		Age:param.Age,
		Gender:param.Gender,
	}
	log.Println("User",&user)
	return &user,nil
}

func (s *UserServer)List(ctx context.Context,void *empty.Empty) (*models.UserList, error) {
	sql:="SELECT id,name,age,gender FROM users"
	userList := models.UserList{}
	rows,err:=s.db.QueryContext(ctx,sql)
	if err != nil {
		return &userList,err
	}
	defer rows.Close()
	for rows.Next() {
		user:=models.UserWithoutPassword{}
		rows.Scan(&user.Id,&user.Name,&user.Age,&user.Gender)
		log.Println(&user)
		userList.List = append(userList.List, &user)
	}
	
	return &userList, nil
}
func main() {
	srv := grpc.NewServer()
	db,err:=connections.NewConnection()
	if err!=nil {
		panic(err)
	}
    userSrv :=NewServer(db)
    models.RegisterUsersServer(srv, userSrv)

    log.Println("Starting RPC server at", config.SERVICE_USER_PORT)

  // more code here ...
  l, err := net.Listen("tcp", config.SERVICE_USER_PORT)
	if err != nil {
	log.Fatalf("could not listen to %s: %v", config.SERVICE_USER_PORT, err)
	}

	log.Fatal(srv.Serve(l))
}