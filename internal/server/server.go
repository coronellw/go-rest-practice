package server

import (
	"log"
	"net/http"

	"github.com/coronellw/go-microservices/internal/database"
	"github.com/labstack/echo"
)

type Server interface {
	Start() error
}

type EchoServer struct {
	echo *echo.Echo
	DB   database.DatabaseClient
}

func NewEchoServer(db database.Client) Server {
	server := &EchoServer{
		echo: echo.New(),
		DB:   db,
	}

	server.registerRoutes()
	return server
}

func (s *EchoServer) Start() error {
	if err := s.echo.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server shutdown ocurred %s", err)
		return err
	}
	return nil
}

func (s *EchoServer) registerRoutes() {}
