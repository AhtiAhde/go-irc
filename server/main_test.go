package main

import (
    "testing"
    "errors"
    "strconv"
    "github.com/ThatGuyFromFinland/server/core"
)

var messageBuffer []string

type MockServer struct {
    t *testing.T
}

func (conn *MockServer) Read(b []byte) (n int, err error) {
    return 1, errors.New("Everything Okay")
}

func (conn *MockServer) Write(b []byte) (n int, err error) {
    messageBuffer = append(messageBuffer, string(b))
    return 1, errors.New("Everything Okay")
}

func (conn *MockServer) Close() error {
    return errors.New("Everything Okay")
}

func getLastMessage() string {
    return messageBuffer[len(messageBuffer) - 1]
}

func contactMockClient(address core.AddressEntry) core.Handler {
	return new(MockServer)
}

func TestHandleClientJoinRequest(t *testing.T) {
    mockConn := new(MockServer)
    mockConn.t = t
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
	
	for _, c := range cases {
	    router.RouteRequest(c.in, mockConn)
	    if getLastMessage() != c.want {
	        t.Errorf("Here with %q, expected %q", getLastMessage(), c.want)
	    }
	}
}