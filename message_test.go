package gomua

import (
	"strings"
	"testing"
)

var msgStr = "From: test1@testing.com\r\nTo: test2@testing.com\r\nDate: Wed, 21 Jan 2015 02:00:03 -0500\r\nSubject: test mail\r\n\r\nTest Content\r\n"

func Test_MessageString(t *testing.T) {
	r := strings.NewReader(msgStr)
	msgT, _ := ReadMessage(r)

	if msgT.String() != msgStr {
		t.Fatalf("\n%v\n did not parse correctly to \n%v\n", []byte(msgStr), []byte(msgT.String()))
	}

	r = strings.NewReader("")
	msgT, _ = ReadMessage(r)
	if msgT != nil {
		t.Fatal("Message parsed from empty string should return nil.")
	}
}
