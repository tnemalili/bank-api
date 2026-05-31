package messaging

import "github.com/sftx/bank-api/core/ports"

type MessaginClient struct {}

func (c *MessaginClient) Publish(topic string, event any) error {
	panic("unimplemented")
}


func NewMessagingClient() *MessaginClient {

	return &MessaginClient{}
}


// Get feedback in case we are not implementing the interface correctly
var _ ports.IMessagingService = (*MessaginClient)(nil)