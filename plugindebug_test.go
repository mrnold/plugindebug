package plugindebug

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestScriptCreation(t *testing.T) {
	tmp, _ := os.CreateTemp("/tmp", "plugindebugtest")
	defer os.Remove(tmp.Name())

	testString := "gRPC socket string"
	done, err := createDebugDummyScript(tmp.Name())
	if err != nil {
		t.FailNow()
	}
	fmt.Println(testString)

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("Failed to create test script within three second timeout")
	case <-done:
	}

	line, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal("Failed to read generated test script file")
	}
	if !strings.Contains(string(line), testString) {
		t.Fatalf("Generated test script file has incorrect contents: [%s]", line)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	setVeleroHandshake()
	if _, present := os.LookupEnv("VELERO_PLUGIN"); !present {
		t.Fail()
	}
}
