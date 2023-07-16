package worker

type ClientMock struct {
	// db *gorm.DB
}

type AsynqClientMock struct {
	client *ClientMock
}

var _ WorkerInterface = &AsynqClientMock{}

func NewMailerMock(client *ClientMock) *AsynqClientMock {
	return &AsynqClientMock{client: client}
}

// Enqueue task to send email
func (ac *AsynqClientMock) SendEmail(payload *Payload) error {
	return nil
}

func (ac *AsynqClientMock) CreatePayload(templateName EmailTemplate, to, subject string, data map[string]interface{}) *Payload {
	return nil
}

func CreateAsynqClientMock() *ClientMock {
	// Create a new Asynq client.
	client := &ClientMock{}
	return client
}
