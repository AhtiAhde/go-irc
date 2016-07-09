package ip

import (
    "testing"
    "strconv"
    "unicode/utf8"
    "strings"
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
		if utf8.RuneCountInString(got) > int(length) {
			t.Errorf("GetIP() returns a string longer (%q) than expected length %q", utf8.RuneCountInString(got), c.want)
		}
		if len(strings.Split(got, ".")) != 4 {
		    t.Errorf("GetIP() does not return four digits separated by dots")
		}
	}
}