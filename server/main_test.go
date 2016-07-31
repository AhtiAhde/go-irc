package main

import (
    "testing"
    "errors"
    "strconv"
    "github.com/ThatGuyFromFinland/server/core"
)

type Container struct {
	connection *MockConnection
}

func (this *Container) GetConnection(address core.AddressEntry) core.Handler {
	return this.connection
}

func (this *Container) SetConnection(conn *MockConnection) {
	this.connection = conn
}

type MockConnection struct {
    t *testing.T
	MessageBuffer []string
}

func (this *MockConnection) Read(b []byte) (n int, err error) {
    return 1, errors.New("Everything Okay")
}

func (this *MockConnection) Write(b []byte) (n int, err error) {
    this.MessageBuffer = append(this.MessageBuffer, string(b))
    return 1, errors.New("Everything Okay")
}

func (this *MockConnection) Close() error {
    return errors.New("Everything Okay")
}

func (this *MockConnection) GetLastMessage() string {
    return this.MessageBuffer[len(this.MessageBuffer) - 1]
}

func contactMockClient(address core.AddressEntry) core.Handler {
	return new(MockConnection)
}

func TestRouter(t *testing.T) {
    mockServer := new(MockConnection)
    mockServer.t = t
    mockClient := new(MockConnection)
    container := new(Container)
    container.SetConnection(mockClient)
    
    longJohn := make([]byte, 1048577)
    tooManyRecipients := ""
    for i := 0; i < 300; i++ {
    	tooManyRecipients = tooManyRecipients + strconv.Itoa(i) + ","
    }
	cases := []struct {
		in, want string
	}{
		{"JOIN:103.23.231.123:4343\n", "Welcome! Your id is: 0, you address is: 103.23.231.123:4343"},
		{"JOIN:123.123.123.123:12345\n", "Welcome! Your id is: 1, you address is: 123.123.123.123:12345"},
		{"PEOPLE:1\n", "0"},
		{"JOIN:76.34.213.124:5678\n", "Welcome! Your id is: 2, you address is: 76.34.213.124:5678"},
		{"PEOPLE:0\n", "1,2"},
		{"MESSAGE:0,1:Where do all the aliens hang out?\n", "Sent: \"Where do all the aliens hang out?\" to users 0,1"},
		{"MESSAGE:1,2:I believe they like it at the Foo Bar.\n", "Sent: \"I believe they like it at the Foo Bar.\" to users 1,2"},
		{"MESSAGE:" + tooManyRecipients + ":Should fail for too many users.\n", "Error: Too many recipients!"},
		{"MESSAGE:1,2:" + string(longJohn) + "\n", "Error: Message too long!"},
	}
	
	router.Init()
	go handleMessages(container.GetConnection, &router)
	
	for _, c := range cases {
	    router.RouteRequest(c.in, mockServer)
	    if mockServer.GetLastMessage() != c.want {
	        t.Errorf("Here with %q, expected %q", mockServer.GetLastMessage(), c.want)
	    }
	}
	
	// Check that our messages have arrived, but no more, no less
	if len(mockClient.MessageBuffer) != 4 {
		t.Errorf("Unexpected amount of messages delivered: %, expected 4", len(mockClient.MessageBuffer))
	}
}