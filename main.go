package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

const (

	// Version defines the gophernotes version.
	Version string = "1.0.0"

	// ProtocolVersion defines the Jupyter protocol version.
	ProtocolVersion string = "5.0"
)

func main() {

	// Parse the connection file.
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalln("Need a command line argument specifying the connection file.")
	}

	var connInfo ConnectionInfo

	connData, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(connData, &connInfo); err != nil {
		log.Fatal(err)
	}

	// Run the kernel.
	RunKernel(connInfo)
}
