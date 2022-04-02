package padutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/miku/runpad/tempfile"
)

// Runner executes a snippet of code and writes combined output of the program
// to the given writer.
type Runner interface {
	Run(w io.Writer, s *Snippet) error
}

// SimpleFileRunner just executes code on the host, without any security features.
type SimpleFileRunner struct {
	Prefix []string
}

func (r *SimpleFileRunner) Run(w io.Writer, s *Snippet) error {
	var suffix string
	if s.Tag == "go" {
		suffix = ".go"
	}
	f, err := tempfile.TempFile("", "runpad-*", suffix)
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	if _, err := io.WriteString(f, s.Text); err != nil {
		return err
	}
	if len(r.Prefix) == 0 {
		return fmt.Errorf("missing command")
	}
	r.Prefix = append(r.Prefix, f.Name())
	cmd := exec.Command(r.Prefix[0], r.Prefix[1:]...)
	b, err := cmd.CombinedOutput()
	if _, err := io.Copy(w, bytes.NewReader(b)); err != nil {
		return err
	}
	return nil
}

// Text represents a blob of text from a pad. This struct allows to find and
// access code snippets in the text.
type Text struct {
	Content string
}

// Block represents a
type Block struct {
	LineStart int
	LineEnd   int
}

// Snippet contains a types piece of text, e.g. a code block in a given
// language.
type Snippet struct {
	Tag   string
	Text  string
	Block Block
}

// NumLines return number of lines of the snippet.
func (s *Snippet) NumLines() int {
	return s.Block.LineEnd - s.Block.LineStart
}

func (t *Text) Snippets() (snippets []*Snippet) {
	var (
		v       *Snippet
		started bool
		lines   []string
	)
	for i, line := range strings.Split(t.Content, "\n") {
		if !strings.HasPrefix(line, "```") {
			if started {
				lines = append(lines, strings.TrimRight(line, "\n"))
			}
			continue
		}
		if started {
			v.Text = strings.Join(lines, "\n")
			v.Block.LineEnd = i + 1
			snippets = append(snippets, v)
			lines, started, v = nil, false, nil
		} else {
			tag := strings.TrimSpace(line[3:])
			v = &Snippet{
				Tag: tag,
				Block: Block{
					LineStart: i + 1,
				},
			}
			started = true
		}
	}
	return
}
