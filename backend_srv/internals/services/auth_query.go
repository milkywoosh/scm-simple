package services

import (
	"context"

	"scm-simple-luke.com/dir/internals/domain"
)

func (a *AuthenticationService) UserRegistration(ctx context.Context, username string, password string) error {
	_, err := a.q.AuthQueries(ctx)
	if err != nil {
		return err
	}

	// err = auth.CreateUser(ctx, username, password)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (a *AuthenticationService) GetInfoUser(ctx context.Context, username string) (domain.User, error) {

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
