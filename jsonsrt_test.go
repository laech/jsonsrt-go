package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestReadWriteToFile(t *testing.T) {
	temp, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := temp.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(temp.Name()); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := temp.WriteString("{ }"); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", "jsonsrt", temp.Name())
	cmd.WaitDelay = time.Second * 5

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	if _, err := temp.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	actual, err := io.ReadAll(temp)
	if err != nil {
		t.Fatal(err)
	}

	expected := "{}\n"
	if string(actual) != expected {
		t.Fatalf(`
expected: %#v
      got: %#v"`, expected, string(actual))
	}
}

func TestReadWriteToStdinStdout(t *testing.T) {
	var stdout strings.Builder

	cmd := exec.Command("go", "run", "jsonsrt")
	cmd.Stdin = strings.NewReader("{ }")
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second * 5

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	expected := "{}\n"
	actual := stdout.String()
	if actual != expected {
		t.Fatalf(`
expected: %#v
      got: %#v"`, expected, actual)
	}
}

func TestCanSortByName(t *testing.T) {
	var stdout strings.Builder

	cmd := exec.Command("go", "run", "jsonsrt", "--sort-by-name")
	cmd.Stdin = strings.NewReader(`{"1":0,"0":0}`)
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second * 5

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	expected := `{
  "0": 0,
  "1": 0
}
`
	actual := stdout.String()
	if actual != expected {
		t.Fatalf(`
expected: %#v
      got: %#v"`, expected, actual)
	}
}

func TestCanSortByValue(t *testing.T) {
	var stdout strings.Builder

	cmd := exec.Command("go", "run", "jsonsrt", "--sort-by-value", "x")
	cmd.Stdin = strings.NewReader(`[{"x":1},{"x":0}]`)
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second * 5

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	expected := `[
  {
    "x": 0
  },
  {
    "x": 1
  }
]
`
	actual := stdout.String()
	if actual != expected {
		t.Fatalf(`
expected: %#v
      got: %#v"`, expected, actual)
	}
}
