package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const defaultBranch string = "main"

func git(args ...string) (string, error) {

	cmd := exec.Command("git", args...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	fmt.Printf("Exec: %v\n", strings.Join(cmd.Args, " "))
	fmt.Printf("%v", outb.String())
	fmt.Printf("---\n")

	if err != nil {
		fmt.Println(errb.String())
	}
	return outb.String(), err
}

// Commit create a new commit
func Commit(files []string, msg string) {
	args := []string{"commit"}
	args = append(args, files...)
	args = append(args, "-m")
	args = append(args, msg)

	out, err := git(args...)
	if err != nil {
		fmt.Println(out)
		panic(err)
	}
}

// Add create a new commit
func Add(files []string) {
	args := []string{"add"}
	args = append(args, files...)
	out, err := git(args...)
	if err != nil {
		fmt.Println(out)
		panic(err)
	}
}

// Pull fetch from origin
func Pull() {
	args := []string{"pull"}
	out, err := git(args...)
	if err != nil {
		fmt.Println(out)
		panic(err)
	}
}

// getRepoOwner returns local branch
func getLocalBranch() string {
	args := []string{"rev-parse", "--abbrev-ref", "HEAD"}
	branch, _ := git(args...)
	fmt.Printf("Current git branch: %v\n", branch)
	return branch
}

// getRepoOwner returns current github owner
func getRepoOwner() (string, error) {
	var owner string
	args := []string{"config", "remote.origin.url"}
	url, _ := git(args...)
	if strings.HasPrefix(url, "git@github.com:") {
		url := strings.Split(url, ":")
		owner = strings.Split(url[1], "/")[0]
	} else if strings.HasPrefix(url, "https://github.com/") {
		owner = strings.Split(url, "/")[3]
	} else {
		err := fmt.Errorf("couldn't find current repository owner in %v", url)
		return "", err
	}
	fmt.Printf("Current repository owner: %v\n", owner)
	return owner, nil
}

// Push the branch on the principal branch
func Push() {
	branch := getLocalBranch()
	owner, err := getRepoOwner()
	if err != nil {
		panic(err)
	}

	args := []string{"push"}
	out, err := git(args...)
	if err != nil {
		fmt.Println(out)
		panic(err)
	}

	fmt.Printf("You can now open your Pull Request via \n\t https://github.com/jenkins-infra/docker-openvpn/compare/%v...%v:%v\n", defaultBranch, owner, branch)
}

// Rebase from origin
func Rebase() {
	args := []string{"rebase", fmt.Sprintf("origin/%v", defaultBranch)}
	out, err := git(args...)
	if err != nil {
		fmt.Println(out)
		panic(err)
	}
}
