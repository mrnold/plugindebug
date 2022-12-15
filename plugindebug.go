package plugindebug

import (
	"fmt"
)

func CreateDebugDummyScript(scriptPath string) error {
	fmt.Printf("Saving gRPC socket path to %s", scriptPath)
	return nil
}
