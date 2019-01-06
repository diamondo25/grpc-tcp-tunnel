package main

import "os"

func main() {
	if len(os.Args) < 2 {
		println(os.Args[0]+" server")
		println(os.Args[0]+" client")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		RunServer()
	case "client":
		RunClient()
	default:
		println(os.Args[1]+"? Thats no option i ever heard of")
		os.Exit(2)
	}
}
