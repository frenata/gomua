package gomua

import (
	"fmt"
	"io"
	"testing"
)

// Test threading

// Get messages

// takes a slice of Messages and prints a numbered list of summaries
func viewMailList(msgs []Mail, w io.Writer) {
	for i, m := range msgs {
		fmt.Fprintf(w, "%d. %s\n", i+1, m.Summary())
	}
}

func TestGetMessages(t *testing.T) {
	const dir string = "./cmd/mua/testmaildir"
	msgs := Scan(dir)
	threads := Thread(msgs)
	for _, m := range threads {
		//fmt.Printf("%s\n", m.Summary())
		_ = m
	}
}
