package structure

type IWhatsAppSendQueueRabbitMQ struct {
	Type         string  `json:"type"`
	TargetNumber string  `json:"targetNumber"`
	Message      *string `json:"message,omitempty"`
	FileName     *string `json:"filename,omitempty"`
}
