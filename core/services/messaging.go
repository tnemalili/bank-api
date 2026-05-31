// Package services
package services

import (
	"github.com/sftx/bank-api/core/ports"
)

type MessagingService struct {
	msgClient ports.IMessagingService
}

// Publish implements [ports.IMessagingService].
func (m *MessagingService) Publish(topic string, event any) error {
	
	return m.msgClient.Publish(topic, event)
}

// Get feedback in case we are not implementing the interface correctly
var _ ports.IMessagingService = (*MessagingService)(nil)

func NewMessagingService(msgClient ports.IMessagingService) *MessagingService {

	return &MessagingService{msgClient: msgClient}
}
