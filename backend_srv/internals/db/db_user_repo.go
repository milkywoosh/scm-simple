package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"scm-simple-luke.com/dir/internals/domain"
)

type DBUserRepository struct {
	Conn DBTX
}

func NewUserRepository(db DBTX) domain.UserRepository {
	return &DBUserRepository{
		Conn: db,
	}
}

func (d *DBUserRepository) CreateUser(ctx context.Context, username string, password string, full_name string, email string) error {

	query :=
		`INSERT INTO users 
			(username, hashed_password,full_name,email,password_changed_at,created_at) 
	 	VALUES
	 		($1      , $2             ,$3       ,$4   ,now()              ,now())
		`

	_, err := d.Conn.Exec(
		ctx,
		query,
		username,
		password,
		full_name,
		email,
	)
	if err != nil {
		return err
	}
	return nil
}


func DTOUser(source domain.UserRow) domain.User {
	return domain.User{
		Username:       source.Username.String,
		HashedPassword: source.HashedPassword.String,
		CreatedAt:      &source.CreatedAt.Time,
		Email:          source.Email.String,
		FullName:       source.Fullname.String,
	}
}

func (d *DBUserRepository) GetUser(ctx context.Context, username string) (domain.User, error) {

	var user domain.UserRow

	query := `select username, full_name, email, hashed_password, created_at from users where username = $1`

	row := d.Conn.QueryRow(ctx, query, username)
	if err := row.Scan(
		&user.Username,
		&user.Fullname,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, err
		}
		return domain.User{}, err
	}

	return DTOUser(user), nil

}

func (d *DBUserRepository) CreateVerifyEmail(ctx context.Context, username string, email string, secret_code string) (domain.InfoEmailVerification, error) {
	query := `
		insert into verify_emails (
			username, email, secret_code
		) values (
			$1, $2, $3 
		) returning username, email, secret_code
	`

	var info domain.InfoEmailVerification

	// note created_at and is_verified is auto default
	row := d.Conn.QueryRow(ctx, query, username, email, secret_code)
	err := row.Scan(
		&info.Username,
		&info.Email,
		&info.SecretCode,
	)

	if err != nil {
		return domain.InfoEmailVerification{}, err
	}

	return info, nil
	// return errors.New("belum implement")
}

func (d *DBUserRepository) VerifyEmail(ctx context.Context, email string) error {

	query := `
		update users
			set is_email_verified = true
		where email = $1
	`
	_, err := d.Conn.Exec(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
	// return errors.New("belum implement")

}

type VerifyEmailRow struct {
	Username   pgtype.Text
	SecretCode pgtype.Text
	Email      pgtype.Text
}

func DTOInfoEmailVerif(row VerifyEmailRow) domain.InfoEmailVerification {

	return domain.InfoEmailVerification{
		Username:   row.Username.String,
		SecretCode: row.SecretCode.String,
		Email:      row.Email.String,
	}
}

func (d *DBUserRepository) VerifySecretCode(ctx context.Context, username string, secret_code string) (string, error) {

	query := `
		update verify_emails
			set is_verified = true
		where email = $1 and
			  secret_code = $2
		returning username, email, secret_code

	`
	var verifyData VerifyEmailRow

	row := d.Conn.QueryRow(ctx, query, username, secret_code)
	if err := row.Scan(
		&verifyData.Username,
		&verifyData.Email,
		&verifyData.SecretCode,
	); err != nil {
		return "", err
	}

	return verifyData.Email.String, nil

}



