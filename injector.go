package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type injector struct {
	stdout    io.Writer
	stderr    io.Writer
	source    string
	hostRoot  string
	container string
	contRoot  string
	isVerbose bool
}

func parseTarget(s string) (string, string, error) {
	strs := strings.Split(s, ":")
	if len(strs) != 2 {
		return "", "", fmt.Errorf("invalid target: %s", s)
	}
	return strs[0], strs[1], nil
}

func newInjector(o, e io.Writer, c *cli.Context) (*injector, error) {
	args := c.Args()
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
		stdout:    o,
		stderr:    e,
		source:    origin,
		hostRoot:  hRoot,
		container: cont,
		contRoot:  cRoot,
		isVerbose: c.Bool("verbose"),
	}, nil
}

func (inj *injector) run() error {
	return filepath.Walk(inj.source, inj.inject)
}

func (inj *injector) inject(curr string, fi os.FileInfo, e error) error {
	if e != nil {
		return e
	}
	rel, err := filepath.Rel(inj.hostRoot, curr)
	if err != nil {
		return err
	}
	tgt := path.Join(inj.contRoot, filepath.ToSlash(rel))
	if !fi.IsDir() {
		dir := path.Dir(tgt)
		if err := inj.injectDir(inj.container, dir); err != nil {
			return err
		}
		if err := inj.showProgress(rel, tgt); err != nil {
			return err
		}
		if err := inj.injectFile(curr, inj.container, tgt); err != nil {
			return err
		}
		return inj.changeMode(inj.container, tgt, fi.Mode())
	}
	fis, err := ioutil.ReadDir(curr)
	if err != nil {
		return err
	}
	if len(fis) == 0 {
		if err := inj.showProgress(rel, tgt); err != nil {
			return err
		}
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

func (inj *injector) showProgress(src, tgt string) error {
	if !inj.isVerbose {
		return nil
	}
	_, err := fmt.Fprintf(inj.stdout, "%s -> %s\n", src, tgt)
	return err
}
