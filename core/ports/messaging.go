// Package ports
package ports

type IMessagingService interface {
	Publish(topic string, event any) error
}

type IMessagingClient interface {
	Publish(topic string, event any) error
}