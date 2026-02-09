package security

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(keyPath string) (*FirebaseAuth, error) {
	if keyPath == "" {
		return nil, nil
	}

	opt := option.WithCredentialsFile(keyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseAuth{client: client}, nil
}

func (f *FirebaseAuth) VerifyToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if f == nil || f.client == nil {
		return nil, nil
	}
	return f.client.VerifyIDToken(ctx, idToken)
}

func (f *FirebaseAuth) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	if f == nil || f.client == nil {
		return nil, nil
	}
	return f.client.GetUser(ctx, uid)
}
