package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// FirebaseService struct
type FirebaseService struct {
	FirebaseApp *firebase.App
}

// NewFirebaseService creates a new instance of FirebaseService
func NewFirebaseService(filename string) (*FirebaseService, error) {
	opt := option.WithCredentialsFile(filename)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return &FirebaseService{FirebaseApp: firebaseApp}, nil
}

func (s *FirebaseService) SendNotification(tokens []string, data map[string]string) error {
	client, err := s.FirebaseApp.Messaging(context.Background())
	if err != nil {
		return err
	}

	// Create the message that will be multicast to multiple tokens.
	message := &messaging.MulticastMessage{
		Data:   data,
		Tokens: tokens,
	}

	// Send a message to the device corresponding to the provided registration token.
	_, err = client.SendMulticast(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
