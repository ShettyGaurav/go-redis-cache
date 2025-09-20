package server

import (
	"bufio"
	"fmt"
	"myredis/protocol"
	"myredis/store"
	"net"
	"strings"
)

type Server struct {
	store *store.RedisStore
}

func NewServer(store *store.RedisStore) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%w", addr, err)
	}
	defer listener.Close()
	fmt.Printf("GoCache Server is Running on %s\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept Error %v\v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		args, err := protocol.ParseRESPCommand(reader)
		if err != nil {
			fmt.Printf("Client sent Invalid Data or Disconnected:%v\n", err)
			return
		}
		if len(args) == 0 {
			continue
		}
		response := s.executeCommand(args)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("Failed to write response:%v\n", err)
			return
		}
	}
}

func (s *Server) executeCommand(args []string) string {
	if len(args) == 0 {
		return protocol.EncodeError("ERR empty command")
	}
	cmd := strings.ToUpper(args[0])
	switch cmd {
	case "PING":
		return protocol.EncodeSimpleString("PONG")
	case "GET":
		if len(args) != 2 {
			return protocol.EncodeError("ERR worng number of Inputs")
		}
		key := args[1]
		value, found := s.store.Get(key)
		if !found {
			return protocol.EncodeBulkString("")
		}
		return protocol.EncodeBulkString(value)
	case "SET":
		if len(args) != 3 {
			return protocol.EncodeError("ERR worng number of Inputs")
		}
		key := args[1]
		value := args[2]
		s.store.Set(key, value, 0)
		return protocol.EncodeSimpleString("OK")

	case "DEL":
		if len(args) < 2 {
			return protocol.EncodeError("ERR wrong  number of Arguments")
		}
		var deleted int
		for i := 1; i < len(args); i++ {
			deleted += s.store.Delete(args[i])
		}
		return protocol.EncodeInteger(deleted)
	case "LPUSH":
		if len(args) != 3 {
			return protocol.EncodeError("ERR wrong number of arguments for LPUSH")
		}
		key := args[1]
		value := args[2]
		newLen := s.store.LPush(key, value)
		return protocol.EncodeInteger(newLen)
	case "LPOP":
		if len(args) != 2 {
			return protocol.EncodeError("ERR wrong number of arguments")
		}
		key := args[1]
		value, ok := s.store.LPop(key)
		if !ok {
			return protocol.EncodeBulkString("")
		}
		return protocol.EncodeBulkString(value)
	default:
		return protocol.EncodeError("ERR unknown command '" + args[0] + "'")
	}
}
