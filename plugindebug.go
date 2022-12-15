package plugindebug

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func CreateDebugDummyScript(scriptPath string) error {
	if _, handshakePresent := os.LookupEnv("VELERO_PLUGIN"); !handshakePresent {
		if err := os.Setenv("VELERO_PLUGIN", "hello"); err != nil {
			return err
		}
	}

	savedPluginBinaryPath := os.Args[0]
	os.Args[0] = scriptPath

	readEnd, writeEnd, err := os.Pipe()
	if err != nil {
		return err
	}

	savedStdout := os.Stdout
	os.Stdout = writeEnd

	go func() {
		scanner := bufio.NewScanner(readEnd)
		for scanner.Scan() {
			line := scanner.Text()
			script := fmt.Sprintf("#!/usr/bin/env bash\necho '%s'", line)
			ioutil.WriteFile(scriptPath, []byte(script), 0755)

			os.Stdout = savedStdout
			os.Args[0] = savedPluginBinaryPath
			fmt.Printf("Wrote gRPC socket path '%s' to dummy plugin script '%s'\n", line, scriptPath)
			break
		}
	}()

	return nil
}
