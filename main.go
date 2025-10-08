package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RootDir       string `yaml:"root_dir"` 
	DefaultType   string `yaml:"default_type"` 
	DefaultDomain string `yaml:"default_domain"` 
	GitInit       bool   `yaml:"git_init"` 
	GitHubToken   string `yaml:"github_token"` 
}

var prefixes = []string{"ARES", "ARGUS", "HYDRA", "PHOENIX", "RAVEN", "NEXUS", "CRONUS", "VORTEX", "SIGMA", "ECHO"}
var types = []string{"LAB", "PROTO", "OPS", "CORE", "SYS", "UI", "AI", "DEMO"}
var domains = []string{"AUTH", "API", "DASH", "BILLING", "ML", "STORAGE", "STREAM", "CRON", "MARKET", "METRICS"}

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("Использование: dnp create [type] [domain] [--prefix=NAME] [--desc='описание'] [--dir=/путь]")
		return
	}

	switch os.Args[1] {
	case "create":
		createProject()
	case "list":
		listProjects()
	default:
		fmt.Println("Команды: create | list")
	}
}

func loadConfig() Config {
	path := filepath.Join(os.Getenv("HOME"), ".dnp", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{
			RootDir:       filepath.Join(os.Getenv("HOME"), "Projects", "D"),
			DefaultType:   "LAB",
			DefaultDomain: "CORE",
			GitInit:       true,
			GitHubToken:   "",
		}
	}
	var cfg Config
	yaml.Unmarshal(data, &cfg)
	return cfg
}

func createProject() {
	cfg := loadConfig()
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	prefixFlag := createCmd.String("prefix", "", "Имя проекта (опционально)")
	descFlag := createCmd.String("desc", "", "Описание проекта")
	dirFlag := createCmd.String("dir", cfg.RootDir, "Путь, где создавать проект")
	createCmd.Parse(os.Args[2:])

	args := createCmd.Args()
	projectType := getArg(args, 0, cfg.DefaultType)
	domain := getArg(args, 1, cfg.DefaultDomain)

	prefix := *prefixFlag
	if prefix == "" {
		prefix = prefixes[rand.Intn(len(prefixes))]
	}

	name := fmt.Sprintf("%s-%s-%s", strings.ToUpper(prefix), strings.ToUpper(projectType), strings.ToUpper(domain))
	projectPath := filepath.Join(*dirFlag, strings.ToLower(name))

	os.MkdirAll(filepath.Join(projectPath, "cmd"), 0755)
	os.MkdirAll(filepath.Join(projectPath, "internal", "core"), 0755)

	// main.go
	mainCode := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Project %s initialized.")
}
`, name)
	os.WriteFile(filepath.Join(projectPath, "cmd", "main.go"), []byte(mainCode), 0644)

	// go.mod
	moduleName := fmt.Sprintf("github.com/%s/%s", os.Getenv("USER"), strings.ToLower(name))
	goMod := fmt.Sprintf("module %s\n\ngo 1.23\n", moduleName)
	os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)

	// Makefile
	makefile := `run:
	go run ./cmd/main.go

build:
	go build -o bin/app ./cmd/main.go

test:
	go test ./...
`
	os.WriteFile(filepath.Join(projectPath, "Makefile"), []byte(makefile), 0644)

	// README
	readme := fmt.Sprintf(`# %s

%s

**Type:** %s  
**Domain:** %s  
**Created:** %s  
`, name, *descFlag, projectType, domain, time.Now().Format("2006-01-02"))
	os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(readme), 0644)

	// .gitignore
	gitignore := "bin/\n*.log\n*.tmp\n.env\n"
	os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte(gitignore), 0644)

	// Git init
	if cfg.GitInit {
		cmd := exec.Command("git", "init")
		cmd.Dir = projectPath
		cmd.Run()
		exec.Command("git", "-C", projectPath, "add", ".").Run()
		exec.Command("git", "-C", projectPath, "commit", "-m", "Initial commit").Run()
		
		// Create GitHub repo and push
		if cfg.GitHubToken != "" {
			createGitHubRepoAndPush(cfg.GitHubToken, projectPath, strings.ToLower(name), *descFlag)
		}
	}

	// Log registry
	logProject(name, projectPath)

	fmt.Printf("Создан проект: %s\n", name)
	fmt.Printf("Расположение: %s\n", projectPath)
	if cfg.GitInit {
		fmt.Println("Git инициализирован: ✅")
	}
}

func listProjects() {
	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".dnp", "projects.log"))
	if err != nil {
		fmt.Println("Проекты пока не создавались.")
		return
	}
	fmt.Println("== D Project Registry ==")
	fmt.Println(string(data))
}

func logProject(name, path string) {
	regDir := filepath.Join(os.Getenv("HOME"), ".dnp")
	os.MkdirAll(regDir, 0755)
	f, _ := os.OpenFile(filepath.Join(regDir, "projects.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s — %s (%s)\n", time.Now().Format("2006-01-02 15:04"), name, path))
}

func getArg(args []string, index int, fallback string) string {
	if len(args) > index {
		return strings.ToUpper(args[index])
	}
	return fallback
}

func createGitHubRepoAndPush(token, projectPath, repoName, description string) {
	ctx := context.Background()
	client := github.NewTokenClient(ctx, token)
	
	// Get current user
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("Ошибка получения пользователя GitHub: %v\n", err)
		return
	}
	
	// Create repository
	repo := &github.Repository{
		Name:        &repoName,
		Description: &description,
		Private:     github.Bool(false), // Set to true if you want private repos
	}
	createdRepo, _, err := client.Repositories.Create(ctx, "", repo)
	if err != nil {
		fmt.Printf("Ошибка создания репозитория GitHub: %v\n", err)
		return
	}
	
	// Set remote origin
	originURL := fmt.Sprintf("https://github.com/%s/%s.git", *user.Login, repoName)
	exec.Command("git", "-C", projectPath, "remote", "add", "origin", originURL).Run()
	
	// Push to GitHub
	exec.Command("git", "-C", projectPath, "push", "-u", "origin", "main").Run()
	
	fmt.Printf("GitHub репозиторий создан: %s\n", *createdRepo.HTMLURL)
	fmt.Println("Первый коммит запушен на GitHub: ✅")
}
