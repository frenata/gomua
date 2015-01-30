package gomua_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/frenata/gomua"
)

var msgStr = "From: test1@testing.com\r\nTo: test2@testing.com\r\nDate: Wed, 21 Jan 2015 02:00:03 -0500\r\nSubject: test mail\r\n\r\nTest Content\r\n"

func scanStr(s string) *gomua.Message {
	gomua.Save("_test.msg", s)
	msgs := gomua.Scan("_test.msg")
	m := msgs[0].(*gomua.Message)

	return m
}

func Test_MessageString(t *testing.T) {
	m := scanStr(msgStr)

	if m.String() != msgStr {
		t.Fatalf("\n%v\n did not parse correctly to \n%v\n", []byte(msgStr), []byte(m.String()))
	}

	sum := m.Summary()
	switch {
	case !strings.Contains(sum, "test mail"):
		t.Fatal("Messsage summary not displaying correctly: ", sum)
	case !strings.Contains(sum, " from "):
		t.Fatal("Messsage summary not displaying correctly: ", sum)
	case !strings.Contains(sum, "test1@testing.com"):
		t.Fatal("Messsage summary not displaying correctly: ", sum)
	}
}

func Test_EmptyMessage(t *testing.T) {
	r := strings.NewReader("")
	m, _ := gomua.ReadMessage(r)
	if m != nil {
		t.Fatal("Message parsed from empty string should return nil.")
	}
}

func Test_WriteMessage(t *testing.T) {
	rStr := "From: from\r\nTo: to\r\nSubject: subject\r\n\r\ncontent\r\n"
	wStr := "from\nto\nsubject\n"
	cStr := "content\nSEND\n"
	m := scanStr(rStr)

	r := strings.NewReader(wStr)
	c := strings.NewReader(cStr)
	mScan := gomua.WriteMessage(r, c)
	if m.String() != mScan.String() {
		fmt.Println([]byte(m.String()))
		fmt.Println([]byte(mScan.String()))
		t.Fatal("Interactive input and reading from file to not produce the same Message.")
	}
}

func Test_MessageFlags(t *testing.T) {
	m := scanStr(msgStr)

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
