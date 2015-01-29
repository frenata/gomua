package gomua

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	newline = "\r\n"
	infotag = ":2,"

	// Maildir flags
	Passed  = "P"
	Replied = "R"
	Seen    = "S"
	Trashed = "T"
	Draft   = "D"
	Flagged = "F"
)

// for sorting strings by bytes, used to reorder file flags correctly.
type byLetter []byte

func (a byLetter) Len() int           { return len(a) }
func (a byLetter) Less(i, j int) bool { return a[i] < a[j] }
func (a byLetter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Mail interface defines mail to be read.
// This can be a threaded list of messages, or a single message.
//
// String() should return a string representing what the user needs to see
// to further interact with the Mail. For a single message, this should be
// simply the message. For a threaded list of messages, this might be a list
// of all messages in the thread, or possibly a just newest unread message?
//
// Summary() should return a short string that can fit on one line, suitable
// for printing as 1 item in a list on a screen.
// For instance, for a single message:
//     Test Message from Frenata <mr.k.frenata@gmail.com>
type Mail interface {
	String() string
	Summary() string
}

// Message implmenets the Mail interface. It holds an embedded mail.Message that it will store the Body
// of on demand.
type Message struct {
	*mail.Message
	isStored bool
	content  string
	filename string
}

// String prints a Message: some basic headers and the Message content.
func (m *Message) String() string {
	var output = fmt.Sprintf("From: %v%s", m.Header.Get("From"), newline) +
		fmt.Sprintf("To: %v%s", m.Header.Get("To"), newline) +
		fmt.Sprintf("Date: %v%s", m.Header.Get("Date"), newline) +
		fmt.Sprintf("Subject: %v%s", m.Header.Get("Subject"), newline)

		//output += fmt.Sprintf("\n%s\n", m.Content)
	output += fmt.Sprintf("%s%s", newline, m.SanitizeContent())

	return output
}

// Summary prints a one line summary of the Message content.
func (m *Message) Summary() string {
	subject := m.Header.Get("Subject")
	from := m.Header.Get("From")
	return fmt.Sprintf("%s from %s", color(subject, "31"), color(from, "33"))
}

// adds ANSI color to text
func color(s string, color string) string {
	return "\033[" + color + "m" + s + "\033[0m"
}

// Filename returns Message's current filename.
func (m *Message) Filename() string { return m.filename }

// Content stores the mesasage content if it hasn't yet been stored, then returns the content.
func (m *Message) Content() string {
	if !m.isStored {
		m.Store()
	}
	return m.content
}

// Store reads from the io.Reader in the embedded mail.Message.Body, then permanently stores this content
// in the Message struct.
func (m *Message) Store() {
	b, err := ioutil.ReadAll(m.Message.Body)
	if err != nil {
		log.Fatal(err)
	}

	m.content = string(b)
	m.isStored = true
}

// ReadMessage reads a Message from the designated Reader, and returns the Message.
func ReadMessage(r io.Reader) (*Message, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}

	return &Message{Message: m}, nil
}

// Move moves the Message to a new directory and appends the Info pre flag, if it has not been previously.
func (m *Message) Move(newpath string) error {
	name := filepath.Base(m.Filename())
	newname := filepath.Join(newpath, name) + infotag
	err := os.Rename(m.filename, newname)
	if err != nil {
		return err
	}
	m.filename = newname
	return nil
}

// Flag sets a filename flag on the message.
func (m *Message) Flag(flag string) {
	s := strings.Split(m.filename, infotag)
	if len(s) != 2 {
		log.Fatal(fmt.Errorf("filename %s does not contain '%s'", m.Filename()), infotag)
	}
	name := s[0]
	flags := s[1]

	if !strings.Contains(flags, flag) {
		//if flags[len(flags)-1] != ',' {
		//	flags += ","
		//}
		flags += flag

		bflags := []byte(flags)
		sort.Sort(byLetter(bflags))
		flags = string(bflags)

		newname := name + infotag + flags
		os.Rename(m.filename, newname)
		m.filename = newname
	}
}

// IsFlagged checks if the filename of the message is flagged with the given flag.
func (m *Message) IsFlagged(flag string) bool {
	s := strings.Split(m.filename, infotag)
	if len(s) != 2 {
		log.Fatal(fmt.Errorf("filename %s does not contain '%s'", m.Filename()), infotag)
	}
	flags := s[1]

	return strings.Contains(flags, flag)
}

// Unread returns true if the message is unread.
func (m *Message) Unread() bool {
	return !m.IsFlagged(Seen)
}

// SanitizeContent returns only text/plain content portions of the email.
// --Probably not fully working.--
func (m *Message) SanitizeContent() string {
	var bound string
	var boundB = false
	t := m.Header.Get("Content-Type")
	if strings.Contains(t, "boundary=") {
		bs := strings.Split(t, "boundary=")
		bound = "--" + bs[1]
		boundB = true
	}

	raw := bufio.NewScanner(strings.NewReader(m.Content()))
	buf := new(bytes.Buffer)
	var write = true
	for raw.Scan() {
		line := raw.Text()
		if strings.Contains(line, "Content-Type:") {
			write = false
		}

		if !strings.Contains(line, "Content-Transfer-Encoding") {
			if !boundB || boundB && line != bound {
				if write {
					buf.WriteString(line + newline)
				}
			}
		}

		if strings.Contains(line, "Content-Type: text/plain") {
			write = true
		}
	}
	return string(buf.Bytes())
}

// WriteMessage interactively prompts the user for an email to send.
func WriteMessage(r io.Reader) *mail.Message {
	cli := bufio.NewScanner(r)

	fmt.Print("To: ")
	cli.Scan()
	to := cli.Text()
	fmt.Print("From: ")
	cli.Scan()
	from := cli.Text()
	fmt.Print("Subject: ")
	cli.Scan()
	subject := cli.Text()
	content := WriteContent(r)

	msg := "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += fmt.Sprintf(
		"To: %v\r\nFrom: %v\r\nSubject: %v\r\n\r\n%v",
		to, from, subject, content)

	m, err := mail.ReadMessage(strings.NewReader(msg))
	if err != nil {
		log.Fatal(err)
	}
	return m
}

// WriteContent interactively prompts the user for email content.
func WriteContent(r io.Reader) string {
	cli := bufio.NewScanner(r)
	fmt.Print("Content: (Enter SEND to finish adding content and send the email.\n")

	var content string
	for {
		cli.Scan()
		line := cli.Text()
		if line == "SEND" {
			break
		} else {
			content += line + "\n"
		}
	}
	return content
}

// Save writes a Message to a file.
func Save(file string, m string) error {
	b := bytes.NewBufferString(m).Bytes()
	ioutil.WriteFile(file, b, 0600)
	return nil
}
