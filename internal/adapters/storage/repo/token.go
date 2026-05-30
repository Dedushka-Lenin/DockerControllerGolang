package repo

import (
	"database/sql"
	"log"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{db}
}

func (tr *TokenRepo) Create(login string, token string) (int, error) {
	id := 0
	err := tr.db.QueryRow(storage.TokenQueryCreate, token, login).Scan(&id)
	if err != nil {
		log.Println("Create. QueryRow. err: " + err.Error())
		return 0, err
	}

	log.Println("Create. err: nil")
	return id, nil
}

func (tr *TokenRepo) Delete(login string) error {
	_, err := tr.db.Exec(storage.TokenQueryDelete, login)
	if err != nil {
		log.Println("Delete. Exec. err: " + err.Error())
		return err
	}

	log.Println("Delete. err: nil")
	return nil
}

func (tr *TokenRepo) Check(login string) (bool, error) {
	var exists bool
	err := tr.db.QueryRow(storage.TokenQueryCheck, login).Scan(&exists)
	if err != nil {
		log.Println("Check. QueryRow. err: " + err.Error())
		return false, err
	}

	log.Println("Check. err: nil")
	return exists, nil
}

func (tr *TokenRepo) GetLogin(token string) (string, error) {
	var login string
	err := tr.db.QueryRow(storage.TokenQueryGet, token).Scan(&login)
	if err != nil {
		log.Println("GetLogin. QueryRow. err: " + err.Error())
		return "", err
	}

	log.Println("GetLogin. err: nil")
	return login, nil
}
