package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"golang.design/x/clipboard"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	link := os.Args[1]
	org, project := extractData(link)
	orgDir := filepath.Join(homedir, "work/github", org)
	if err := os.MkdirAll(orgDir, 0o755); err != nil {
		panic(err)
	}

	if err := exe(fmt.Sprintf("gh repo clone %s/%s", org, project), orgDir); err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("cd %s", path.Join(orgDir, project))

	fmt.Printf("\nCopied \n%s\nto clipboard", color.HiGreenString(msg))
	_ = clipboard.Write(clipboard.FmtText, []byte(msg))
}

func extractData(link string) (string, string) {
	parsedUrl, err := url.ParseRequestURI(link)
	if err != nil {
		panic(err)
	}
	org, projectRaw, found := strings.Cut(strings.TrimPrefix(parsedUrl.Path, "/"), "/")
	if !found {
		panic(fmt.Errorf("couldnt cut org and project from %q, full: %#v", parsedUrl.Path, parsedUrl))
	}

	if !strings.Contains(projectRaw, "/") {
		return org, projectRaw
	}

	project, _, found := strings.Cut(projectRaw, "/")
	if !found {
		panic(fmt.Errorf("couldnt cut project from %q", projectRaw))
	}

	return org, project
}

func exe(command, dir string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
