package interactors

import (
	"log"
	"net/http"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/domain"
	tokenM "github.com/Dedushka-Lenin/DockerControllerGolang/internal/domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepo interface {
	Create(login string, password string) (int, error)
	Get(login string) (string, error)
	Check(login string) (bool, error)
}

type Token interface {
	Generate(login string) (string, error)
	Invalidation(token string) error

	GetToken(c *gin.Context) (string, error)
	GetLogin(token string) (string, error)

	Set(c *gin.Context, login, token string) error
}

type UsersHandlers struct {
	cfg *config.Config
	ur  UsersRepo
	tkn Token
}

func NewUsersHandlers(cfg *config.Config, ur UsersRepo, tkn Token) *UsersHandlers {
	return &UsersHandlers{cfg: cfg, ur: ur, tkn: tkn}
}

func (uh *UsersHandlers) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Register. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	if exists, err := uh.ur.Check(user.Login); exists {
		log.Println("Register. Check. err: user with this login already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.RegistrationError})
		return
	} else if err != nil {
		log.Println("Register. Check. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.RegistrationError})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Register. GenerateFromPassword. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.RegistrationError})
		return
	}

	if _, err := uh.ur.Create(user.Login, string(hashedPassword)); err != nil {
		log.Println("Register. Create. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.RegistrationError})
		return
	}

	log.Println("Logout. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (uh *UsersHandlers) Login(c *gin.Context) {
	var userData domain.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		log.Println("Login. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	password, err := uh.ur.Get(userData.Login)
	if err != nil {
		log.Println("Login. Get. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.IncorrectLoginPassword})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(userData.Password)); err != nil {
		log.Println("Login. CompareHashAndPassword. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.IncorrectLoginPassword})
		return
	}

	token, err := uh.tkn.Generate(userData.Login)
	if err != nil {
		log.Println("Login. Generate. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.GenerationError})
		return
	}

	if err = uh.tkn.Set(c, userData.Login, token); err != nil {
		log.Println("Login. Set. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "TODO"})
		return
	}

	log.Println("Logout. err: nil")
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (uh *UsersHandlers) Logout(c *gin.Context) {
	token, err := uh.tkn.GetToken(c)
	if err != nil {
		log.Println("Logout. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	if err := uh.tkn.Invalidation(token); err != nil {
		log.Println("Logout. Invalidation. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidationError})
		return
	}

	log.Println("Logout. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (uh *UsersHandlers) Status(c *gin.Context) {
	if !uh.status(c) {
		c.JSON(http.StatusOK, gin.H{"Status": "не в аккаунте"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "в аккаунте"})
}

func (uh *UsersHandlers) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isLoggedIn := uh.status(c)

		if !isLoggedIn {
			c.Redirect(http.StatusSeeOther, "/register")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (uh *UsersHandlers) status(c *gin.Context) bool {
	token, err := uh.tkn.GetToken(c)
	if err != nil {
		log.Println("status. GetToken. err: " + err.Error())
		return false
	}

	_, err = uh.tkn.GetLogin(token)
	if err != nil {
		log.Println("status. GetLogin. err: " + err.Error())
		return false
	}

	log.Println("status. err: nil")
	return true
}
