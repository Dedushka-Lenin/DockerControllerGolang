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

func TestUsersRepo_Create(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		password     string
		mockBehavior func(mock sqlmock.Sqlmock, login, password string)
		wantId       int
		wantErr      bool
	}{
		{
			name:     "success",
			login:    "test_user",
			password: "secure_password",
			mockBehavior: func(mock sqlmock.Sqlmock, login, password string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryCreate)).
					WithArgs(login, password).
					WillReturnRows(rows)
			},
			wantId:  1,
			wantErr: false,
		},
		{
			name:     "db_error",
			login:    "test_user",
			password: "secure_password",
			mockBehavior: func(mock sqlmock.Sqlmock, login, password string) {
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryCreate)).
					WithArgs(login, password).
					WillReturnError(errors.New("db error"))
			},
			wantId:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock, tt.login, tt.password)

			r := repo.NewUsersRepo(db)
			id, err := r.Create(tt.login, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantId, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersRepo_Delete(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		mockBehavior func(mock sqlmock.Sqlmock, login string)
		wantErr      bool
	}{
		{
			name:  "success",
			login: "delete_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				mock.ExpectExec(regexp.QuoteMeta(storage.UsersQueryDelete)).
					WithArgs(login).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:  "db_error",
			login: "delete_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				mock.ExpectExec(regexp.QuoteMeta(storage.UsersQueryDelete)).
					WithArgs(login).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock, tt.login)

			r := repo.NewUsersRepo(db)
			err = r.Delete(tt.login)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersRepo_Check(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		mockBehavior func(mock sqlmock.Sqlmock, login string)
		wantExists   bool
		wantErr      bool
	}{
		{
			name:  "exists_true",
			login: "existing_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryCheck)).
					WithArgs(login).
					WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:  "exists_false",
			login: "missing_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryCheck)).
					WithArgs(login).
					WillReturnRows(rows)
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name:  "db_error",
			login: "any_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryCheck)).
					WithArgs(login).
					WillReturnError(errors.New("db error"))
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock, tt.login)

			r := repo.NewUsersRepo(db)
			exists, err := r.Check(tt.login)

			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExists, exists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersRepo_Get(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		mockBehavior func(mock sqlmock.Sqlmock, login string)
		wantPassword string
		wantErr      error
	}{
		{
			name:  "success",
			login: "user1",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				rows := sqlmock.NewRows([]string{"password"}).AddRow("secret_hash")
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryGet)).
					WithArgs(login).
					WillReturnRows(rows)
			},
			wantPassword: "secret_hash",
			wantErr:      nil,
		},
		{
			name:  "not_found",
			login: "unknown_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				mock.ExpectQuery(regexp.QuoteMeta(storage.UsersQueryGet)).
					WithArgs(login).
					WillReturnError(sql.ErrNoRows)
			},
			wantPassword: "",
			wantErr:      sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock, tt.login)

			r := repo.NewUsersRepo(db)
			pass, err := r.Get(tt.login)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, pass)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPassword, pass)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
