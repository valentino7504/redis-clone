package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error encountered: ", err)
		return
	}
	fmt.Println("Listening on port 6379")
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("Error closing the connection:", err)
		}
	}()
	for {
		resp := NewResp(conn)
		val, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		if val.t != "array" {
			fmt.Println("Invalid request, expected array of bulks")
			continue
		}
		if len(val.array) == 0 {
			fmt.Println("Error, expected array of length >= 1")
			continue
		}

		// convert command to uppercase - redis is not case-sensitive
		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		writer := NewWriter(conn)
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command - valid commands include SET, GET and PING")
			_ = writer.Write(Value{t: "string", str: ""})
			continue
		}

		// write the response
		response := handler(args)
		err = writer.Write(response)
		if err != nil {
			fmt.Println("Error writing output:", err)
		}
	}
}
