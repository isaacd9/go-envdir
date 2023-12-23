package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const code = 111

func usage() {
	fmt.Println("envdir: usage: envdir dir child")
	os.Exit(code)
}

func die(err error) {
	fmt.Println(err)
	os.Exit(code)
}

func envVar(name, val string) string {
	return name + "=" + val
}

func envv(environ []string, dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory: %w", err)
	}

	var extra [][2]string

	for _, file := range files {
		contents, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("unable to read file: %w", err)
		}
		name := file.Name()

		trimmed := strings.TrimRight(string(contents), " \t\n")
		extra = append(extra, [2]string{name, trimmed})
	}

	out := make([]string, 0, len(environ)+len(extra))

	for _, e := range environ {
		found := false
		for _, v := range extra {
			if strings.HasPrefix(e, v[0]) {
				found = true
			}
		}

		if !found {
			out = append(out, e)
		}
	}

	for _, e := range extra {
		if e[1] != "" {
			out = append(out, envVar(e[0], e[1]))
		}
	}

	return out, nil
}

func main() {
	runtime.LockOSThread()

	args := os.Args
	if len(args) < 3 {
		usage()
	}

	dir, exe, argv := args[1], args[2], args[2:]

	ev, err := envv(os.Environ(), dir)
	if err != nil {
		die(err)
		return
	}

	pathExe, err := exec.LookPath(exe)
	if err != nil {
		die(fmt.Errorf("unable to find executable in path: %w", err))
	}

	if err := syscall.Exec(pathExe, argv, ev); err != nil {
		die(fmt.Errorf("unable to exec: %w", err))
	}
}
