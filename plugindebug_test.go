package plugindebug

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestScriptCreation(t *testing.T) {
	tempDir := t.TempDir()
	scriptName := "cleanuptest"
	pluginPath := path.Join(tempDir, scriptName)

	testString := "gRPC socket string"
	done, err := createDebugDummyScript(pluginPath)
	if err != nil {
		t.Fatalf("Failed to create test script: %s", err)
	}

	// Write test contents to stdout to be picked up by the goroutine
	// waiting inside createDebugDummyScript.
	fmt.Println(testString)

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("Failed to create test script within three second timeout")
	case <-done:
	}

	line, err := os.ReadFile(pluginPath)
	if err != nil {
		t.Fatal("Failed to read generated test script file")
	}
	if !strings.Contains(string(line), testString) {
		t.Fatalf("Generated test script file has incorrect contents: [%s]", line)
	}

	Cleanup()
	_, err = os.Stat(pluginPath)
	if !os.IsNotExist(err) {
		t.Fatalf("Failed to clean up temporary file %s", pluginPath)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	setVeleroHandshake()
	if _, present := os.LookupEnv("VELERO_PLUGIN"); !present {
		t.Fail()
	}
}
