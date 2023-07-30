package worker

type ClientMock struct {
	// pool *pgxpool.Pool
}

type AsynqClientMock struct {
	client *ClientMock
}

var _ MailerInterface = &AsynqClientMock{}

func CreateMailerMock(client *ClientMock) {
	Mailer = &AsynqClientMock{client: client}
}

// Enqueue task to send email
func (ac *AsynqClientMock) SendEmail(payload *EmailPayload) error {
	return nil
}

func (ac *AsynqClientMock) CreateEmailPayload(templateName EmailTemplate, to, subject string, data map[string]interface{}) *EmailPayload {
	return nil
}

func CreateMailerClientMock() *ClientMock {
	// Create a new Asynq client.
	client := &ClientMock{}
	return client
}
