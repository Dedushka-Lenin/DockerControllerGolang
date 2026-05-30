package storage

const (
	TokenQueryCreate = `INSERT INTO token (token, login) VALUES ($1, $2) RETURNING id;`
	TokenQueryDelete = `DELETE FROM token WHERE login = $1;`
	TokenQueryCheck  = `SELECT EXISTS (SELECT 1 FROM token WHERE login = $1)`
	TokenQueryGet    = `SELECT login FROM token WHERE token = $1`
)

const (
	UsersQueryCreate = `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	UsersQueryDelete = `DELETE FROM users WHERE login = $1`
	UsersQueryCheck  = `SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)`
	UsersQueryGet    = `SELECT password FROM users WHERE login = $1`
)

const (
	ContainersQueryCreate  = `INSERT INTO containers (login, container_name, container_id) VALUES ($1, $2, $3) RETURNING id;`
	ContainersQueryDelete  = `DELETE FROM containers WHERE id = $1`
	ContainersQueryCheck   = `SELECT EXISTS (SELECT 1 FROM containers WHERE id = $1 AND login = $2)`
	ContainersQueryGetList = `SELECT container_id, container_name FROM containers WHERE login = $1`
	ContainersQueryGet     = `SELECT containers FROM users WHERE id = $1`
)
