package kafka

type ProducerMock struct {
	PublishMessageMock func(topic, message string) error
}

func (p *ProducerMock) PublishMessage(topic, message string) error {
	return p.PublishMessageMock(topic, message)
}
