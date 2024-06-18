package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"golang.design/x/clipboard"
)

func main() {
	root := "work"
	if envRoot := os.Getenv("GITCLONE_ROOT_DIR"); envRoot != "" {
		root = envRoot
	}

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
	orgDir := filepath.Join(homedir, root, "github.com", org)

	exists, err := checkIfExists(filepath.Join(orgDir, project))
	if err != nil {
		panic(err)
	}

	changeDirMsg := fmt.Sprintf("cd %s", filepath.Join(orgDir, project))
	fmt.Printf("Cloning %s into %s\n", color.GreenString(fmt.Sprintf("%s/%s", org, project)), filepath.Join(orgDir, project))
	if exists {
		fmt.Printf("Project already cloned, copied %s to clipboard\n", color.HiGreenString(changeDirMsg))
		_ = clipboard.Write(clipboard.FmtText, []byte(changeDirMsg))
		return
	}

	if err := os.MkdirAll(orgDir, 0o755); err != nil {
		panic(err)
	}

	if err := clone(org, project, orgDir); err != nil {
		panic(err)
	}

	fmt.Printf("\nCopied \n%s\nto clipboard\n", color.HiGreenString(changeDirMsg))
	_ = clipboard.Write(clipboard.FmtText, []byte(changeDirMsg))
}

func checkIfExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func extractData(link string) (string, string) {
	if strings.HasPrefix(link, "git@") && strings.HasSuffix(link, ".git") {
		return handleSSHLink(link)
	}

	linkWithoutGitExtension := strings.TrimSuffix(link, ".git")
	if strings.HasPrefix(linkWithoutGitExtension, "github.com") {
		linkWithoutGitExtension = fmt.Sprintf("https://%s", linkWithoutGitExtension)
	}
	parsedUrl, err := url.ParseRequestURI(linkWithoutGitExtension)
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
