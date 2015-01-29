package gomua_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/frenata/gomua"
)

var msgStr = "From: test1@testing.com\r\nTo: test2@testing.com\r\nDate: Wed, 21 Jan 2015 02:00:03 -0500\r\nSubject: test mail\r\n\r\nTest Content\r\n"

func Test_MessageString(t *testing.T) {
	gomua.Save("_test.msg", msgStr)
	msgs := gomua.Scan("_test.msg")
	m := msgs[0].(*gomua.Message)

	if m.String() != msgStr {
		t.Fatalf("\n%v\n did not parse correctly to \n%v\n", []byte(msgStr), []byte(m.String()))
	}

}

func Test_EmptyMessage(t *testing.T) {
	r := strings.NewReader("")
	m, _ := gomua.ReadMessage(r)
	if m != nil {
		t.Fatal("Message parsed from empty string should return nil.")
	}
}

func Test_MessageFlags(t *testing.T) {
	gomua.Save("_test.msg", msgStr)
	msgs := gomua.Scan("_test.msg")
	m := msgs[0].(*gomua.Message)

	fmt.Println(m.Filename())
	err := m.Move(".")
	if err != nil {
		t.Fatalf("error moving message to '.'")
	}
	if !m.Unread() {
		t.Fatalf("Newly scanned message %s does not show as unread.\n", m.Filename())
	}
	m.Flag(gomua.Seen)
	if !m.IsFlagged(gomua.Seen) {
		t.Fatalf("Flagging then checking for flag fails.")
	}
	m.Flag(gomua.Replied)
	if !m.IsFlagged(gomua.Replied) {
		t.Fatalf("Flagging then checking for flag fails.")
	}
	if !strings.Contains(m.Filename(), "RS") {
		t.Fatalf("flags are not being reordered.")
	}
}
