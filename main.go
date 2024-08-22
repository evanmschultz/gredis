package main

import (
	"fmt"
	"net"
	"strings"
)

// main is the entry point for the Redis-compatible server. It listens on port :6379 for incoming connections,
// reads commands from the connection, and executes the appropriate handler for the command. It also reads
// commands from an append-only file (AOF) and replays them on startup.
func main() {
	fmt.Println("Listening on port :6379")

	// Listen listens on the default Redis port (:6379) for incoming TCP connections.
 	// If an error occurs while listening, it is printed to the console and the program exits.
 l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// NewAof creates a new append-only file (AOF) at the specified path. If the file does not exist, it is created.
	// If an error occurs while opening or creating the file, it is returned.
	// The AOF is used to store and replay commands executed by the Redis-compatible server.
 aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// aof.Read reads commands from the append-only file (AOF) and executes them. For each command read from the AOF:
	// - The command name is extracted from the first element of the command array.
	// - The command arguments are extracted from the remaining elements of the command array.
	// - The appropriate command handler is looked up in the Handlers map.
	// - If the command handler is found, it is called with the extracted arguments.
	// - If the command handler is not found, an error message is printed.
 aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	// Accept accepts an incoming TCP connection on the listener l. If an error occurs while accepting the connection,
	// it is printed to the console and the function returns.
	conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

	// Close the connection when the function returns.
	defer conn.Close()

	// The main loop of the Redis-compatible server. It reads requests from the client connection,
	// processes the commands, and writes the responses back to the client.
	// For each request:
	// - The request is read from the connection using NewResp().
	// - The command name and arguments are extracted from the request.
	// - The appropriate command handler is looked up in the Handlers map.
	// - If the command handler is found, it is called with the extracted arguments, and the result is written back to the client using NewWriter().
	// - If the command handler is not found, an error message is written back to the client.
	// - If the command is "SET" or "HSET", the request is also written to the append-only file (AOF) using aof.Write().
 	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}