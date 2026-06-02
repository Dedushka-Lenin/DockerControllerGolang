package token

import (
	"fmt"
	"log"
	"time"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenRepo interface {
	Create(login, token string) (int, error)
	Delete(login string) error
	Check(login string) (bool, error)
	GetLogin(token string) (string, error)
}

type Token struct {
	cfg *config.Config
	tr  TokenRepo
}

func NewToken(cfg *config.Config, tr TokenRepo) *Token {
	return &Token{cfg: cfg, tr: tr}
}

func (t *Token) Generate(login string) (string, error) {
	duration := time.Duration(t.cfg.Token.ExpireSeconds) * time.Second

	claims := jwt.MapClaims{
		"user_login": login,
		"exp":        time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err := token.SignedString([]byte(t.cfg.Token.SecretKey))
	return strToken, err
}

func (t *Token) Invalidation(token string) error {
	login, err := t.GetLogin(token)
	if err != nil {
		log.Println("Invalidation. GetLogin. err: " + err.Error())
		return err
	}

	return t.tr.Delete(login)
}

func (t *Token) GetToken(c *gin.Context) (string, error) {
	cook, err := c.Cookie("auth_token")
	return cook, err
}

func (t *Token) GetLogin(token string) (string, error) {
	login, err := t.tr.GetLogin(token)
	if err != nil {
		log.Println("GetLogin. GetLogin. err: " + err.Error())
		return "", err
	}

	if exists, err := t.check(login, token); !exists {
		err = fmt.Errorf("There is no user with this login")
		log.Println("GetLogin. check. err: " + err.Error())
		return "", err
	} else if err != nil {
		log.Println("GetLogin. check. err: " + err.Error())
		return "", err
	}

	log.Println("GetLogin. err: nil")
	return login, nil
}

func (t *Token) Set(c *gin.Context, login, token string) error {
	if exists, err := t.check(login, token); exists {
		err = t.tr.Delete(token)
		log.Println("Set. Delete. err: " + err.Error())
		return err
	} else if err != nil {
		log.Println("Set. Create. err: " + err.Error())
		return err
	}

	if _, err := t.tr.Create(login, token); err != nil {
		log.Println("Set. Create. err: " + err.Error())
		return err
	}

	c.SetCookie(
		"auth_token",
		token,
		t.cfg.Token.ExpireSeconds,
		"/",
		"localhost",
		false, // Переключите в false для HTTP (локально)
		true,
	)

	log.Println("Set. err: nil")
	return nil
}

func (t *Token) check(login, tokenString string) (bool, error) {
	if exists, err := t.tr.Check(login); !exists {
		log.Println("check. Check. err: не найденно")
		return false, err
	} else if err != nil {
		log.Println("check. Check. err: " + err.Error())
		return false, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(t.cfg.Token.SecretKey), nil
	})

	if err != nil || !token.Valid {
		log.Println("check. Parse. err: " + err.Error())
		err = t.tr.Delete(login)
		log.Println("check. Delete. err: " + err.Error())
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, err
	}

	_, ok = claims["user_login"].(string)
	if !ok {
		err = t.tr.Delete(login)
		log.Println("check. Delete. err: " + err.Error())
		return false, err
	}

	log.Println("check. err: nil")
	return true, nil
}
