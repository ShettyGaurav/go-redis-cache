package test

import (
	"bufio"
	"fmt"
	"myredis/protocol"
	"strings"
	"testing"
)

func TestReadBulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	got, err := protocol.ReadBulkString(reader)
	want := "hello"

	fmt.Println("TestReadBulkString output:", got, "error:", err)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseRESPCommand(t *testing.T) {
	input := "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	got, err := protocol.ParseRESPCommand(reader)
	want := []string{"ECHO", "hello"}

	fmt.Println("TestParseRESPCommand output:", got, "error:", err)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d elements, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("at index %d got %q want %q", i, got[i], want[i])
		}
	}
}

func TestEncodeSimpleString(t *testing.T) {
	input := "OK"
	want := "+OK\r\n"
	got := protocol.EncodeSimpleString(input)

	fmt.Println("TestEncodeSimpleString output:", got)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeError(t *testing.T) {
	input := "ERR something went wrong"
	want := "-ERR something went wrong\r\n"
	got := protocol.EncodeError(input)

	fmt.Println("TestEncodeError output:", got)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeBulkString(t *testing.T) {
	input := "hello"
	want := "$5\r\nhello\r\n"
	got := protocol.EncodeBulkString(input)

	fmt.Println("TestEncodeBulkString output:", got)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncodeInteger(t *testing.T) {
	input := 123
	want := ":123\r\n"
	got := protocol.EncodeInteger(input)

	fmt.Println("TestEncodeInteger output:", got)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
