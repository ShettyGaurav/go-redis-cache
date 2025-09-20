package protocol

import "strconv"

func EncodeSimpleString(s string) string {
	return "+" + s + "\r\n"
}

func EncodeError(msg string) string {
	return "-" + msg + "\r\n"
}

func EncodeBulkString(s string) string {
	if s == "" {
		return "$-1\r\n"
	}
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

func EncodeInteger(n int) string {
	return ":" + strconv.Itoa(n) + "\r\n"
}
