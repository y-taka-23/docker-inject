package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type injector struct {
}

func newInjector() *injector {
	return &injector{}
}

func (inj *injector) run(args []string) error {
	if len(args) != 2 {
		errors.New("two arguments required")
	}
	srcPath := args[0]
	container, tgtPath, err := inj.parseTarget(args[1])
	if err != nil {
		return err
	}
	fmt.Println("%s ==> &s:&s", srcPath, container, tgtPath)
	return nil
}

func (inj *injector) parseTarget(tgt string) (string, string, error) {
	strs := strings.Split(tgt, ":")
	if len(strs) != 2 {
		return "", "", fmt.Errorf("invalid target: %s", tgt)
	}
	return strs[0], strs[1], nil
}

func (inj *injector) inject(src, con, tgt string) error {
	// TODO revise with cmd.Stdin
	cmd := exec.Command(
		"docker", "exec", "-it", con,
		"/bin/bash", "-c", "'cat > "+tgt+"'", "<", src,
	)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (inj *injector) mkdir(con, tgt string) error {
	cmd := exec.Command("docker", "exec", "-it", con, "mkdir", "-p", tgt)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
