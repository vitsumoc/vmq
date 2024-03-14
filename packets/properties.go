package packets

import (
	"errors"
	"io"

	t "github.com/vitsumoc/vmq/types"
)

type PROPERTIES struct {
	PropertyLength *t.MQTT_VAR_INT              // length of Properties + UserProperties, not self
	Properties     map[*t.MQTT_BYTE]t.MQTT_TYPE // any types of MQTT_*, depending on key byte
	UserProperties []*t.MQTT_UTF8_PAIR          // The User Property is allowed to appear multiple times
}

func NewProperties() *PROPERTIES {
	return &PROPERTIES{
		PropertyLength: t.NewVarInt(),
		Properties:     make(map[*t.MQTT_BYTE]t.MQTT_TYPE, 0),
		UserProperties: make([]*t.MQTT_UTF8_PAIR, 0),
	}
}

func (p *PROPERTIES) SetProperty(key MQTT_PROPERTY_KEY, v1 any, v2 any) error {
	pk := t.NewByte().FromValue(byte(key))
	if key == PROPERTY_PAYLOAD_FORMAT_INDICATOR ||
		key == PROPERTY_REQUEST_PROBLEM_INFORMATION ||
		key == PROPERTY_REQUEST_RESPONSE_INFORMATION ||
		key == PROPERTY_MAXIMUM_QOS ||
		key == PROPERTY_RETAIN_AVAILABLE ||
		key == PROPERTY_WILDCARD_SUBSCRIPTION_AVAILABLE ||
		key == PROPERTY_SUBSCRIPTION_IDENTIFIER_AVAILABLE ||
		key == PROPERTY_SHARED_SUBSCRIPTION_AVAILABLE {
		v, ok := v1.(byte)
		if !ok {
			return errors.New("SetProperty value error")
		}
		p.Properties[pk] = t.NewByte().FromValue(v)
	} else if key == PROPERTY_SERVER_KEEP_ALIVE ||
		key == PROPERTY_RECEIVE_MAXIMUM ||
		key == PROPERTY_TOPIC_ALIAS_MAXIMUM ||
		key == PROPERTY_TOPIC_ALIAS {
		v, ok := v1.(uint16)
		if !ok {
			return errors.New("SetProperty value error")
		}
		p.Properties[pk] = t.NewU16().FromValue(v)
	} else if key == PROPERTY_MESSAGE_EXPIRY_INTERVAL ||
		key == PROPERTY_SESSION_EXPIRY_INTERVAL ||
		key == PROPERTY_WILL_DELAY_INTERVAL ||
		key == PROPERTY_MAXIMUM_PACKET_SIZE {
		v, ok := v1.(uint32)
		if !ok {
			return errors.New("SetProperty value error")
		}
		p.Properties[pk] = t.NewU32().FromValue(v)
	} else if key == PROPERTY_CONTENT_TYPE ||
		key == PROPERTY_RESPONSE_TOPIC ||
		key == PROPERTY_ASSIGNED_CLIENT_IDENTIFIER ||
		key == PROPERTY_AUTHENTICATION_METHOD ||
		key == PROPERTY_RESPONSE_INFORMATION ||
		key == PROPERTY_SERVER_REFERENCE ||
		key == PROPERTY_REASON_STRING {
		v, ok := v1.(string)
		if !ok {
			return errors.New("SetProperty value error")
		}
		pv, err := t.NewUtf8().FromValue(v)
		if err != nil {
			return errors.New("SetProperty value error")
		}
		p.Properties[pk] = pv
	} else if key == PROPERTY_CORRELATION_DATA ||
		key == PROPERTY_AUTHENTICATION_DATA {
		v, ok := v1.([]byte)
		if !ok {
			return errors.New("SetProperty value error")
		}
		pv, err := t.NewBin().FromValue(v)
		if err != nil {
			return errors.New("SetProperty value error")
		}
		p.Properties[pk] = pv
	} else if key == PROPERTY_USER_PROPERTY {
		userKey, ok := v1.(string)
		if !ok {
			return errors.New("SetProperty value error")
		}
		userValue, ok := v2.(string)
		if !ok {
			return errors.New("SetProperty value error")
		}
		pv, err := t.NewUtf8Pair().FromValue(userKey, userValue)
		if err != nil {
			return errors.New("SetProperty value error")
		}
		p.UserProperties = append(p.UserProperties, pv)
	} else {
		return errors.New("SetProperty key error")
	}
	// after change, recal length
	p.calLength()
	return nil
}

func (p *PROPERTIES) Length() int {
	return p.PropertyLength.Length() + p.PropertyLength.ToValue()
}

