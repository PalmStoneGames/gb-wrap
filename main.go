/*
Copyright 2015 Palm Stone Games, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/constabulary/gb"
	"github.com/constabulary/gb/cmd"
)

const DocUsage = `gb wrap is a tool for running GOPATH aware tools in a gb project

Usage:

	gb wrap [command] [arguments]

The command can be any tool available in PATH
`

var (
	projectRoot string
)

func main() {
	// Setup flags
	fs := flag.NewFlagSet("gb-wrap", flag.ExitOnError)
	fs.StringVar(&projectRoot, "R", os.Getenv("GB_PROJECT_DIR"), "set the project root")

	err := cmd.RunCommand(fs, &cmd.Command{
		Run: run,
	}, os.Getenv("GB_PROJECT_DIR"), "", os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func run(ctx *gb.Context, args []string) error {
	if len(args) == 0 {
		return errors.New(DocUsage)
	}

	// Build up the fake env
	env := cmd.MergeEnv(os.Environ(), map[string]string{
		"GOPATH": fmt.Sprintf("%s%s%s", projectRoot, string(os.PathListSeparator), path.Join(projectRoot, "vendor")),
	})

	app := exec.Command(args[0], args[1:]...)
	app.Stdin = os.Stdin
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	app.Env = env

	if err := app.Run(); err != nil {
		return fmt.Errorf("Failed to run command: %v", err)
	}

	return nil
}
