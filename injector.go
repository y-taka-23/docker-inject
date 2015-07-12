package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type injector struct {
	stderr    io.Writer
	source    string
	hostRoot  string
	container string
	contRoot  string
}

func parseTarget(s string) (string, string, error) {
	strs := strings.Split(s, ":")
	if len(strs) != 2 {
		return "", "", fmt.Errorf("invalid target: %s", s)
	}
	return strs[0], strs[1], nil
}

func newInjector(w io.Writer, args []string) (*injector, error) {
	if len(args) != 2 {
		return nil, errors.New("two arguments required")
	}
	origin, err := filepath.Abs(args[0])
	if err != nil {
		return nil, err
	}
	hRoot := filepath.Dir(origin)
	cont, cRoot, err := parseTarget(args[1])
	if err != nil {
		return nil, err
	}
	return &injector{
		stderr:    w,
		source:    origin,
		hostRoot:  hRoot,
		container: cont,
		contRoot:  cRoot,
	}, nil
}

func (inj *injector) run() error {
	return filepath.Walk(inj.source, inj.inject)
}

func (inj *injector) inject(path string, fi os.FileInfo, e error) error {
	if e != nil {
		return e
	}
	rel, err := filepath.Rel(inj.hostRoot, path)
	if err != nil {
		return err
	}
	tgt := filepath.Join(inj.contRoot, rel)
	if !fi.IsDir() {
		dir := filepath.Dir(tgt)
		if err := inj.injectDir(inj.container, dir); err != nil {
			return err
		}
		if err := inj.injectFile(path, inj.container, tgt); err != nil {
			return err
		}
		return inj.changeMode(inj.container, tgt, fi.Mode())
	}
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if len(fis) == 0 {
		return inj.injectDir(inj.container, tgt)
	}
	return nil
}

func (inj *injector) injectFile(src, con, tgt string) error {
	cmd := exec.Command(
		"docker", "exec", "-i", con,
		"/bin/bash", "-c", "cat > "+tgt,
	)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	cmd.Stdin = f
	cmd.Stderr = inj.stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (inj *injector) changeMode(con, tgt string, fm os.FileMode) error {
	cmd := exec.Command(
		"docker", "exec", con, "chmod", fmt.Sprintf("%o", fm), tgt,
	)
	cmd.Stderr = inj.stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (inj *injector) injectDir(con, dir string) error {
	cmd := exec.Command("docker", "exec", con, "mkdir", "-p", dir)
	cmd.Stderr = inj.stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
