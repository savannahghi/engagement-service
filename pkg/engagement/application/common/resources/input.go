package resources

// SendSMSPayload is used to serialise an SMS sent through the twilio service REST API
type SendSMSPayload struct {
	To      []string `json:"to"`
	Message string   `json:"message"`
}

// EMailMessage holds data required to send emails
type EMailMessage struct {
	Subject string   `json:"subject,omitempty"`
	Text    string   `json:"text,omitempty"`
	To      []string `json:"to,omitempty"`
}
