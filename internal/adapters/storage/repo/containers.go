package repo

import (
	"database/sql"
	"log"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/domain"
)

type ContainersRepo struct {
	db *sql.DB
}

func NewContainersRepo(db *sql.DB) *ContainersRepo {
	return &ContainersRepo{db}
}

func (tr *ContainersRepo) Create(login, container_id, container_name string) (int, error) {
	id := 0
	err := tr.db.QueryRow(storage.ContainersQueryCreate, login, container_name, container_id).Scan(&id)
	if err != nil {
		log.Println("Create. QueryRow. err: " + err.Error())
		return 0, err
	}

	log.Println("Check. err: nil")
	return id, nil
}

func (tr *ContainersRepo) Delete(id string) error {
	_, err := tr.db.Exec(storage.ContainersQueryDelete, id)
	if err != nil {
		log.Println("Delete. Exec. err: " + err.Error())
		return err
	}

	log.Println("Check. err: nil")
	return nil
}

func (tr *ContainersRepo) Check(login string, id string) (bool, error) {
	var exists bool

	err := tr.db.QueryRow(storage.ContainersQueryCheck, id, login).Scan(&exists)
	if err != nil {
		log.Println("Check. QueryRow. err: " + err.Error())
		return false, err
	}

	log.Println("Check. err: nil")
	return exists, nil
}

func (tr *ContainersRepo) GetList(login string) ([]domain.Container, error) {
	rows, err := tr.db.Query(storage.ContainersQueryGetList, login)
	if err != nil {
		log.Println("GetList. Query. err: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var containers []domain.Container
	for rows.Next() {
		var c domain.Container
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}

	return containers, rows.Err()
}

func (tr *ContainersRepo) GetById(id string) (*domain.Container, error) {
	var container domain.Container
	err := tr.db.QueryRow(storage.ContainersQueryGet, id).Scan(&container.Id, &container.Name)
	if err != nil {
		log.Println("GetById. QueryRow. err: " + err.Error())
		return nil, err
	}

	log.Println("GetById. err: nil")
	return &container, nil
}
