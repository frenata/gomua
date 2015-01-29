package gomua

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// A mailDir defines a mail directory and a list of messages in that directory.
type mailDir struct {
	dir  string
	msgs []*Message
}

// newMailDir returns a new mailDir ready to read a given directory.
func newMailDir(dir string) *mailDir {
	msgs := make([]*Message, 0)
	return &mailDir{dir: dir, msgs: msgs}
}

// processFile reads from a given file and creates a new mail Message.
func (md *mailDir) processFile(filename string, in io.Reader) error {
	if in == nil {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		in = bytes.NewReader(b)
	}

	msg, err := ReadMessage(in)
	msg.filename, _ = filepath.Abs(filename)
	if err != nil {
		return err
	}

	md.msgs = append(md.msgs, msg)
	return nil
}

// visitFile processes the file if there is no error.
func (md *mailDir) visitFile(path string, f os.FileInfo, err error) error {
	if err != nil {
		// This was printing errors when reaading a directory, nerfing the msg for now.
		fmt.Printf("%s\n", err)
	}
	err = md.processFile(path, nil)

	return nil
}

// walkDir walks the directory structure and processes each file it finds.
func (md *mailDir) walkDir(path string) {
	filepath.Walk(path, md.visitFile)
}

// Scan checks all files for the given path and returns a slice of Messages.
func Scan(path string) []*Message {
	md := newMailDir(path)

	switch dir, err := os.Stat(md.dir); {
	case err != nil:
		fmt.Printf("%s", err)
	case dir.IsDir():
		md.walkDir(md.dir)
	default:
		if err := md.processFile(md.dir, nil); err != nil {
			fmt.Printf("%s", err)
		}
	}
	return md.msgs
}
