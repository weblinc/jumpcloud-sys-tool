package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/weblinc/jumpcloud-sys-tool/jc"
	"io/ioutil"
	"os"
)

const version = "1.0.0"

func main() {

	fmt.Println("jcsystool - JumpCloud System Tool", version)

	var opts struct {
		HttpAction  string `short:"X" long:"action" description:"HTTP method to use e.g. GET/PUT/DELETE" required:"true"`
		JSONContent string `short:"J" long:"json" description:"JSON string to use for PUT actions to system API. Alternatively, use STDIN."`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	httpAction := opts.HttpAction
	jsonContent := opts.JSONContent

	// PUT requests are the only actions that need json to send
	if httpAction == "PUT" {
		if jsonContent == "" {
			jsonContent = readStdin()
		}

		if jsonContent == "" {
			fmt.Println("--json or -J option missing. Data is required to perform a PUT action. Alternatively, you can pipe a JSON string into jcsystool.")
			os.Exit(1)
		}
	} else {
		jsonContent = ""
	}

	// Create the client with config/client key overrides if neccessary
	client, err := jc.NewSystemClient(
		os.Getenv("JC_CONFIG_PATH"),
		os.Getenv("JC_CLIENT_KEY_PATH"),
	)

	if err != nil {
		exitWithError("Error creating client.", err)
	}

	resp, err := client.Do(httpAction, jsonContent)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	body := string(bodyBytes)
	if resp.StatusCode == 200 {
		fmt.Println("Successfully executed", httpAction, "for this system. Response:")
	} else {
		fmt.Println("[", resp.StatusCode, "]", "There was an error executing your request:")
	}

	fmt.Println(body)
}

func exitWithError(desc string, err error) {
	fmt.Println(desc, "\n", err)
	os.Exit(1)
}

// Returns a string from STDIN
func readStdin() string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			exitWithError("Error reading STDIN.", err)
		}
		return string(bytes)
	}

	return ""
}
