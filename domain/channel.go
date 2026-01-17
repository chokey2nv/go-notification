package domain

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

type DeliveryStatus string

const (
	DeliveryPending  DeliveryStatus = "pending"
	DeliveryRetrying DeliveryStatus = "retrying"
	DeliverySent     DeliveryStatus = "sent"
	DeliveryFailed   DeliveryStatus = "failed"
)
