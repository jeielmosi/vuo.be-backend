package firestore

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func GetApp() (*firebase.App, error) {
	//TODO: Define FIREBASE_PATH
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_PATH"))
	return firebase.NewApp(context.Background(), nil, opt)
}
