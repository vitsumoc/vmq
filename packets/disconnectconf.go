package packets

type DisconnectConf struct {
	disconnectReasonCode MQTT_REASON_CODE
	disconnectProperties PROPERTIES
}

func NewDisconnectConf() *DisconnectConf {
	return &DisconnectConf{
		disconnectReasonCode: RC_NORMAL,
		disconnectProperties: *NewProperties(),
	}
}

func (dc *DisconnectConf) SetReasonCode(rc MQTT_REASON_CODE) {
	dc.disconnectReasonCode = rc
}

func (dc *DisconnectConf) SetDisconnectProperties(pp *PROPERTIES) {
	dc.disconnectProperties = *pp
}

func (dc *DisconnectConf) SetProperty(key MQTT_PROPERTY_KEY, v1 any, v2 any) error {
	return dc.disconnectProperties.SetProperty(key, v1, v2)
}
