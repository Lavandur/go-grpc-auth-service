package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/common"
	"auth-service/internal/grpc"
	"auth-service/internal/models"
	repository3 "auth-service/internal/roles/repository"
	repository2 "auth-service/internal/users/repository"
	"auth-service/internal/users/user_service"
	"auth-service/pkg/config"
	"auth-service/pkg/postgres"
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	perm := models.Permission{
		PermissionID: uuid.New().String(),
		Name:         "fds",
		Description: common.LocalizedString{
			"sfdfs": "FSDFSF",
		},
	}

	json.Marshal(common.LocalizedString{
		"sfdfs": "FSDFSF",
	})
	query := goqu.Insert("permissions").Rows(
		perm,
	)
	fmt.Println(query.ToSQL())
	q, _, err := query.ToSQL()
	if err != nil {
		panic(err)
	}
	fmt.Println(q)

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
