package gomua

import (
	"strings"
	"testing"
)

func TestGetMessages(t *testing.T) {
	const dir string = "./cmd/mua/testmaildir"
	msgs := Scan(dir)
	threads := Thread(msgs)

	if len(threads) != 26 {
		t.Errorf("Incorrect number of threads, expected %v got %v", 26, len(threads))
	}

	wantSubject := "Hannover BSD meetup"
	msg := threads["<20150122140230.GA18519@autobahn.atexit.net>"]
	if strings.Contains(msg.Summary(), wantSubject) == false {
		t.Errorf("Could not find thread with subject %v, found %v instead.", wantSubject, msg.Summary())
	}
}
