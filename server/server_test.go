package main

import (
    "testing"
    "errors"
)

type MockServer struct {
    MessageBuffer []string
    t *testing.T
}

func (conn *MockServer) Read(b []byte) (n int, err error) {
    return 1, errors.New("Everything Okay")
}

func (conn *MockServer) Write(b []byte) (n int, err error) {
    conn.MessageBuffer = append(conn.MessageBuffer, string(b))
    return 1, errors.New("Everything Okay")
}

func (conn *MockServer) Close() error {
    return errors.New("Everything Okay")
}

func (conn *MockServer) GetLastMessage() string {
    return conn.MessageBuffer[len(conn.MessageBuffer) - 1]
}

func TestHandleClientJoinRequest(t *testing.T) {
    mockConn := new(MockServer)
    mockConn.t = t
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
	}
	
	for _, c := range cases {
	    routeRequest(c.in, mockConn)
	    if mockConn.GetLastMessage() != c.want {
	        t.Errorf("Here with %q, expected %q", mockConn.GetLastMessage(), c.want)
	    }
	}
	
	if (clients.MessageQueue[4] != Message{}) {
		t.Errorf("Unexpected amount of messages")
	}
}