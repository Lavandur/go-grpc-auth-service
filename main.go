package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/auth/auth_service"
	"auth-service/internal/common"
	"auth-service/internal/delivery"
	repository3 "auth-service/internal/roles/repository"
	"auth-service/internal/roles/role_service"
	"auth-service/internal/users"
	repository2 "auth-service/internal/users/repository"
	"auth-service/internal/users/user_service"
	"auth-service/pkg/config"
	"auth-service/pkg/postgres"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/*perm := models.Permission{
		PermissionID: uuid.New().String(),
		Title:         "fds",
		Description: common.LocalizedString{
			"sfdfs": "FSDFSF",
		},
	}*/

	json.Marshal(common.LocalizedString{
		"sfdfs": "FSDFSF",
	})

	conf, _ := config.SetupConfiguration()
	pg, _ := postgres.NewPG(conf)
	logger := logrus.New()

	//permRepos := repository.NewPermissionRepository(pg, logger)
	rolesRepos := repository3.NewRoleRepository(pg, logger)
	userRepos := repository2.NewUsersRepository(pg, rolesRepos, logger)

	roleService := role_service.NewRoleService(rolesRepos, logger, conf)
	userService := user_service.NewUserService(userRepos, roleService, logger)
	paseto := auth_service.NewPaseto()

	authS := auth_service.NewAuthServiceImpl(paseto, userService, logger)

	authServ := auth.NewAuthService(authS, logger)
	usServ := users.NewUserGrpcService(userService, logger)
	server := delivery.NewGRPCServer(authServ, usServ)

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
