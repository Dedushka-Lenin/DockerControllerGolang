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

func TestTokenRepo_Create(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		token        string
		mockBehavior func(mock sqlmock.Sqlmock, login, token string)
		wantId       int
		wantErr      bool
	}{
		{
			name:  "success",
			login: "user_login",
			token: "jwt_token_string",
			mockBehavior: func(mock sqlmock.Sqlmock, login, token string) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryCreate)).
					WithArgs(login, token).
					WillReturnRows(rows)
			},
			wantId:  1,
			wantErr: false,
		},
		{
			name:  "db_error",
			login: "user_login",
			token: "jwt_token_string",
			mockBehavior: func(mock sqlmock.Sqlmock, login, token string) {
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryCreate)).
					WithArgs(login, token).
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

			tt.mockBehavior(mock, tt.login, tt.token)

			r := repo.NewTokenRepo(db)
			id, err := r.Create(tt.login, tt.token)

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

func TestTokenRepo_Delete(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		mockBehavior func(mock sqlmock.Sqlmock, login string)
		wantErr      bool
	}{
		{
			name:  "success",
			login: "user_login",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				mock.ExpectExec(regexp.QuoteMeta(storage.TokenQueryDelete)).
					WithArgs(login).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:  "db_error",
			login: "user_login",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				mock.ExpectExec(regexp.QuoteMeta(storage.TokenQueryDelete)).
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

			r := repo.NewTokenRepo(db)
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

func TestTokenRepo_Check(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		mockBehavior func(mock sqlmock.Sqlmock, login string)
		wantExists   bool
		wantErr      bool
	}{
		{
			name:  "exists_true",
			login: "active_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryCheck)).
					WithArgs(login).
					WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:  "exists_false",
			login: "expired_user",
			mockBehavior: func(mock sqlmock.Sqlmock, login string) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryCheck)).
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
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryCheck)).
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

			r := repo.NewTokenRepo(db)
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

func TestTokenRepo_GetLogin(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		mockBehavior func(mock sqlmock.Sqlmock, token string)
		wantLogin    string
		wantErr      error
	}{
		{
			name:  "success",
			token: "valid_token",
			mockBehavior: func(mock sqlmock.Sqlmock, token string) {
				rows := sqlmock.NewRows([]string{"login"}).AddRow("found_user")
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryGet)).
					WithArgs(token).
					WillReturnRows(rows)
			},
			wantLogin: "found_user",
			wantErr:   nil,
		},
		{
			name:  "not_found",
			token: "invalid_or_expired_token",
			mockBehavior: func(mock sqlmock.Sqlmock, token string) {
				mock.ExpectQuery(regexp.QuoteMeta(storage.TokenQueryGet)).
					WithArgs(token).
					WillReturnError(sql.ErrNoRows)
			},
			wantLogin: "",
			wantErr:   sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock, tt.token)

			r := repo.NewTokenRepo(db)
			login, err := r.GetLogin(tt.token)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, login)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLogin, login)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
