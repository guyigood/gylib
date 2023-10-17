package mqttpack

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MqttPack struct {
	MqClient *Client
}

type Client struct {
	MQTT.Client
	mqttfun             MQTT.MessageHandler
	UserName, Pass, Url string
	MqttQos             uint8
	MqttRetained        bool
}

//var

func NewMqttPack() *MqttPack {
	this := new(MqttPack)
	this.MqClient = new(Client)
	this.MqClient.MqttQos = 0
	this.MqClient.MqttRetained = false
	return this
}

func (this *MqttPack) SetUrl(url string) *MqttPack {
	this.MqClient.Url = url
	return this
}

func (this *MqttPack) SetUserInfo(u_name, u_pass string) *MqttPack {
	this.MqClient.UserName = u_name
	this.MqClient.Pass = u_pass
	return this
}

func (this *MqttPack) SetMqttInitMQ() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(this.MqClient.Url)
	if this.MqClient.UserName != "" {
		opts.Username = this.MqClient.UserName
	}
	if this.MqClient.Pass != "" {
		opts.Password = this.MqClient.Pass
	}
	//opts.SetAutoReconnect(true)
	this.MqClient.Client = MQTT.NewClient(opts)
}

func (this *MqttPack) SetMsgPayLoad(mqfun func(client MQTT.Client, msg MQTT.Message)) {
	this.MqClient.mqttfun = mqfun
	//func(client MQTT.Client, msg MQTT.Message) {
	//	fmt.Printf("date: %s\n", msg.Payload())
	//}
}

func (this *MqttPack) GetConnectStatus() bool {
	return this.MqClient.IsConnected()
}

/*func SetMqtt(url string) *Client {
	c := &Client{
		Url: url,
	}
	c.initMqtt()
	return c
}

func (m *Client) initMqtt() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.Url)
	m.Client = MQTT.NewClient(opts)
}*/

func (m *Client) Conn() bool {
	token := m.Client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println("mqtt connnect error:", m.Url, token.Error())
		return false
	}
	return true
}

func (m *Client) DisConn() {
	m.Client.Disconnect(250)
}

func (m *Client) Sub(topic string) bool {
	token := m.Client.Subscribe(topic, m.MqttQos, m.mqttfun)
	if token.Wait() && token.Error() != nil {
		fmt.Println("mqtt sub error:", topic, token.Error())
		return false
	}
	return true
}

func (m *Client) Pub(topic, send_str string) bool {
	//QoS0，At most once，至多一次；
	//QoS1，At least once，至少一次；
	//QoS2，Exactly once，确保只有一次。
	//  为什么要设置retained？
	//  1.当消息发布到MQTT服务器时，我们需要保留最新的消息到服务器上，以免订阅时丢失上一次最新的消息；
	//  当订阅消费端服务器重新连接MQTT服务器时，总能拿到该主题最新消息， 这个时候我们需要把retained设置为true;
	// 2.当消息发布到MQ服务器时，我们不需要保留最新的消息到服务器上；
	//  当订阅消费端服务器重新连接MQTT服务器时，不能拿到该主题最新消息，只能拿连接后发布的消息，这个时候我们需要把  retained设置为false;
	token := m.Client.Publish(topic, m.MqttQos, m.MqttRetained, send_str) //只能收到连接后发布的消息，至多一次
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return false
	}
	return true
}

func (m *Client) Pubbyte(topic string, buff []byte) bool {
	token := m.Client.Publish(topic, m.MqttQos, m.MqttRetained, buff)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return false
	}
	return true
}

func test() {
	c := NewMqttPack()
	c.SetUrl("tcp://broker.mqttdashboard.com:1883")
	c.MqClient.Conn()
	c.MqClient.Sub("date")
	c.MqClient.Pub("date", "this is a good")
	time.Sleep(time.Second)
	c.MqClient.DisConn()
	/*c := SetMqtt()
	c.Conn()
	c.Sub()
	c.Pub()
	time.Sleep(time.Second)
	c.DisConn()*/
}
