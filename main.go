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

func goexec(in io.Reader, out, errw io.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("Pattern: go-envdir /path/to/env/dir some_prog")
	}

	path := args[0]
	progname := args[1]

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(progname)
	cmd.Env = env(path, files)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errw
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func env(path string, files []os.FileInfo) []string {
	env := make([]string, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := filepath.Join(path, file.Name())
		data, err := ioutil.ReadFile(name)
		if err != nil {
			continue
		}

		var builder strings.Builder
		builder.WriteString(file.Name())
		builder.WriteRune('=')
		builder.WriteString(string(data))
		env = append(env, builder.String())
	}

	return env
}

func main() {
	err := goexec(os.Stdin, os.Stdout, os.Stderr, os.Args[1:])
	if err != nil {
		panic(err)
	}
}
