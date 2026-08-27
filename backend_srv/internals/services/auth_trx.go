package services

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"scm-simple-luke.com/dir/internals/domain"
)

func (a *AuthenticationService) UserRegistrationTx(ctx context.Context, username string, password string, fullName string, email string, afterCreate func(username string) error) error {
	tx, err := a.uow.BeginAuth(ctx)

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	errCreate := tx.CreateUser(ctx, username, password, fullName, email)
	if errCreate != nil {
		return errCreate
	}

	// off
	/*
		errAfterCreate := afterCreate(username)
		if errAfterCreate != nil {
			return errAfterCreate
		}
	*/

	return tx.Commit(ctx)

}

func (a *AuthenticationService) InfoUser(ctx context.Context, username string) (domain.User, error) {

	auth, err := a.q.AuthQueries(ctx)
	if err != nil {
		return domain.User{}, err
	}

	userInfo, err := auth.GetUser(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return userInfo, err
}

// create row email verification
func (a *AuthenticationService) CreateEmailVerificationTx(ctx context.Context, username string, email string, secret_code string) (domain.InfoEmailVerification, error) {
	tx, err := a.uow.BeginAuth(ctx)
	if err != nil {
		return domain.InfoEmailVerification{}, err
	}

	defer tx.Rollback(ctx)

	infoEmailVerif, err := tx.CreateVerifyEmail(ctx, username, email, secret_code)
	if err != nil {
		return domain.InfoEmailVerification{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.InfoEmailVerification{}, err
	}

	return infoEmailVerif, nil
}

func (a *AuthenticationService) VerifyEmailSecretCodeTx(ctx context.Context, username string, secret_code string) error {

	tx, err := a.uow.BeginAuth(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	email, err := tx.VerifySecretCode(ctx, username, secret_code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("gagal verifikasi secret code, secret code is invalid")
		}
		log.Printf("pass err VerifySecretCode\n")
		return err
	}

	log.Printf("pass VerifySecretCode\n")
	// set users to be verified => true
	err = tx.VerifyEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("gagal verifikasi secret code, email is invalid")
		}
		log.Printf("pass err VerifyEmail\n")
		return err
	}

	log.Printf("pass VerifyEmail\n")

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
