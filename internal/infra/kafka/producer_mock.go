package kafka

type ProducerMock struct {
	PublishMessageMock func(topic string, message []byte) error
}

func (p *ProducerMock) PublishMessage(topic string, message []byte) error {
	return p.PublishMessageMock(topic, message)
}
