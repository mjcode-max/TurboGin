package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ProjectName    = "TurboGin"
	BinaryName     = "turbogin"
	Version        = "0.0.1"
	DockerImage    = "turbogin-app"
	DockerTag      = "latest"
	WireCmd        = "wire"
	WireGenPath    = "./internal/wire"
	MainModulePath = "./cmd/server/main.go"
	BinDir         = "bin"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	switch command {
	case "init":
		initProject()
	case "deps":
		downloadDependencies()
	case "install-tools":
		installTools()
	case "generate":
		generateWire()
	case "build":
		build(false)
	case "build-linux":
		build(true)
	case "run":
		run()
	case "clean":
		clean()
	case "docker-build":
		dockerBuild()
	case "docker-run":
		dockerRun()
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
	}
}

func initProject() {
	fmt.Println("Initializing project...")
	runCommand("go", "mod", "tidy")
	installTools()
}

func downloadDependencies() {
	fmt.Println("Downloading dependencies...")
	runCommand("go", "mod", "download")
}

func installTools() {
	fmt.Println("Installing required tools...")
	runCommand("go", "install", "github.com/google/wire/cmd/wire@latest")
}

func generateWire() {
	fmt.Println("Generating Wire dependencies...")
	runCommand(WireCmd, "gen", WireGenPath)
}

func build(linux bool) {
	generateWire()
	fmt.Println("Building application with optimization...")

	ldflags := "-s -w"
	output := filepath.Join(BinDir, BinaryName)

	if linux {
		fmt.Println("Target: Linux (amd64)")
		os.Setenv("GOOS", "linux")
		os.Setenv("GOARCH", "amd64")
		output = filepath.Join(BinDir, BinaryName+"-linux")
	}

	args := []string{
		"build",
		"-ldflags=" + ldflags,
		"-trimpath",
		"-o", output,
		MainModulePath,
	}

	runCommand("go", args...)
}

func run() {
	generateWire()
	fmt.Println("Starting application...")
	runCommand("go", "run", MainModulePath)
}

func clean() {
	fmt.Println("Cleaning up...")
	os.RemoveAll(BinDir)
	os.RemoveAll("docs")
}

func dockerBuild() {
	fmt.Println("Building Docker image...")
	runCommand("docker", "build", "-t", fmt.Sprintf("%s:%s", DockerImage, DockerTag), ".")
}

func dockerRun() {
	fmt.Println("Running Docker container...")
	runCommand("docker", "run", "-p", "8080:8080", "--name", BinaryName, "--rm", fmt.Sprintf("%s:%s", DockerImage, DockerTag))
}

func runCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}

func printHelp() {
	fmt.Printf("%s Scaffolding - Turbo CLI Help\n\n", ProjectName)
	fmt.Println("Commands:")
	fmt.Println("  init           - Initialize the project (go mod init)")
	fmt.Println("  deps           - Download all dependencies")
	fmt.Println("  install-tools  - Install required tools (wire)")
	fmt.Println("  generate       - Generate Wire dependencies")
	fmt.Println("  build          - Build the application")
	fmt.Println("  build-linux    - Build the application for linux")
	fmt.Println("  run            - Run the application")
	fmt.Println("  clean          - Clean generated files")
	fmt.Println("  docker-build   - Build Docker image")
	fmt.Println("  docker-run     - Run Docker container")
	fmt.Println("  help           - Show this help message")
}
