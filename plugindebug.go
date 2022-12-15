package plugindebug

import (
	"fmt"
	"os"
)

func CreateDebugDummyScript(scriptPath string) error {
	fmt.Printf("Saving gRPC socket path to %s", scriptPath)

	if _, handshakePresent := os.LookupEnv("VELERO_PLUGIN"); !handshakePresent {
		if err := os.Setenv("VELERO_PLUGIN", "hello"); err != nil {
			return err
		}
	}

	os.Args[0] = scriptPath

	return nil
}
