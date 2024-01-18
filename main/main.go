package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/KevinZonda/go-jupyter"
	"log"
	"os"
	"runtime"
)

func main() {
	var Version string = "1.0.0"

	// Parse the connection file.
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalln("Need a command line argument specifying the connection file.")
	}

	var connInfo jupyter.ConnectionInfo

	connData, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(connData, &connInfo); err != nil {
		log.Fatal(err)
	}

	// Run the kernel.
	jupyter.RunKernel(connInfo, jupyter.KernelInfo{
		ProtocolVersion:       jupyter.ProtocolVersion,
		Implementation:        "gophernotes",
		ImplementationVersion: Version,
		Banner:                fmt.Sprintf("Go kernel: gophernotes - v%s", Version),
		LanguageInfo: jupyter.KernelLanguageInfo{
			Name:          "go",
			Version:       runtime.Version(),
			FileExtension: ".go",
		},
		HelpLinks: []jupyter.KernelInfoHelpLink{
			{Text: "Go", URL: "https://golang.org/"},
			{Text: "gophernotes", URL: "https://github.com/gopherdata/gophernotes"},
		},
	})
}
