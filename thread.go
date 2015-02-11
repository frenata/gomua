package gomua

import (
	"fmt"
	"strings"
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

//  Display a single message
func (t *MessageThread) String() string {
	node := t.head
	var output string
	output = fmt.Sprintf("From: %v\n", node.msg.Header.Get("From")) +
		fmt.Sprintf("To: %v\n", node.msg.Header.Get("To")) +
		fmt.Sprintf("Date: %v\n", node.msg.Header.Get("Date")) +
		fmt.Sprintf("Subject: %v\n", node.msg.Header.Get("Subject")) +
		fmt.Sprintf("\n%s\n", node.msg.Content())

	return output
}

// Summary returns a subject - from summary of the messages in this thread
func (t *MessageThread) Summary() string {
	if t.head == nil {
		return "No message."
	}

	node := t.head
	var output string
	subject := node.msg.Header.Get("Subject")
	from := node.msg.Header.Get("From")
	output += fmt.Sprintf("%s from %s", color(subject, "31"), color(from, "33"))

	for node.next != nil {
		node = node.next
		subject := node.msg.Header.Get("Subject")
		from := node.msg.Header.Get("From")
		output += fmt.Sprintf("\n\t%s from %s", color(subject, "31"), color(from, "33"))
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

// Thread takes a Mail slice and returns a map of subjects to Threads.
func Thread(msgs []Mail) map[string]*MessageThread {
	threads := make(map[string]*MessageThread)

	for _, m := range msgs {
		node := new(ThreadNode)
		node.msg = m.(*Message)
		refs := node.msg.Header.Get("References")
		refSlice := strings.Split(refs, " ")
		parent := refSlice[0]
		if len(parent) == 0 {
			parent = node.msg.Header.Get("Message-Id")
		}
		if threads[parent] == nil {
			thread := new(MessageThread)
			thread.head = node
			threads[parent] = thread
		} else {
			existingThread := threads[parent]
			existingThread.appendNode(node)
		}
	}
	return threads
}
