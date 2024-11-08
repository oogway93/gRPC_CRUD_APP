package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/oogway93/gRPC_CRUD_APP/config"
	database "github.com/oogway93/gRPC_CRUD_APP/db"
	pb "github.com/oogway93/gRPC_CRUD_APP/proto"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	port = flag.Int("p", 50020, "Port for listening a server")
)

type Server struct {
	pb.UnimplementedCRUDServer
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}
func main() {
	flag.Parse()
	config, err := config.NewCfg()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	db := database.DatabaseConnection(database.Config{
		Username: config.User,
		Password: config.Password,
		Host:     config.Host,
		Port:     config.Port,
		DBName:   config.Name,
	},
	)
	db.AutoMigrate(&database.User{})
	log.Println("Successfully migrated the database")
	port := *port
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error starting and listening the TCP server: %v", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterCRUDServer(s, NewServer(db))
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		os.Exit(1)
	}
	log.Printf("Hosting server on: %s", listen.Addr().String())
}

func Validator(inUser *pb.UserMessage) error {
	if inUser.Age <= 0 || inUser.Email == "" || inUser.Name == "" {
		return fmt.Errorf("error Validating Input Data, retry again")
	}
	return nil
}

func (s *Server) Create(ctx context.Context, inUser *pb.UserMessage) (*pb.ID, error) {
	if err := Validator(inUser); err != nil {
		return nil, err
	}
	dataModel := &database.User{
		Name:  inUser.Name,
		Age:   inUser.Age,
		Email: inUser.Email,
	}
	tx := s.db.Begin()
	result := tx.Create(&dataModel)
	if result.Error != nil {
		log.Fatalf("Error saving INPUT USER Data: %v", result.Error)
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &pb.ID{Id: uint32(dataModel.ID)}, nil
}
func (s *Server) Delete(ctx context.Context, inID *pb.ID) (*pb.ID, error) {
	var userModel database.User
	tx := s.db.Begin()
	result := tx.Where("id = ?", inID.Id).Delete(&userModel)
	if result.Error != nil {
		log.Fatalf("Error deleting USER BY CURRENT ID: %v", result.Error)
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &pb.ID{Id: uint32(userModel.ID)}, nil
}

func (s *Server) Read(ctx context.Context, inID *pb.ID) (*pb.UserMessage, error) {
	var userModel database.User
	result := s.db.Where("id = ?", inID.Id).First(&userModel)
	if result.Error != nil {
		log.Fatalf("Error reading USER BY ID: %v", result.Error)
		return nil, result.Error
	}
	userResponse := &pb.UserMessage{
		Id:    uint32(userModel.ID),
		Name:  userModel.Name,
		Email: userModel.Email,
		Age:   userModel.Age,
	}
	return userResponse, nil
}
func (s *Server) Update(ctx context.Context, inUser *pb.UserMessage) (*pb.ID, error) {
	if err := Validator(inUser); err != nil {
		return nil, err
	}
	userModel := &database.User{}
	tx := s.db.Begin()
	checkUser := tx.Where("id = ?", inUser.Id).First(&userModel)
	if checkUser.Error != nil {
		log.Fatalf("Error checking THE USER BY ID: %v", checkUser.Error)
		tx.Rollback()
		return nil, checkUser.Error
	}
	userModel.Name = inUser.Name
	userModel.Age = inUser.Age
	userModel.Email = inUser.Email
	userModel.CreatedAt = time.Now()
	result := tx.Save(&userModel)
	if result.Error != nil {
		log.Fatalf("Error saving THE NEW DATA OF USER: %v", checkUser.Error)
		tx.Rollback()
		return nil, checkUser.Error
	}
	tx.Commit()
	return &pb.ID{Id: uint32(userModel.ID)}, nil
}
