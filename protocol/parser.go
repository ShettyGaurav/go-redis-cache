package protocol

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func ReadBulkString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\r\n")
	if len(line) < 2 || line[0] != '$' {
		return "", errors.New("invalid bulk string header")
	}
	lengthStr := line[1:]
	length, err := strconv.Atoi(lengthStr)
	if err != nil || length < -1 {
		return "", errors.New("invalid bulk string length")
	}
	if length == -1 {
		return "", nil
	}
	buf := make([]byte, length+2)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return "", err
	}
	if buf[length] != '\r' || buf[length+1] != '\n' {
		return "", errors.New("missing CRLF after bulk String")
	}
	return string(buf[:length]), nil
}

func ParseRESPCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSuffix(line, "\r\n")
	if len(line) < 2 || line[0] != '*' {
		return nil, errors.New("invalid RESP array header")
	}
	count, err := strconv.Atoi(line[1:])
	if err != nil || count < 0 {
		return nil, errors.New("invalid array count")
	}
	var args []string
	for i := 0; i < count; i++ {
		arg, err := ReadBulkString(reader)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}
