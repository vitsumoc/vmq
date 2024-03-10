package types

import (
	"io"
)

func NewUtf8Pair() *MQTT_UTF8_PAIR {
	return &MQTT_UTF8_PAIR{}
}

func (mutf8p *MQTT_UTF8_PAIR) FromStream(input io.Reader) (*MQTT_UTF8_PAIR, int, error) {
	// total length
	len := 0
	// get the key
	key, n, err := NewUtf8().FromStream(input)
	if err != nil {
		return nil, n, err
	}
	mutf8p.key = *key
	len += n
	// get the value
	value, n, err := NewUtf8().FromStream(input)
	if err != nil {
		return nil, len + n, err
	}
	mutf8p.value = *value
	len += n
	return mutf8p, len, nil
}

func (mutf8p *MQTT_UTF8_PAIR) ToStream(output io.Writer) (int, error) {
	len := 0
	// key
	n, err := mutf8p.key.ToStream(output)
	if err != nil {
		return n, err
	}
	len += n
	// value
	n, err = mutf8p.value.ToStream(output)
	if err != nil {
		return len + n, err
	}
	len += n
	return len, nil
}

func (mutf8p *MQTT_UTF8_PAIR) FromValue(k string, v string) (*MQTT_UTF8_PAIR, error) {
	key, err := NewUtf8().FromValue(k)
	if err != nil {
		return nil, err
	}
	value, err := NewUtf8().FromValue(v)
	if err != nil {
		return nil, err
	}
	mutf8p.key = *key
	mutf8p.value = *value
	return mutf8p, nil
}

func (mutf8p *MQTT_UTF8_PAIR) ToValue() (k string, v string) {
	k = mutf8p.key.ToValue()
	v = mutf8p.value.ToValue()
	return k, v
}
