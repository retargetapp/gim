package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

const CALL_GIM_GIT_HOOK_CODE = "$GOPATH/bin/gim sync"

func installCmd(c *cli.Context) error {
	fmt.Println("Install Git hooks...")

	hp, err := filepath.Abs("./.git/hooks")
	if err != nil {
		fmt.Println("Unable to check Git hooks directory: " + err.Error())
	}

	err = installGitHook(hp + string(os.PathSeparator) + "post-merge")
	if err != nil {
		return cli.NewExitError("Unable to insatll Git post-merge hook: "+err.Error(), 1)
	}
	fmt.Println("Post-merge Git hook installed")

	installGitHook(hp + string(os.PathSeparator) + "post-checkout")
	if err != nil {
		return cli.NewExitError("Unable to insatll Git post-checkout hook: "+err.Error(), 1)
	}
	fmt.Println("Post-checkout Git hook installed")

	return nil
}

func installGitHook(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	i, sc := false, bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Text() == CALL_GIM_GIT_HOOK_CODE {
			i = true
			break
		}
	}

	if i {
		return nil
	}

	_, err = f.WriteString(CALL_GIM_GIT_HOOK_CODE + "\n")
	return err
}
