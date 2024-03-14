package types

import (
	"io"
)

func NewUtf8Pair() *MQTT_UTF8_PAIR {
	return &MQTT_UTF8_PAIR{}
}

func (mutf8p *MQTT_UTF8_PAIR) FromStream(input io.Reader) (int, error) {
	// get the key
	key := NewUtf8()
	_, err := key.FromStream(input)
	if err != nil {
		return 0, err
	}
	mutf8p.key = *key
	// get the value
	value := NewUtf8()
	_, err = value.FromStream(input)
	if err != nil {
		return 0, err
	}
	mutf8p.value = *value
	return mutf8p.Length(), nil
}

func (mutf8p *MQTT_UTF8_PAIR) ToStream(output io.Writer) (int, error) {
	// key
	_, err := mutf8p.key.ToStream(output)
	if err != nil {
		return 0, err
	}
	// value
	_, err = mutf8p.value.ToStream(output)
	if err != nil {
		return 0, err
	}
	return mutf8p.Length(), nil
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

func (mutf8p *MQTT_UTF8_PAIR) Length() int {
	return mutf8p.key.Length() + mutf8p.value.Length()
}
