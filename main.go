package main

import (
	"gopkg.in/urfave/cli.v2"
	// 	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		if !check(c) {
			return nil
		}

		slice(c.Args().Get(0))

		return nil
	}

	app.Run(os.Args)
}

func slice(path string) {
}

func check(c *cli.Context) bool {
	if c.NArg() == 0 {
		println("args must one")
		return false
	}

	if c.NArg() > 1 {
		println("args must one")
		return false
	}

	targetDir := c.Args().Get(0)

	if !exists(targetDir) {
		println("not exist dir")
		return false
	}

	result := false

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		dotIndex := strings.LastIndex(path, ".")
		ext := path[dotIndex:]
		if ext == ".png" {
			result = true
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	if result == false {
		fmt.Println("no .png files")
	}

	return result
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
