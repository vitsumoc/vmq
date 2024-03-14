package packets

// all member will be private, so SetXXX() will be the only entrance
type ConnectConf struct {
	// connect flags
	cfCleanStart bool
	cfWillFlag   bool
	cfWillQos    int
	cfWillRetain bool
	cfPassword   bool
	cfUsername   bool

	// variable header
	keepAlive  uint16
	properties PROPERTIES

	// payload
	clientID       string
	willProperties PROPERTIES
	willTopic      string
	willPayload    []byte
	username       string
	password       []byte
}

func NewConnectConf() *ConnectConf {
	return &ConnectConf{
		cfCleanStart:   false,
		cfWillFlag:     false,
		cfWillQos:      0,
		cfWillRetain:   false,
		cfPassword:     false,
		cfUsername:     false,
		keepAlive:      0,
		properties:     PROPERTIES{},
		clientID:       "",
		willProperties: PROPERTIES{},
		willTopic:      "",
		willPayload:    []byte{},
		username:       "",
		password:       []byte{},
	}
}

func (cc *ConnectConf) SetCleanStart(b bool) {
	cc.cfCleanStart = b
}

func (cc *ConnectConf) SetWill(qos int, retain bool, properties PROPERTIES, topic string, payload []byte) {
	cc.cfWillFlag = true
	cc.cfWillQos = qos
	cc.cfWillRetain = retain
	cc.willProperties = properties
	cc.willTopic = topic
	cc.willPayload = payload
}

func (cc *ConnectConf) SetPassword(password []byte) {
	cc.cfPassword = true
	cc.password = password
}

func (cc *ConnectConf) SetUsername(username string) {
	cc.cfUsername = true
	cc.username = username
}

func (cc *ConnectConf) SetKeepAlive(keepAlive uint16) {
	cc.keepAlive = keepAlive
}

func (cc *ConnectConf) SetProperties(pp *PROPERTIES) {
	cc.properties = *pp
}

func (cc *ConnectConf) SetClientID(clientID string) {
	cc.clientID = clientID
}
