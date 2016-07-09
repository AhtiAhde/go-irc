package ip

import (
    "testing"
    "strconv"
    "unicode/utf8"
)

func TestGetIP(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Length Test", "13"},
	}
	for _, c := range cases {
		got := GetIP()
		length, _ := strconv.ParseInt(c.want, 10, 0)
		if utf8.RuneCountInString(got) != int(length) {
			t.Errorf("GetIP() returns a string of length %q", c.want)
		}
	}
}