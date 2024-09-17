package chat

import (
	"context"
	"fmt"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
)

var (
	ctx = context.Background()
)

// ChatServiceInterface defines the methods that a chat service should implement.
type ChatServiceInterface interface {
	GetStreamToken(userID string) (*TokenDetails, error)
	CreateChannel(channelID string, userID string) error
	AddMember(channelID string, userID string) error
	RemoveMember(channelID string, userID string) error
}

// ChatService represents the chat service implementation
type ChatService struct {
	streamClient *stream.Client
}

// TokenDetails represents the details of a stream token.
type TokenDetails struct {
	Token string
}

var _ ChatServiceInterface = &ChatService{}

// NewChatService creates and returns a new instance of ChatService
func NewChatService(streamClient *stream.Client) *ChatService {
	return &ChatService{streamClient: streamClient}
}

func (s *ChatService) GetStreamToken(userID string) (*TokenDetails, error) {
	token, err := s.streamClient.CreateToken(userID, time.Time{})
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		Token: token,
	}, nil
}

// CreateChannel creates a new Stream Chat channel with the given parameters.
func (s *ChatService) CreateChannel(channelID string, userID string) error {
	resp, err := s.streamClient.CreateChannel(ctx, "messaging", channelID, userID, nil)
	if err != nil {
		fmt.Println("Error creating channel:", err)
		return err
	}

	_, err = resp.Channel.AddMembers(ctx, []string{userID})
	if err != nil {
		fmt.Println("Error adding members:", err)
		return err
	}

	return nil
}

func (s *ChatService) AddMember(channelID string, userID string) error {
	// Check if the user exists
	userExists, err := s.checkUserExists(userID)
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return err
	}

	if !userExists {
		// Handle the case where the user doesn't exist
		// User does not exist, create the user
		err = s.createUser(userID)
		if err != nil {
			fmt.Println("Error creating user:", err)
			return err
		}
	}

	channel := s.streamClient.Channel("messaging", channelID)

	_, err = channel.AddMembers(ctx, []string{userID})
	if err != nil {
		fmt.Println("Error adding members:", err)
		return err
	}

	return nil
}

func (s *ChatService) RemoveMember(channelID string, userID string) error {
	channel := s.streamClient.Channel("messaging", channelID)

	_, err := channel.RemoveMembers(ctx, []string{userID}, nil)
	if err != nil {
		fmt.Println("Error removing members:", err)
		return err
	}

	return nil
}

func (s *ChatService) checkUserExists(userID string) (bool, error) {
	// Create a query option to search for the specific user by their ID
	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"id": userID,
		},
	}

	// Make a request to the GetStream.io API to search for users
	users, err := s.streamClient.QueryUsers(ctx, query)
	if err != nil {
		fmt.Println("Error querying users:", err)
		return false, err
	}

	// Check if the user exists
	if len(users.Users) > 0 {
		// User exists
		return true, nil
	}

	// User does not exist
	return false, nil
}

func (s *ChatService) createUser(userID string) error {
	// Create a new user object
	user := &stream.User{
		ID: userID,
	}

	// Make a request to the GetStream.io API to create or update the user
	_, err := s.streamClient.UpsertUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
