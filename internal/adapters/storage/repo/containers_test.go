package repo_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/storage/repo"
	"github.com/stretchr/testify/assert"
)

func TestContainersRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryCreate)).
			WithArgs("user_login", "cont_123", "my_cont").
			WillReturnRows(rows)

		id, err := repo.Create("user_login", "cont_123", "my_cont")

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryCreate)).
			WithArgs("user_login", "cont_123", "my_cont").
			WillReturnError(errors.New("db error"))

		id, err := repo.Create("user_login", "cont_123", "my_cont")

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestContainersRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)

		mock.ExpectExec(regexp.QuoteMeta(storage.ContainersQueryDelete)).
			WithArgs("user_login").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = repo.Delete("user_login")

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)

		mock.ExpectExec(regexp.QuoteMeta(storage.ContainersQueryDelete)).
			WithArgs("user_login").
			WillReturnError(errors.New("db error"))

		err = repo.Delete("user_login")

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestContainersRepo_Check(t *testing.T) {
	t.Run("exists_true", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryCheck)).
			WithArgs("id", "user_login").
			WillReturnRows(rows)

		exists, err := repo.Check("id", "user_login")

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryCheck)).
			WithArgs("id", "user_login").
			WillReturnError(errors.New("db error"))

		exists, err := repo.Check("id", "user_login")

		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestContainersRepo_GetList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("1", "first").
			AddRow("2", "second")

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryGetList)).
			WithArgs("user_login").
			WillReturnRows(rows)

		list, err := repo.GetList("user_login")

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, "1", list[0].Id)
		assert.Equal(t, "second", list[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)
		rows := sqlmock.NewRows([]string{"id"}).AddRow("only_id_no_name")

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryGetList)).
			WithArgs("user_login").
			WillReturnRows(rows)

		list, err := repo.GetList("user_login")

		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestContainersRepo_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("10", "test_name")

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryGet)).
			WithArgs(10).
			WillReturnRows(rows)

		res, err := repo.GetById("id")

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "10", res.Id)
		assert.Equal(t, "test_name", res.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := repo.NewContainersRepo(db)

		mock.ExpectQuery(regexp.QuoteMeta(storage.ContainersQueryGet)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		res, err := repo.GetById("id")

		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
