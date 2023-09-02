package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"

	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/model"
)

var _ Storage = &sqlStorage{}

type sqlStorage struct {
	*sql.DB
}

func (ss sqlStorage) AddUser(ctx context.Context, u *model.User) (*model.User, error) {
	_, err := ss.ExecContext(ctx, sqlInsertUser, u.Login, u.HashedPassword, "", u.DataKey)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return nil, constant.ErrDuplicateRecord
			}
		}
		return nil, err
	}

	return u, nil
}

func (ss sqlStorage) GetUser(ctx context.Context, login string) (*model.User, error) {
	u := model.User{}
	err := ss.QueryRowContext(ctx, sqlSelectUser, login).
		Scan(&u.ID, &u.Login, &u.HashedPassword, &u.RestorePassword, &u.DataKey)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, constant.ErrNotFound
	case err != nil:
		return nil, err
	}

	return &u, nil
}

func (ss sqlStorage) SaveSecret(ctx context.Context, s *model.Secret, login string) (*model.Secret, error) {
	var result sql.Result
	var err error

	if s.ID != 0 {
		result, err = ss.ExecContext(
			ctx,
			sqlUpdateSecret,
			s.Type,
			s.Title,
			s.Login,
			s.Password,
			s.Note,
			s.FileName,
			s.CardholderName,
			s.CardNumber,
			s.Expiration,
			s.SecurityCode,
			s.ID,
			login,
		)
	} else {
		result, err = ss.ExecContext(
			ctx,
			sqlInsertSecret,
			s.Type,
			s.Title,
			s.Login,
			s.Password,
			s.Note,
			s.File,
			s.FileName,
			s.CardholderName,
			s.CardNumber,
			s.Expiration,
			s.SecurityCode,
			login,
		)
	}

	if err != nil {
		return nil, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	s.ID = insertedID

	return s, nil
}

func (ss sqlStorage) GetSecretsByUser(ctx context.Context, login string) (map[int]*model.Secret, error) {
	result := make(map[int]*model.Secret)
	rows, err := ss.QueryContext(ctx, sqlFindSecretsByUser, login)
	if errors.Is(err, sql.ErrNoRows) {
		// This case is valid, though user exists, she doesn't have any secrets
		return result, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var secret model.Secret

		err = rows.Scan(
			&secret.ID,
			&secret.Type,
			&secret.Title,
			&secret.Login,
			&secret.Password,
			&secret.Note,
			&secret.File,
			&secret.FileName,
			&secret.CardholderName,
			&secret.CardNumber,
			&secret.Expiration,
			&secret.SecurityCode,
		)
		if err != nil {
			return nil, err
		}

		result[int(secret.ID)] = &secret
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (ss sqlStorage) DeleteSecret(ctx context.Context, id, login string) error {
	result, err := ss.ExecContext(ctx, sqlDeleteSecret, id, login)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (ss sqlStorage) GetSecret(ctx context.Context, secretID, login string) (*model.Secret, error) {
	secret := model.Secret{}

	err := ss.QueryRowContext(ctx, sqlGetSecretByID, secretID, login).Scan(
		&secret.ID,
		&secret.Type,
		&secret.Title,
		&secret.Login,
		&secret.Password,
		&secret.Note,
		&secret.File,
		&secret.FileName,
		&secret.CardholderName,
		&secret.CardNumber,
		&secret.Expiration,
		&secret.SecurityCode,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, constant.ErrNotFound
	case err != nil:
		return nil, err
	}

	return &secret, nil
}

func (ss sqlStorage) Close() error {
	return ss.DB.Close()
}

func NewSQLStorage(ctx context.Context, dsn string) Storage {
	db, err := sql.Open("sqlite3", "./kms.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(sqlCreateUserTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(sqlCreateSecretTable)
	if err != nil {
		panic(err)
	}

	return sqlStorage{
		DB: db,
	}
}
