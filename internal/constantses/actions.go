package constantses

type Action int

const (
	ActionClientSentPhone Action = iota
	ActionClientSentAppointment
	ActionSendBookingMessage
	ActionSendDeclineMessage
)
