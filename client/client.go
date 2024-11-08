package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/oogway93/gRPC_CRUD_APP/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("p", 50020, "Port for listening a server")
)

func main() {
	flag.Parse()
	port := *port
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewCRUDClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Println("\nWelcome to a simple gRPC/PostgreSQL based app that performs CRUD operations!")
	log.Println("Enter the one of the following choices below:")
	log.Print("1 to create an item; 2 to read; 3 to update and 4 to remove: ")

	choice := bufio.NewReader(os.Stdin)
	text, _ := choice.ReadString('\n')

	switch text {
	case "1\n":
		// CreateUser operation
		// Read the Name
		log.Print("\nEnter the Name: ")
		name := bufio.NewReader(os.Stdin)
		n, _ := name.ReadString('\n')
		n = strings.Trim(n, "\n")

		// Read the Age
		log.Print("Enter the Age: ")
		age := bufio.NewReader(os.Stdin)
		a, _ := age.ReadString('\n')
		a = strings.Trim(a, "\n")
		Age, _ := strconv.Atoi(a)

		// Read the Email
		log.Print("Enter the Email: ")
		email := bufio.NewReader(os.Stdin)
		email2, _ := email.ReadString('\n')
		email2 = strings.Trim(email2, "\n")

		// Populate the User struct
		id, err := c.Create(ctx, &pb.UserMessage{
			Name:  n,
			Age:   uint32(Age),
			Email: email2,
		})
		if err != nil {
			log.Fatalf("Could not create a new User: %v", err)
		}
		log.Printf("\nInserted %s age of %d and email: %s | User ID is %d\n", n, Age, email2, id.Id)

	case "2\n":
		// ReadUser operation
		fmt.Print("\nEnter the ID: ")
		id := bufio.NewReader(os.Stdin)
		i, _ := id.ReadString('\n')
		i = strings.Trim(i, "\n")
		ID, _ := strconv.Atoi(i)

		read, err := c.Read(ctx, &pb.ID{Id: uint32(ID)})
		if err != nil {
			log.Fatalf("Error reading the item: %v", err)
		}
		fmt.Println("\nUser found!")
		fmt.Println("ID:", read.Id)
		fmt.Println("Name:", read.Name)
		fmt.Println("Email:", read.Email)
		fmt.Println("Age:", read.Age)

	case "3\n":
		// UpdateUser operation
		// Read the ID
		fmt.Print("\nEnter the existing ID: ")
		id := bufio.NewReader(os.Stdin)
		i, _ := id.ReadString('\n')
		i = strings.Trim(i, "\n")
		ID, _ := strconv.Atoi(i)

		// Read the name
		log.Print("\nEnter the Name: ")
		name := bufio.NewReader(os.Stdin)
		n, _ := name.ReadString('\n')
		n = strings.Trim(n, "\n")

		// Read the Age
		log.Print("Enter the New Age: ")
		age := bufio.NewReader(os.Stdin)
		a, _ := age.ReadString('\n')
		a = strings.Trim(a, "\n")
		Age, _ := strconv.Atoi(a)

		// Read the Email
		log.Print("Enter the New Email: ")
		email := bufio.NewReader(os.Stdin)
		email2, _ := email.ReadString('\n')
		email2 = strings.Trim(email2, "\n")

		up, err := c.Update(ctx, &pb.UserMessage{Name: n, Id: uint32(ID), Age: uint32(Age), Email: email2})
		if err != nil {
			log.Fatalf("Error updating the item: %v", err)
		}
		log.Printf("\nItem updated with the ID: %d", up.Id)

	case "4\n":
		// DeleteUser operation
		// Read the ID
		fmt.Print("\nEnter the existing ID: ")
		id := bufio.NewReader(os.Stdin)
		i, _ := id.ReadString('\n')
		i = strings.Trim(i, "\n")
		ID, _ := strconv.Atoi(i)
		del, _ := c.Delete(ctx, &pb.ID{Id: uint32(ID)})
		if del != nil {
			log.Printf("\nItem with the ID %d deleted", del.Id)
		} else {
			log.Printf("Successful delete (even though ID didn't exist)")
		}

	default:
		log.Println("\nWrong option!")
	}
}
