package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ProjectConfig struct {
	name         string
	template     string
	shouldInstallDeps bool
	packageManager string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	projectConfig := getProjectConfig(reader)
	switch projectConfig.template {
	case "react-ts":
		createReactTsProject(projectConfig)
	case "vite-react":
		createViteReactProject(projectConfig)
	case "nextjs-router":
		createNextJsProject(projectConfig)
	}
	createReactTsProject(projectConfig)
}

func getProjectConfig(reader *bufio.Reader) ProjectConfig {
	config := ProjectConfig{}
	config.name = prompt(reader, "📘 Enter project name: ")
	fmt.Println("\n📋 Select template:")
	fmt.Println("=======================================")
	fmt.Println("\n1️⃣ \tMy React-ts with Rsbuild")
	fmt.Println("\n2️⃣ \tVite + React")
	fmt.Println("\n3️⃣ \tNext.js app router")

	templateChoice := prompt(reader, "\n💻 Choose your template (1-3): ")
	switch templateChoice {
	case "1":
		config.template = "react-ts"
	case "2":
		config.template = "vite-react"
	case "3":
		config.template = "nextjs-router"
	default:
		config.template = "react-ts"
		fmt.Println("\nUsing default template: react-ts")
	}
	installDepsCheck := prompt(reader, "\n📦 Do you want to install dependencies (y/n): ")
	if strings.ToLower(installDepsCheck) == "y" || installDepsCheck == "" {
		config.shouldInstallDeps = true
		fmt.Println("\n📦 Choose your package manager: ")
		fmt.Println("=======================================")
		fmt.Println("\n1️⃣ \tnpm")
		fmt.Println("\n2️⃣ \tyarn")
		fmt.Println("\n3️⃣ \tpnpm")
		packageManagerChoice := prompt(reader, "\n🤔 Choose between 1-3: ")
		switch packageManagerChoice {
		case "1":
			config.packageManager = "npm"
		case "2":
			config.packageManager = "yarn"
		case "3":
			config.packageManager = "pnpm"
		default:
			config.packageManager = "npm"
			fmt.Println("Using default package manager: npm")
		}
	} else {
		config.shouldInstallDeps = false
	}
	fmt.Println("\n=======================================")
	fmt.Println("\n✅ Project Confirmation:")
	fmt.Printf("\nProject name: %s\n", config.name)
	fmt.Printf("\nTemplate: %s\n", config.template)
	fmt.Printf("\nInstall dependencies: %t\n", config.shouldInstallDeps)
	if config.shouldInstallDeps {
		fmt.Printf("\nPackage manager: %s\n", config.packageManager)
	}
	confirm := prompt(reader, "\nIs this correct? (Y/n): ")
	if strings.ToLower(confirm) == "n" {
		fmt.Println("\n\nExiting... 👋")
		os.Exit(0)
		} else {
		fmt.Println("\n✨Creating project...")
		return config
	}
	return config
}

func createReactTsProject(config ProjectConfig) error{
	if err := runCommand("git", "clone", "https://github.com/aryan1306/react-ts.git", config.name); err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}
	if err := os.Chdir(config.name); err != nil {
		return fmt.Errorf("error changing directory: %v", err)
	}
	if err := os.RemoveAll(".git"); err != nil {
		return fmt.Errorf("error removing .git directory: %v", err)
	}
	if err := runCommand("git", "init"); err != nil {
		return fmt.Errorf("error initializing new git repository: %v", err)
	}
	if config.shouldInstallDeps{
		if err:= runCommand(config.packageManager, "install"); err != nil {
			return fmt.Errorf("error installing dependencies: %v", err)
		}
	}
	return nil
}

func createViteReactProject(config ProjectConfig) error{
	if err:= runCommand(config.packageManager, "create", "vite@latest", config.name, "--template", "react-ts"); err != nil {
		return fmt.Errorf("error creating vite project: %v", err)
	}
	if config.shouldInstallDeps {
		if err:= os.Chdir(config.name); err != nil {
			return fmt.Errorf("error changing directory: %v", err)
		}
		if err:= runCommand(config.packageManager, "install"); err != nil {
			return fmt.Errorf("error installing dependencies: %v", err)
		}
	}
	return nil
}

func createNextJsProject(config ProjectConfig) error{
	var pkgManager string
	switch config.packageManager {
	case "npm":
		pkgManager = "--use-npm"
	case "yarn":
		pkgManager = "--use-yarn"
	case "pnpm":
		pkgManager = "--use-pnpm"
	default:
		pkgManager = "--use-npm"
	}
	if err := runCommand("npx", "create-next-app@latest", config.name, "--ts", "--tailwind", "--eslint", "--app", "--src-dir", "--skip-install", "--turbopack no", `--import-alias "@/*"`, pkgManager); err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}
	if config.shouldInstallDeps{
		if err:= runCommand(config.packageManager, "install"); err != nil {
			return fmt.Errorf("error installing dependencies: %v", err)
		}
	}
	return nil
}

func prompt(reader *bufio.Reader, question string) string {
	fmt.Print(question)
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v", err)
		os.Exit(1)
	}
	return strings.TrimSpace(answer)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}