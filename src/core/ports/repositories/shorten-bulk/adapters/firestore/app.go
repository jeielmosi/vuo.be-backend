package firestore

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func getApp(envName string) (*firebase.App, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv(envName + "_FIREBASE_PATH"))

	return firebase.NewApp(ctx, nil, opt)
}
