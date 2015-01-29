package gomua

import (
	"fmt"
)

// MessageThread is a linked list of Messages
type MessageThread struct {
	head *ThreadNode
}

// A ThreadNode is one element of the MessageThread linked list.
type ThreadNode struct {
	msg  *Message
	next *ThreadNode
}

// Take a slice of gomua.Messages and sort them into threads
// A thread starts as a linked-list of gomua.Messages

// Display a threaded conversation
func (t *MessageThread) String() string {
	node := t.head
	var output string
	if node.next == nil {
		output = fmt.Sprintf("From: %v\n", t.head.msg.Header.Get("From")) +
			fmt.Sprintf("To: %v\n", t.head.msg.Header.Get("To")) +
			fmt.Sprintf("Date: %v\n", t.head.msg.Header.Get("Date")) +
			fmt.Sprintf("Subject: %v\n", t.head.msg.Header.Get("Subject")) +
			fmt.Sprintf("\n%s\n", t.head.msg.Content())
	} else {
		for node.next != nil {
			output += fmt.Sprintf("From: %v\n", t.head.msg.Header.Get("From")) +
				fmt.Sprintf("To: %v\n", t.head.msg.Header.Get("To")) +
				fmt.Sprintf("Date: %v\n", t.head.msg.Header.Get("Date")) +
				fmt.Sprintf("Subject: %v\n", t.head.msg.Header.Get("Subject")) +
				fmt.Sprintf("\n%s\n|-\t", t.head.msg.Content())
			node = node.next
		}
	}
	return output
}

// Summary gets a summarized subject for this thread (i.e. remove Re: re: Re:)
func (t *MessageThread) Summary() string {
	node := t.head
	var output string
	if node.next == nil {
		subject := t.head.msg.Header.Get("Subject")
		from := t.head.msg.Header.Get("From")
		return fmt.Sprintf("%s from %s", color(subject, "31"), color(from, "33"))
	}

	for node.next != nil {
		subject := node.msg.Header.Get("Subject")
		from := node.msg.Header.Get("From")
		output += fmt.Sprintf("\t%s from %s", color(subject, "31"), color(from, "33"))
		node = node.next
	}
	return output
}

func (t *MessageThread) appendNode(n *ThreadNode) {
	node := t.head
	for node.next != nil {
		node = node.next
	}
	node.next = n
}

// Thread takes a Mail slice and returns a map of ??? to Threads.
func Thread(msgs []Mail) map[string]*MessageThread {
	threads := make(map[string]*MessageThread)
	// take a slice of mails
	// make a hash table keyed off of subject for now

	// thread them!
	for _, m := range msgs {
		// check the subject -- if we've seen it, put it in the list
		node := new(ThreadNode)
		node.msg = m.(*Message)
		if threads[m.Summary()] == nil {
			thread := new(MessageThread)
			thread.head = node
			threads[m.Summary()] = thread
		} else {
			existingThread := threads[m.Summary()]
			existingThread.appendNode(node)
		}
	}
	return threads
}
