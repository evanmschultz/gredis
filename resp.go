package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// The constants STRING, ERROR, INTEGER, BULK, and ARRAY represent the different types of values that can be returned in a RESP (Redis Serialization Protocol) response.
// STRING represents a string value, ERROR represents an error value, INTEGER represents an integer value, BULK represents a bulk string value, and ARRAY represents an array of values.
// These constants are used throughout the package to handle and parse RESP responses.
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value represents a value in the RESP (Redis Serialization Protocol) format. It can be one of several types:
// - STRING: a string value
// - ERROR: an error value
// - INTEGER: an integer value
// - BULK: a bulk string value
// - ARRAY: an array of values
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

// Resp is a struct that holds a bufio.Reader for reading RESP (Redis Serialization Protocol) responses.
type Resp struct {
	reader *bufio.Reader
}

// NewResp creates a new Resp instance that reads from the provided io.Reader.
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads a line of text from the Resp's reader, excluding the trailing newline characters.
// It returns the line as a byte slice, the number of bytes read, and any error that occurred during the read.
// The function reads bytes from the reader until it encounters a newline character, and returns the line
// excluding the trailing newline characters.
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger reads an integer value from the Resp's reader.
// It reads a line of text from the reader, converts it to an integer,
// and returns the integer value, the number of bytes read, and any error that occurred.
// The function assumes the line contains a valid integer representation.
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Read reads a RESP value from the Resp's reader. It determines the type of the value
// based on the first byte read, and then calls the appropriate parsing function to
// read the value. If the type is unknown, it prints a message and returns an empty
// Value and a nil error.
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// readArray reads an array value from the Resp's reader. It reads the length of the
// array, then reads each element of the array and appends it to the array field of
// the returned Value. If any errors occur during reading, the function returns the
// error.
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.array = append(v.array, val)
	}

	return v, nil
}

// readBulk reads a bulk value from the Resp's reader. It reads the length of the
// bulk string, then reads the bytes of the string and stores them in the bulk
// field of the returned Value. If any errors occur during reading, the function
// returns the error.
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}


// Marshal returns the RESP representation of the Value. The representation
// depends on the type of the Value, which is stored in the typ field.
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

// marshalString returns the RESP representation of a string value. It prepends
// the string type identifier, appends the string value, and adds the trailing
// CRLF.
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalBulk returns the RESP representation of a bulk string value. It prepends
// the bulk string type identifier, appends the length of the string value, adds
// the trailing CRLF, and then appends the string value followed by another CRLF.
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalArray returns the RESP representation of an array value. It prepends
// the array type identifier, appends the length of the array, adds the trailing
// CRLF, and then appends the RESP representation of each element in the array.
func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

// marshallError returns the RESP representation of an error value. It prepends
// the error type identifier, appends the error message, and adds the trailing
// CRLF.
func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshallNull returns the RESP representation of a null value. It prepends the
// null type identifier and adds the trailing CRLF.
func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}


// Writer is a struct that wraps an io.Writer and provides a Write method to write RESP-encoded values.
type Writer struct {
	writer io.Writer
}

// NewWriter creates a new Writer that writes RESP-encoded values to the provided io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write writes the RESP-encoded representation of the provided Value to the
// underlying io.Writer. It returns an error if the write operation fails.
func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}