package gomua

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"strings"
)

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

	Save("test.mail", msg)
	m, err := mail.ReadMessage(strings.NewReader(msg))
	if err != nil {
		log.Fatal(err)
	}
	return m
}

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

func Save(file string, m string) error {
	b := bytes.NewBufferString(m).Bytes()
	ioutil.WriteFile(file, b, 0600)
	return nil
}
