package gomua

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetMessages(t *testing.T) {
	const dir string = "./cmd/mua/testmaildir"
	msgs := Scan(dir)
	threads := Thread(msgs)
	for _, m := range threads {
		fmt.Printf("%s\n", m.Summary())
	}

	if len(threads) != 34 {
		t.Error("Incorrect number of threads")
	}

	msg := threads["Re: Hannover BSD meetup"]
	if strings.Contains(msg.Summary(), "Hannover BSD meetup") == false {
		t.Error("Could not find thread.")
	}
}