func (p *PROPERTIES) FromStream(input io.Reader) (int, error) {
	// varlength and length represent the length represented in the stream
	varlength := t.NewVarInt()
	_, err := varlength.FromStream(input)
	if err != nil {
		return 0, err
	}
	length := varlength.ToValue() + varlength.Length()
	// read the data need to consider type
	for p.Length() < length {
		pk := t.NewByte()
		_, err := pk.FromStream(input)
		if err != nil {
			return 0, errors.New("PROPERTIES parse error")
		}
		key := MQTT_PROPERTY_KEY(pk.ToValue())
		if key == PROPERTY_PAYLOAD_FORMAT_INDICATOR ||
			key == PROPERTY_REQUEST_PROBLEM_INFORMATION ||
			key == PROPERTY_REQUEST_RESPONSE_INFORMATION ||
			key == PROPERTY_MAXIMUM_QOS ||
			key == PROPERTY_RETAIN_AVAILABLE ||
			key == PROPERTY_WILDCARD_SUBSCRIPTION_AVAILABLE ||
			key == PROPERTY_SUBSCRIPTION_IDENTIFIER_AVAILABLE ||
			key == PROPERTY_SHARED_SUBSCRIPTION_AVAILABLE {
			b := t.NewByte()
			_, err := b.FromStream(input)
			if err != nil {
				return 0, errors.New("PROPERTIES parse error")
			}
			p.Properties[pk] = b
		} else if key == PROPERTY_SERVER_KEEP_ALIVE ||
			key == PROPERTY_RECEIVE_MAXIMUM ||
			key == PROPERTY_TOPIC_ALIAS_MAXIMUM ||
			key == PROPERTY_TOPIC_ALIAS {
			u16 := t.NewU16()
			_, err := u16.FromStream(input)
			if err != nil {
				return 0, errors.New("PROPERTIES parse error")
			}
			p.Properties[pk] = u16
		} else if key == PROPERTY_MESSAGE_EXPIRY_INTERVAL ||
			key == PROPERTY_SESSION_EXPIRY_INTERVAL ||
			key == PROPERTY_WILL_DELAY_INTERVAL ||
			key == PROPERTY_MAXIMUM_PACKET_SIZE {
			u32 := t.NewU32()
			_, err := u32.FromStream(input)
			if err != nil {
				return 0, errors.New("PROPERTIES parse error")
			}
			p.Properties[pk] = u32
		} else if key == PROPERTY_CONTENT_TYPE ||
			key == PROPERTY_RESPONSE_TOPIC ||
			key == PROPERTY_ASSIGNED_CLIENT_IDENTIFIER ||
			key == PROPERTY_AUTHENTICATION_METHOD ||
			key == PROPERTY_RESPONSE_INFORMATION ||
			key == PROPERTY_SERVER_REFERENCE ||
			key == PROPERTY_REASON_STRING {
			utf8 := t.NewUtf8()
			_, err := utf8.FromStream(input)
			if err != nil {
				return 0, errors.New("PROPERTIES parse error")
			}
			p.Properties[pk] = utf8
		} else if key == PROPERTY_CORRELATION_DATA ||
			key == PROPERTY_AUTHENTICATION_DATA {
			ubin := t.NewBin()
			_, err := ubin.FromStream(input)
			if err != nil {
				return 0, errors.New("PROPERTIES parse error")
			}
			p.Properties[pk] = ubin
		} else if key == PROPERTY_USER_PROPERTY {
			utf8p := t.NewUtf8Pair()
			_, err := utf8p.FromStream(input)
			if err != nil {
				return 0, errors.New("PROPERTIES parse error")
			}
			p.UserProperties = append(p.UserProperties, utf8p)
		} else {
			return 0, errors.New("GetProperty key error")
		}
		p.calLength()
	}
	if p.Length() != length {
		return 0, errors.New("PROPERTIES length error")
	}
	return p.Length(), nil
}

func (p *PROPERTIES) ToStream(output io.Writer) (int, error) {
	_, err := p.PropertyLength.ToStream(output)
	if err != nil {
		return 0, err
	}
	for k, v := range p.Properties {
		_, err = k.ToStream(output)
		if err != nil {
			return 0, err
		}
		_, err = v.ToStream(output)
		if err != nil {
			return 0, err
		}
	}
	for _, v := range p.UserProperties {
		userKey := t.NewByte()
		userKey.FromValue(byte(PROPERTY_USER_PROPERTY))
		_, err = userKey.ToStream(output)
		if err != nil {
			return 0, err
		}
		_, err = v.ToStream(output)
		if err != nil {
			return 0, err
		}
	}
	return p.Length(), nil
}

func (p *PROPERTIES) calLength() {
	length := 0
	for _, v := range p.Properties {
		// key length
		length += 1
		// value length
		length += v.Length()
	}
	// user properties length
	for _, u := range p.UserProperties {
		length += 1
		length += u.Length()
	}
	p.PropertyLength.FromValue(length)
}
