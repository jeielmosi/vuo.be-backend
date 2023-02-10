package firestore

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func GetApp(envName string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(os.Getenv(envName + "_FIREBASE_PATH"))
	return firebase.NewApp(context.Background(), nil, opt)
}
