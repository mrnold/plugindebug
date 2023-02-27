package plugindebug

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

// Save the script path so it can be cleaned up afterward.
// Justification for global is that EnableDebug should only be getting called
// once for the runtime of a particular plugin.
var generatedScript string

// EnableDebug enables local debugging of a Velero plugin, assuming
// Velero itself is also running locally in a separate debug session.
// Accepts a path to a plugin directory, which must also be provided
// to that other running Velero server instance with the "--plugin-dir"
// option, and accepts a name for the current plugin. The name does not
// have to match the actual plugin name.
//
// Usage: in the plugin source, import "github.com/mrnold/plugindebug".
//
//	Then call plugindebug.EnableDebug before calling the Serve
//	function of your registered plugin.
//
// This will set the VELERO_PLUGIN environment variable to "hello" to
// allow the plugin to start up, then create an executable script in the
// specified plugin directory. The script will contain the gRPC socket
// information printed by the plugin code, so that the main Velero process
// connects to that same socket for every plugin operation. This allows
// you to run the plugin under a debugger, setting breakpoints and
// inspecting values and whatever else. This is very much a hack and
// should not be expected to work long-term or for every conceivable use case.
func EnableDebug(pluginDir, pluginName string) error {
	if err := setVeleroHandshake(); err != nil {
		return err
	}
	plugin := path.Join(pluginDir, pluginName)
	if _, err := createDebugDummyScript(plugin); err != nil {
		return err
	}
	return nil
}

// Cleanup removes the generated script file.
// Call this after the plugin's Serve is done.
func Cleanup() {
	if generatedScript == "" {
		return
	}
	_, err := os.Stat(generatedScript)
	if os.IsNotExist(err) {
		return
	}
	os.Remove(generatedScript)
}

func setVeleroHandshake() error {
	if _, handshakePresent := os.LookupEnv("VELERO_PLUGIN"); !handshakePresent {
		if err := os.Setenv("VELERO_PLUGIN", "hello"); err != nil {
			return err
		}
	}
	return nil
}

func createDebugDummyScript(scriptPath string) (chan struct{}, error) {
	savedPluginBinaryPath := os.Args[0]
	generatedScript = scriptPath
	os.Args[0] = scriptPath
	done := make(chan struct{}, 1)

	readEnd, writeEnd, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	savedStdout := os.Stdout
	os.Stdout = writeEnd

	go func() {
		scanner := bufio.NewScanner(readEnd)
		for scanner.Scan() {
			line := scanner.Text()
			script := fmt.Sprintf("#!/usr/bin/env bash\necho '%s'", line)
			os.WriteFile(scriptPath, []byte(script), 0755)

			os.Stdout = savedStdout
			os.Args[0] = savedPluginBinaryPath
			fmt.Printf("Wrote gRPC socket path '%s' to dummy plugin script '%s'\n", line, scriptPath)
			done <- struct{}{}
			break
		}
	}()

	return done, nil
}
