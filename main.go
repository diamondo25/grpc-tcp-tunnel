package main

import "os"

func main() {
	if len(os.Args) < 3 {
		println("Usage:")
		println(os.Args[0] + " server server-addr")
		println(os.Args[0] + " client server-addr ip port")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		RunServer()
	case "client":
		RunClient()
	default:
		println(os.Args[1] + "? Thats no option i ever heard of")
		os.Exit(2)
	}
}
