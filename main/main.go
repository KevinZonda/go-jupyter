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
	jupyter.RunKernel(miniInterpreter{}, connInfo, jupyter.KernelInfo{
		ProtocolVersion:       jupyter.ProtocolVersion,
		Implementation:        "Mini Kernel",
		ImplementationVersion: Version,
		Banner:                fmt.Sprintf("Go kernel: minikernel - v%s", Version),
		LanguageInfo: jupyter.KernelLanguageInfo{
			Name:          "minikernel",
			Version:       runtime.Version(),
			FileExtension: ".mini",
		},
		HelpLinks: []jupyter.KernelInfoHelpLink{
			{Text: "Go", URL: "https://golang.org/"},
			{Text: "gophernotes", URL: "https://github.com/gopherdata/gophernotes"},
		},
	})
}

type miniInterpreter struct{}

func (miniInterpreter) CompleteWords(code string, cursorPos int) (prefix string, completions []string, tail string) {
	return "", nil, ""
}

func (miniInterpreter) Eval(code string) (values []interface{}, err error) {
	return []interface{}{code}, nil
}

func (miniInterpreter) Close() error {
	return nil
}
