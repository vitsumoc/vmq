package packets

type PublishConf struct {
	dup                bool
	qos                PUBLISH_FLAG_QOS
	retain             bool
	topic              string
	publishProperties  PROPERTIES
	applicationMessage []byte
}

func NewPublishConf() *PublishConf {
	return &PublishConf{
		qos:                PUBLISH_FLAG_QOS_0,
		publishProperties:  *NewProperties(),
		applicationMessage: make([]byte, 0),
	}
}

func (pc *PublishConf) SetDup(b bool) {
	pc.dup = b
}

func (pc *PublishConf) SetQos(qos PUBLISH_FLAG_QOS) {
	pc.qos = qos
}

func (pc *PublishConf) SetRetain(b bool) {
	pc.retain = b
}

func (pc *PublishConf) SetTopic(topic string) {
	pc.topic = topic
}

func (pc *PublishConf) SetConnectProperties(pp *PROPERTIES) {
	pc.publishProperties = *pp
}

func (pc *PublishConf) SetProperty(key MQTT_PROPERTY_KEY, v1 any, v2 any) error {
	return pc.publishProperties.SetProperty(key, v1, v2)
}

func (pc *PublishConf) SetApplicationMessage(b []byte) {
	pc.applicationMessage = b
}
