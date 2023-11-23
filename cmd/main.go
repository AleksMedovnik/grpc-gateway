package main

import (
	"context"
	"log"
	"time"

	pb "github.com/AleksMedovnik/grpc-server/pkg/auth/proto"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)
 
const (
   ADDRESS = "localhost:50051"
)

type TodoTask struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
 }

func postTodo(m *fiber.Ctx) error{
    conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
 
    if err != nil {
        log.Fatalf("did not connect : %v", err)
    }
  
    defer conn.Close()
  
    c := pb.NewTodoServiceClient(conn)
  
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  
    defer cancel()

    p := new(TodoTask)
    if err := m.BodyParser(p); err != nil {
        return err
    }
    res, err := c.CreateTodo(ctx, &pb.NewTodo{Name: p.Name, Description: p.Description, Done: p.Done})

    if err != nil {
        log.Fatalf("could not create user: %v", err)
    }

    return m.JSON(res)
}
 
func main() {
    app := fiber.New()

    app.Post("/tasks", postTodo)

    log.Fatal(app.Listen(":3000"))
 
}

