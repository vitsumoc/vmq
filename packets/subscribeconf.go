package packets

import (
	t "github.com/vitsumoc/vmq/types"
)

type SubscribeConf struct {
	// variable header
	subscribeProperties PROPERTIES
	subPairs            []SUB_PAIR
}

func NewSubscribeConf() *SubscribeConf {
	return &SubscribeConf{
		subscribeProperties: *NewProperties(),
		subPairs:            make([]SUB_PAIR, 0),
	}
}

func (sc *SubscribeConf) SetSubscribeProperties(pp *PROPERTIES) {
	sc.subscribeProperties = *pp
}

func (sc *SubscribeConf) SetProperty(key MQTT_PROPERTY_KEY, v1 any, v2 any) error {
	return sc.subscribeProperties.SetProperty(key, v1, v2)
}

func (sc *SubscribeConf) SetTopic(topicfilter string) error {
	return sc.SetTopicAndOption(topicfilter, SUBSCRIBE_OPTION_QOS_0, false, false, SUBSCRIBE_OPTION_RETAIN_HANDLING_0)
}

func (sc *SubscribeConf) SetTopicAndOption(topicfilter string, qos SUBSCRIBE_OPTION_QOS, noLocal bool,
	retainAsPublished bool, retainHandling SUBSCRIBE_OPTION_RETAIN_HANDLING) error {
	var subOption byte = 0x00
	subTopic, err := t.NewUtf8().FromValue(topicfilter)
	if err != nil {
		return err
	}
	subOption = subOption | byte(qos)
	if noLocal {
		subOption = subOption | SUBSCRIBE_OPTION_NO_LOCAL
	}
	if retainAsPublished {
		subOption = subOption | SUBSCRIBE_OPTION_RETAIN_AS_PUBLISHED
	}
	subOption = subOption | byte(retainHandling)
	sc.subPairs = append(sc.subPairs, SUB_PAIR{
		TopicFilter:         *subTopic,
		SubscriptionOptions: *t.NewByte().FromValue(subOption),
	})
	return nil
}
