package auth

import (
	"context"
	"log"

	"blinders/packages/utils"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseManager struct {
	App    *firebase.App
	Client *auth.Client
}

func (m FirebaseManager) Verify(jwt string) (*UserAuth, error) {
	authToken, err := m.Client.VerifyIDToken(context.Background(), jwt)
	if err != nil {
		return nil, err
	}

	firebaseUID := authToken.UID
	email := authToken.Claims["email"].(string)
	name := authToken.Claims["name"].(string)

	userAuth := UserAuth{
		Email:  email,
		Name:   name,
		AuthID: firebaseUID,
	}

	return &userAuth, nil
}

func NewFirebaseManager(adminConfig []byte) (*FirebaseManager, error) {
	manager := FirebaseManager{}
	opt := option.WithCredentialsJSON(adminConfig)
	newApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	manager.App = newApp

	newClient, err := manager.App.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	manager.Client = newClient

	return &manager, nil
}

// return FirebaseManager channel instead
func InitFirebaseManagerFromFile(filename string) chan Manager {
	ch := make(chan Manager)
	go func() {
		adminConfig, err := utils.GetFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		m, err := NewFirebaseManager(adminConfig)
		if err != nil {
			log.Fatalf("can not init firebase manager: %v", err)
		}

		ch <- m
	}()

	return ch
}
