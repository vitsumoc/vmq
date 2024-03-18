package vmq

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	p "github.com/vitsumoc/vmq/packets"
)

// go test -v --count=1
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
	// Manually kick out on the server side
	for x := 0; x < 20; x++ {
		time.Sleep(1 * time.Second)
		fmt.Println(v.status)
	}
	fmt.Println(v.packetDisconn)
}
