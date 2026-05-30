package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/containers"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage/repo"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/token"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/interactors"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// router.Use(CORSMiddleware())

type Server struct {
	cfg    *config.Config
	db     *sql.DB
	router *gin.Engine
}

func NewServer(cfg *config.Config) (*Server, error) {
	ctx := context.Background()

	db, err := storage.GetDB(cfg)
	if err != nil {
		return nil, err
	}

	ur := repo.NewUsersRepo(db)
	tr := repo.NewTokenRepo(db)
	cr := repo.NewContainersRepo(db)

	tkn := token.NewToken(cfg, tr)
	ctnr := containers.NewContainers(ctx, cfg, cr)

	users := interactors.NewUsersHandlers(cfg, ur, tkn)
	containers := interactors.NewContainersHandlers(cfg, ctnr, tkn)

	router := gin.Default()
	router.Static("/static", "./frontend")
	router.StaticFile("/favicon.ico", "./frontend/favicon.ico")

	router.GET("/status", users.Status)

	router.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"error": nil})
	})

	router.GET("/register", func(c *gin.Context) {
		c.File("./frontend/account/register.html")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("./frontend/account/login.html")
	})

	usersGroup := router.Group("")
	{
		usersGroup.POST("/registration/", users.Register)
		usersGroup.POST("/account-login/", users.Login)
	}

	usersProtectedGroup := usersGroup.Group("")
	usersProtectedGroup.Use(users.AuthRequired())
	usersProtectedGroup.POST("/logout", users.Logout)

	protectedGroup := router.Group("")
	protectedGroup.Use(users.AuthRequired())
	protectedGroup.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	containersGroup := protectedGroup.Group("containers")
	{
		containersGroup.POST("/create/", containers.Create)
		containersGroup.DELETE("/delete/:id", containers.Delete)

		containersGroup.GET("/status/:id", containers.GetStatus)
		containersGroup.GET("/get", containers.GetList)

		containersGroup.GET("/logs/:id", containers.Logs)

		containersGroup.POST("/start/:id", containers.Start)
		containersGroup.POST("/stop/:id", containers.Stop)
		containersGroup.POST("/restart/:id", containers.Restart)
	}

	return &Server{
		db:     db,
		cfg:    cfg,
		router: router,
	}, nil
}

func (s Server) Run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)

	defer func() {
		stop()
		fmt.Println("db close")
		s.db.Close()
	}()

	go func() {
		s.router.Run(s.cfg.API.Url)
	}()

	<-ctx.Done()
	return nil
}
