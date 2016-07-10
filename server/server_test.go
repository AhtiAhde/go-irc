package main

import (
    "testing"
    "errors"
)

type MockConnection struct {
    MessageBuffer []string
    t *testing.T
}

func (conn *MockConnection) Read(b []byte) (n int, err error) {
    return 1, errors.New("Everything Okay")
}

func (conn *MockConnection) Write(b []byte) (n int, err error) {
    conn.MessageBuffer = append(conn.MessageBuffer, string(b))
    return 1, errors.New("Everything Okay")
}

func (conn *MockConnection) Close() error {
    return errors.New("Everything Okay")
}

func (conn *MockConnection) GetLastMessage() string {
    return conn.MessageBuffer[len(conn.MessageBuffer) - 1]
}

func TestHandleClientJoinRequest(t *testing.T) {
    mockConn := new(MockConnection)
    mockConn.t = t
	cases := []struct {
		in, want string
	}{
		{"JOIN:103.23.231.123:4343\n", "Welcome! Your id is: 0, you address is: 103.23.231.123:4343"},
		{"JOIN:123.123.123.123:12345\n", "Welcome! Your id is: 1, you address is: 123.123.123.123:12345"},
		{"PEOPLE:1\n", "0"},
		{"JOIN:76.34.213.124:5678\n", "Welcome! Your id is: 2, you address is: 76.34.213.124:5678"},
		{"PEOPLE:0\n", "1,2"},
	}
	for _, c := range cases {
	    routeRequest(c.in, mockConn)
	    if mockConn.GetLastMessage() != c.want {
	        t.Errorf("Here with %q, expected %q", mockConn.GetLastMessage(), c.want)
	    }
	}
}