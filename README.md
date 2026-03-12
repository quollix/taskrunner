# Taskrunner

Handling CI-related tasks such as building, testing, and deploying can be automated using bash scripts, but these can quickly get messy. Managing multiple processes, handling cleanup, aborting on failure, adding pretty colored output, and more can lead to complex bash code. That's why Taskrunner was created: a simple, open source library to replace bash scripts with cleaner, more maintainable code written in Go.

Note: API is not stable yet.

### Installation

```bash
go get github.com/ocelot-cloud/taskrunner
```

### Usage Example

```go

var ( 
    backendDir = "../backend"
    frontendDir = "../frontend"
    acceptanceTestsDir = "../acceptance"
)

func TestFrontend() {
    tr := taskrunner.GetTaskRunner()
	
    tr.Log.Info("Testing Integrated Components")
    defer tr.Cleanup() // shuts down the daemon processes at the end
    tr.ExecuteInDir(backendDir, "go build")
    tr.StartDaemon(backendDir, "./backend")
    tr.WaitUntilPortIsReady("8080")

    tr.ExecuteInDir(frontendDir, "npm install")
    tr.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE=TEST")
    tr.WaitForWebPageToBeReady("http://localhost:8081/")
    tr.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE=TEST")
}
```

The idea is to write simple functions like this and build a CLI tool, e.g., by using [cobra](https://github.com/spf13/cobra), to call these functions. The final use of the CLI tool might look like this: 

```bash
go build
./my-task-runner test frontend
```

This approach helps you create a modern and scalable CI infrastructure.

### Contributing

Please read the [Community](https://ocelot-cloud.org/docs/community/) articles for more information on how to contribute to the project and interact with others.

### License

This project is licensed under a permissive open source license, the [0BSD License](https://opensource.org/license/0bsd/). See the [LICENSE](LICENSE) file for details.