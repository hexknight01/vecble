package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/cockroachdb/pebble"
)

const (
	redisOK     = "+OK\r\n"
	redisNil    = "$-1\r\n"
	redisPrefix = "*"
)

var (
	db   *pebble.DB
	lock sync.RWMutex
)

func main() {
	var err error
	db, err = pebble.Open("pebble_data", &pebble.Options{})
	if err != nil {
		log.Fatalf("Failed to open Pebble DB: %v", err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("Redis-compatible server running on :6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		log.Printf("Client disconnected: %s", conn.RemoteAddr().String())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	// conn.Write([]byte("+Hello! Welcome to Pebble-Redis.\r\n"))

	for {
		// Read the first line to determine the command type
		line, err := reader.ReadString('\n')
		// line, err := reader.ReadString('\n')
		// if err != nil {
		// 	log.Printf("Connection error: %v", err)
		// 	return
		// }
		log.Printf("Received: %s", line)

		// Check if the command follows RESP format
		// if !strings.HasPrefix(firstLine, redisPrefix) {
		// log.Println("Invalid command format")
		// conn.Write([]byte("-ERR invalid command\r\n"))
		// continue
		// }

		// Parse the full RESP command
		command, args := parseRESP(line)
		log.Printf("Parsed Command: %s, Args: %v", command, args)

		// Execute and return response
		response := handleCommand(command, args)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
	}
}

func parseRESP(line string) (string, []string) {
	line = strings.TrimSpace(line) // Trim \r\n
	println(line)
	// Check if it's an array-type command (starts with "*")
	if !strings.HasPrefix(line, "*") {
		log.Printf("Invalid RESP format: %q", line)
		return "", nil
	}

	// Read the command name (e.g., "$4\r\nPING\r\n")
	// cmdLength, _ := reader.ReadString('\n') // Read "$4\r\n"
	cmd, _ := reader.ReadString('\n') // Read "PING\r\n"

	cmd = strings.TrimSpace(cmd)
	cmd = strings.ToUpper(cmd) // Convert to uppercase

	// Read additional arguments if available
	var args []string
	for {
		argLength, err := reader.ReadString('\n')
		if err != nil || !strings.HasPrefix(argLength, "$") {
			break
		}
		arg, _ := reader.ReadString('\n')
		args = append(args, strings.TrimSpace(arg))
	}

	return cmd, args
}

// Handles Redis commands
func handleCommand(cmd string, args []string) string {
	cmd = strings.ToLower(cmd) // Ensure case-insensitive matching
	log.Printf("Executing command: %s, Args: %v", cmd, args)
	switch cmd {
	case "ping":
		log.Print("Executing command")
		return "+OK\r\n"
	// case "SET":
	// 	if len(args) != 2 {
	// 		return "-ERR wrong number of arguments for 'set' command\r\n"
	// 	}
	// 	lock.Lock()
	// 	db.Set([]byte(args[0]), []byte(args[1]), nil)
	// 	lock.Unlock()
	// 	return redisOK
	// case "GET":
	// 	if len(args) != 1 {
	// 		return "-ERR wrong number of arguments for 'get' command\r\n"
	// 	}
	// 	lock.RLock()
	// 	value, closer, err := db.Get([]byte(args[0]))
	// 	lock.RUnlock()
	// 	if err != nil {
	// 		return redisNil
	// 	}
	// 	defer closer.Close()
	// 	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	default:
		return "-ERR unknown command\r\n"
	}
}
