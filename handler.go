package main

import (
	"sync"
)

// Handlers is a map of command names to their corresponding handler functions.
// The handlers are used to process different types of commands that can be
// executed by the application.
var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

// ping is a command handler that responds with "PONG" if no arguments are provided,
// or echoes the first argument back as a string.
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

// SETs is a map that stores key-value pairs for the "SET" command.
var SETs = map[string]string{}

// SETsMu is a read-write mutex that protects access to the SETs map.
var SETsMu = sync.RWMutex{}

// set is a command handler that sets a key-value pair in the SETs map.
// It takes two arguments: the key and the value to be set.
// If the number of arguments is not exactly 2, it returns an error.
// The function acquires a write lock on the SETsMu mutex before modifying the SETs map,
// and releases the lock after the operation is complete.
// It returns a Value with a "string" type and the value "OK" upon successful completion.
func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

// get is a command handler that retrieves the value associated with a given key
// from the SETs map. It takes one argument: the key to retrieve.
// If the number of arguments is not exactly 1, it returns an error.
// The function acquires a read lock on the SETsMu mutex before accessing the SETs map,
// and releases the lock after the operation is complete.
// If the key is not found in the SETs map, it returns a Value with a "null" type.
// Otherwise, it returns a Value with a "bulk" type containing the value associated with the key.
func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// HSETs is a map that stores hash sets. The outer map maps hash names to inner maps,
// and the inner maps map keys to values within each hash set.
var HSETs = map[string]map[string]string{}

// HSETsMu is a read-write mutex that protects access to the HSETs map.
var HSETsMu = sync.RWMutex{}

// hset is a command handler that adds or updates a key-value pair in a hash set.
// It takes three arguments: the name of the hash set, the key, and the value.
// If the number of arguments is not exactly 3, it returns an error.
// The function acquires a write lock on the HSETsMu mutex before modifying the HSETs map,
// and releases the lock after the operation is complete.
// If the hash set does not exist, it creates a new one before adding the key-value pair.
// It returns a Value with a "string" type and the value "OK" upon successful completion.
func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

// hget is a command handler that retrieves the value associated with a key in a hash set.
// It takes two arguments: the name of the hash set and the key.
// If the number of arguments is not exactly 2, it returns an error.
// The function acquires a read lock on the HSETsMu mutex before accessing the HSETs map,
// and releases the lock after the operation is complete.
// If the key does not exist in the hash set, it returns a null value.
// Otherwise, it returns the value associated with the key as a bulk string.
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// hgetall is a command handler that retrieves all key-value pairs in a hash set.
// It takes one argument: the name of the hash set.
// If the number of arguments is not exactly 1, it returns an error.
// The function acquires a read lock on the HSETsMu mutex before accessing the HSETs map,
// and releases the lock after the operation is complete.
// If the hash set does not exist, it returns a null value.
// Otherwise, it returns an array of all the key-value pairs in the hash set.
func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}