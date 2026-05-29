package bootstrap

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/config"
)

func InitFirebase() (userAuth *auth.Client) {
	config.InitFirebase()

	userApp, err := config.FirebaseAppUser.Auth(context.Background())
	if err != nil {
		log.Fatal("Failed to init Firebase User:", err)
	}

	return userApp
}
