package genie

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestLamp_TraverseCommands(t *testing.T) {
	t.Run("validate traverse", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := `lamp
lamp sub1
lamp sub1 sub1-1
lamp sub2
`
		subject := NewLamp("lamp", "", true)
		subject.RootCommand.SubCommands = []*Command{
			{
				Name: "sub1",
				SubCommands: []*Command{
					{
						Name: "sub1-1",
					},
				},
			},
			{
				Name: "sub2",
			},
		}

		subject.TraverseCommands(func(c *Command) {
			b.Write([]byte(c.Path() + "\n"))
		})

		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("want %s, got %s", want, string(got))
		}
	})
}
