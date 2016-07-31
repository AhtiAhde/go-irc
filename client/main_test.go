package main

import (
    "testing"
    "github.com/ThatGuyFromFinland/client/core"
)

var messageBuffer []string

func getLastMessage() string {
    return messageBuffer[len(messageBuffer) - 1]
}

func sendMockRequest(input string, serverIp string) string {
    switch {
        case len(input) > 7 && input[:8] == "PEOPLE:1":
            messageBuffer = append(messageBuffer, "0,2")
            return "0,2"
        case len(input) > 6 && input[:7] == "WHOAMI:":
            messageBuffer = append(messageBuffer, "1")
            return "1"
        case len(input) > 9 && input[:10] == "MESSAGE:2:":
            messageBuffer = append(messageBuffer, "To 2: " + input[10:])
            return "To 2: " + input[9:]
        case len(input) > 11 && input[:12] == "MESSAGE:0,2:":
            messageBuffer = append(messageBuffer, "All: " + input[12:])
            return "To 2: " + input[9:]
        default:
            messageBuffer = append(messageBuffer, "Unexpected input: " + input)
            return "Unexpected input: " + input
    }
}

func TestClient(t *testing.T) {
    cases := []struct {
		in, want string
	}{
		{"/WHOIS\n", "0,2"},
		{"/WHOAMI\n", "1"},
		{"/PRIVATE 2 There are many paths to testing\n", "To 2: There are many paths to testing\n"},
		{"Final test to rule them all\n", "All: Final test to rule them all\n"},
	}
	
	client := core.NewClient("1", "123.123.123.123", sendMockRequest)
	
	for _, c := range cases {
	    client.HandleCommands(c.in)
	    if getLastMessage() != c.want {
	        t.Errorf("Here with %q, expected %q", getLastMessage(), c.want)
	    }
	}
}