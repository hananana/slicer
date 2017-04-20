package main

import (
	"fmt"
	"github.com/oliamb/cutter"
	"gopkg.in/urfave/cli.v2"
	"image"
	"image/png"
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

func slice(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		dotIndex := strings.LastIndex(path, ".")
		ext := path[dotIndex:]
		if ext == ".png" {
			doSlice(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func doSlice(path string) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:   32,
		Height:  48,
		Options: cutter.Copy,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	if !exists("out") {
		if err := os.Mkdir("out", 0777); err != nil {
			fmt.Println(err)
		}
	}

	out, _ := os.Create("out/test.png")
	_ = png.Encode(out, croppedImg)
	defer out.Close()
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
