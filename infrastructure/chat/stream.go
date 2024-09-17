package chat

import stream "github.com/GetStream/stream-chat-go/v5"

// StreamService struct contains the stream client, and chat service
type StreamService struct {
	StreamClient *stream.Client
	ChatService  ChatServiceInterface
}

// NewStreamService creates a new instance of StreamService with the provided host, port, and password
func NewStreamService(apiKey, apiSecret string) (*StreamService, error) {
	streamClient, _ := stream.NewClient(apiKey, apiSecret)
	return &StreamService{
		StreamClient: streamClient,
		ChatService:  NewChatService(streamClient),
	}, nil
}
