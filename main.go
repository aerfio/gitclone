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
	orgDir := filepath.Join(homedir, "work/github.com", org)
	if err := os.MkdirAll(orgDir, 0o755); err != nil {
		panic(err)
	}

	if err := clone(org, project, orgDir); err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("cd %s", path.Join(orgDir, project))

	fmt.Printf("\nCopied \n%s\nto clipboard", color.HiGreenString(msg))
	_ = clipboard.Write(clipboard.FmtText, []byte(msg))
}

func extractData(link string) (string, string) {
	if strings.HasPrefix(link, "git@") && strings.HasSuffix(link, ".git") {
		return handleSSHLink(link)
	}

	parsedUrl, err := url.ParseRequestURI(strings.TrimSuffix(link, ".git"))
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

func handleSSHLink(link string) (string, string) {
	orgRepo := strings.Split(
		strings.TrimSuffix(
			strings.TrimPrefix(link, "git@github.com:"),
			".git"),
		"/")
	if len(orgRepo) != 2 {
		panic(fmt.Errorf("%+v should have 2 elements", orgRepo))
	}
	return orgRepo[0], orgRepo[1]
}

func clone(org, repo string, dir string) error {
	cmd := exec.Command("gh", "repo", "clone", org+"/"+repo)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
