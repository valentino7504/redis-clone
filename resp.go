package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	t     string
	str   string
	num   int
	bulk  string
	array []Value
}

func (v Value) String() string {
	return fmt.Sprintf(
		"Type: %v, String: %v, Num: %v, Bulk: %v, Array: %+v\n",
		v.t, v.str, v.num, v.bulk, v.array,
	)
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// Reads a line of RESP up until \r\n and parses it.
// Returns the line without the CLRF, no of characters read including CLRF.
// If any errors it returns those instead
func (r *Resp) readLine() ([]byte, int, error) {
	line := make([]byte, 0)
	n := 0
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

// Helper function to read integers - this will be used for size of arrays and
// bulk strings
func (r *Resp) readInteger() (int, int, error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	intVal, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(intVal), n, nil
}

// readArray function. Reads size of the array first then populates the array
// by reading all the bulks in the command
func (r *Resp) readArray() (Value, error) {
	v := Value{t: "array"}
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	v.array = make([]Value, length)
	for i := range v.array {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array[i] = val
	}
	return v, nil
}

// readBulk is a helper to read bulks. Supports readArray when parsing syntax from the CLI.
func (r *Resp) readBulk() (Value, error) {
	v := Value{t: "bulk"}
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	bulk := make([]byte, length)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return v, err
	}
	v.bulk = string(bulk)
	_, _, err = r.readLine()
	if err != nil {
		return v, err
	}
	return v, nil
}

// generic read function that utilises readArray and readBulk
func (r *Resp) Read() (Value, error) {
	// Read the first character of the RESP string to tell the type
	t, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch t {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unrecognised type: %v", string(t))
		return Value{}, nil
	}
}

// marshal Arrays to an array of all the bytes in the command
func (v Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, '\r', '\n')
	for _, val := range v.array {
		bytes = append(bytes, val.Marshal()...)
	}
	return bytes
}

// marshals bulks into arrays of bytes
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshals strings into arrays of bytes
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshals nulls
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// errors will just be stored in v.str
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// generic marshalling function that calls other specific ones
func (v Value) Marshal() []byte {
	switch v.t {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

// writing RESP responses for the client
func (w *Writer) Write(v Value) error {
	bytes := v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
