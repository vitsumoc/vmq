package vmq

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	p "github.com/vitsumoc/vmq/packets"
)

// go test -v --count=1 -run TestConnectPacket
func TestConnectPacket(test *testing.T) {
	cc := p.NewConnectConf()
	cc.SetClientID("imvc")
	cc.SetKeepAlive(10)
	err := cc.SetProperty(p.PROPERTY_SESSION_EXPIRY_INTERVAL, uint32(10), nil)
	if err != nil {
		test.Error()
	}

	cp := p.NewConnectPacketFromConf(cc)
	buffer := bytes.NewBuffer(nil)
	n, err := cp.ToStream(buffer)
	if err != nil || n != 24 || buffer.Bytes()[2] != 0x00 {
		test.Error()
	}
}

// go test -v --count=1 -run TestConnDisconnAction
func TestConnDisconnAction(test *testing.T) {
	// vmq
	v := New("test", "tcp", "36.33.24.191", 1883)
	// conn conf
	cc := p.NewConnectConf()
	cc.SetClientID("imcc")
	cc.SetKeepAlive(10)
	err := cc.SetProperty(p.PROPERTY_SESSION_EXPIRY_INTERVAL, uint32(10), nil)
	if err != nil {
		test.Error()
	}
	// conn to server
	err = v.Connect(cc)
	if err != nil {
		test.Error()
	}
	for x := 0; x < 1; x++ {
		time.Sleep(1 * time.Second)
	}
	if v.status != STATUS_CONNECTED {
		test.Error()
	}
	// disconn to server
	dc := p.NewDisconnectConf()
	dc.SetReasonCode(p.RC_ADMINISTRATIVE_ACTION)
	dc.SetProperty(p.PROPERTY_REASON_STRING, "i will", nil)
	err = v.Disconnect(dc)
	if err != nil {
		test.Error()
	}
	for x := 0; x < 1; x++ {
		time.Sleep(1 * time.Second)
	}
	if v.status != STATUS_DISCONNECTED {
		test.Error()
	}
}

// go test -v --count=1 -run TestDisconnByServer
func TestDisconnByServer(test *testing.T) {
	// vmq
	v := New("test", "tcp", "36.33.24.191", 1883)
	// conn conf
	cc := p.NewConnectConf()
	cc.SetClientID("imqq")
	cc.SetKeepAlive(10)
	err := cc.SetProperty(p.PROPERTY_SESSION_EXPIRY_INTERVAL, uint32(10), nil)
	if err != nil {
		test.Error()
	}
	// conn to server
	err = v.Connect(cc)
	if err != nil {
		test.Error()
	}
	// keep alive 10, so when x > 15, server will send disconn (time out)
	for x := 0; x < 20; x++ {
		time.Sleep(1 * time.Second)
		fmt.Println(v.status)
	}
	fmt.Println(v.packetDisconn)
}

// go test -v --count=1 -run TestSub
func TestSub(test *testing.T) {
	// conn sub then check sub and suback
	// vmq
	v := New("test", "tcp", "36.33.24.191", 1883)
	// conn conf
	cc := p.NewConnectConf()
	cc.SetClientID("imsub")
	// conn to server
	err := v.Connect(cc)
	if err != nil {
		test.Error()
	}
	// do sub after 3 second
	for x := 0; x < 20; x++ {
		time.Sleep(1 * time.Second)
		if x == 3 {
			subconf := p.NewSubscribeConf()
			subconf.SetTopic("vctest/vc")
			v.Subscribe(subconf)
		}
	}
	v.Disconnect(p.NewDisconnectConf())
}

// go test -v --count=1 -run TestSubOption
func TestSubOption(test *testing.T) {
	// conn sub then check sub and suback
	// vmq
	v := New("test", "tcp", "36.33.24.191", 1883)
	// conn conf
	cc := p.NewConnectConf()
	cc.SetClientID("imsub")
	// conn to server
	err := v.Connect(cc)
	if err != nil {
		test.Error()
	}
	// do sub after 3 second
	for x := 0; x < 10; x++ {
		time.Sleep(1 * time.Second)
		if x == 3 {
			subconf := p.NewSubscribeConf()
			subconf.SetTopic("vctest/vc")
			subconf.SetTopicAndOption("vctest/vc2", p.SUBSCRIBE_OPTION_QOS_2, true, true, p.SUBSCRIBE_OPTION_RETAIN_HANDLING_1)
			subconf.SetTopic("vctest/vc3")
			v.Subscribe(subconf)
		}
	}
	v.Disconnect(p.NewDisconnectConf())
}

// go test -v --count=1 -run TestPublish
func TestPublish(test *testing.T) {
	// vmq
	v := New("test", "tcp", "36.33.24.191", 1883)
	v.OnMessage(func(v *vmq, topic string, message []byte) {
		fmt.Println(topic)
		fmt.Println(string(message))
	})
	// conn conf
	cc := p.NewConnectConf()
	cc.SetClientID("impub")
	// conn to server
	err := v.Connect(cc)
	if err != nil {
		test.Error()
	}
	// do sub after 3 second
	for x := 0; x < 10; x++ {
		time.Sleep(1 * time.Second)
		if x == 3 {
			subconf := p.NewSubscribeConf()
			subconf.SetTopic("vctest/vc")
			subconf.SetTopicAndOption("vctest/vc2", p.SUBSCRIBE_OPTION_QOS_2, true, true, p.SUBSCRIBE_OPTION_RETAIN_HANDLING_1)
			subconf.SetTopic("vctest/vc3")
			v.Subscribe(subconf)
		}
		if x == 5 {
			pc := p.NewPublishConf()
			pc.SetTopic("vctest/vc")
			pc.SetApplicationMessage([]byte("hello !!!! i ! am ! vmq !!!"))
			v.Publish(pc)
		}
	}
	v.Disconnect(p.NewDisconnectConf())
	time.Sleep(1 * time.Second)
}
