package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/grpc"
	repository3 "auth-service/internal/roles/repository"
	repository2 "auth-service/internal/users/repository"
	"auth-service/internal/users/user_service"
	"auth-service/pkg/config"
	"auth-service/pkg/postgres"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, _ := config.SetupConfiguration()
	pg, _ := postgres.NewPG(conf)
	logger := logrus.New()

	//permRepos := repository.NewPermissionRepository(pg, logger)
	rolesRepos := repository3.NewRoleRepository(pg, logger)
	userRepos := repository2.NewUsersRepository(pg, rolesRepos, logger)

	userService := user_service.NewUserService(userRepos, rolesRepos, logger)
	authServ := auth.NewAuthService(userService, nil, logger)
	server := grpc.NewGRPCServer(authServ)

	l, _ := net.Listen("tcp", ":8080")
	defer l.Close()

	go func() {
		logger.Infof("Server is listening on port: %v", ":8080")
		if err := server.Serve(l); err != nil {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	server.GracefulStop()
	logger.Info("Server Exited Properly")
}
