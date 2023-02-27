# Velero plugin debug
This package is a hack to enable local debugging of Velero plugins.

## Quick start
Get a Velero server up and running, with the "--plugin-dir" option set somewhere on the local file system. Then, call plugindebug.EnableDebug inside the code of the plugin you want to debug, using that same plugin directory:
```
import "github.com/mrnold/plugindebug"

func main() {
    plugindebug.EnableDebug("/home/marnold/Work/velero/plugin-directory", "velero-plugin-test")
    defer plugindebug.Cleanup()
    ...
    framework.NewServer()....Serve()
}
```
You should now be able to run the plugin under a debugger and trigger breakpoints when the main velero process asks the plugin to do something.

## What this does
This package attempts to trick a Velero plugin into running as a standalone local program. The first thing it does is to set the VELERO_PLUGIN handshake to "hello", to avoid "This binary is a plugin" messages. Then, it temporarily redirects STDOUT so that it can copy the gRPC socket information printed by the go-plugin Serve() call. It copies the socket path into a small shell script in the Velero plugin directory, so that the main velero process will see the executable and attempt to run it as a plugin. Normally, running the plugin itself would change the gRPC socket on every run, but this way forces velero to connect to the long-running plugin repeatedly. Finally, it temporarily sets os.Args[0] to point to the script path, so that Serve() responds with the path to the script file when the main Velero process asks.

This sequence of hacks just about gets Velero plugins to work with a debugger, although there is currently plenty of logging from Velero indicating that it is not happy about being run this way. This technique is sort of useful for exploration or small fixes, but probably should not be relied upon for serious development tasks.

## Problems
* In VSCode, the "Stop Debugging" button ends up sending a SIGKILL to the plugin process, so the plugin never gets a chance to run the clean up code. Debugging multitple plugins this way will leave the plugin directory full of scripts pointing to closed gRPC sockets, and Velero will not be able to start up all the way. To get around this, plugindebug.EnableDebug creates a cleanup script that can be used as a "Post Debug Task", like this:

    ```
    launch.json
    {
        ...
        "configurations": [
            {
                ...
                "postDebugTask": "Clean Up"
            }
        ]
    }
    ```
    ```
    tasks.json
    {
        ...
        "tasks": [
            {
                "label": "Clean Up",
                "type": "shell",
                "command": "./.cleanup.sh && rm .cleanup.sh"
            }
        ]
    }
    ```
    This workaround should remove the script after debugging stops.
