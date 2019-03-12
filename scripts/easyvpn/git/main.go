package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

var debug = false

func git(args ...string) {

	cmd := exec.Command("git", args...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if debug {
		fmt.Printf("Weird: %v", outb.String())
	}

	if err != nil {
		fmt.Println(errb.String())
	}
}

// Commit create a new commit
func Commit(files []string, msg string) {
	args := []string{"commit"}
	for _, file := range files {
		args = append(args, file)
	}
	args = append(args, "-m")
	args = append(args, msg)

	git(args...)
}

// Add create a new commit
func Add(files []string) {
	args := []string{"add"}
	for _, file := range files {
		args = append(args, file)
	}
	git(args...)
}

// Pull fetch from origin
func Pull() {
	args := []string{"pull"}
	git(args...)
}

// Push the branch on remote master branch
func Push() {
	args := []string{"push", "origin", "master"}
	git(args...)
}

// Rebase from origin
func Rebase() {
	args := []string{"rebase", "origin/master"}
	git(args...)
}
