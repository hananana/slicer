package main

import (
	"fmt"
	"github.com/oliamb/cutter"
	"gopkg.in/urfave/cli.v2"
	"image"
	"image/draw"
	// 	"image/png"
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
			images := doSlice(path)
			combine(images)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func doSlice(path string) []image.Image {
	file, _ := os.Open(path)
	defer file.Close()

	img, _, _ := image.Decode(file)

	croppedImages := make([]image.Image, 12)

	for y := 0; y < 4; y++ {
		for x := 0; x < 3; i++ {
			croppedImg, _ := cutter.Crop(img, cutter.Config{
				Width:   32,
				Height:  48,
				Anchor:  image.Point{32 * x, 48 * y},
				Options: cutter.Copy,
			})
			croppedImages[i] = croppedImg
		}
	}

	return croppedImages
}

func combine(images []image.Image) {
	if !exists("out") {
		if err := os.Mkdir("out", 0777); err != nil {
			fmt.Println(err)
		}
	}
	// 	dp := image.Point{images[0].Bounds().Dx, 0}
	// 	r := image.Rect(0, 0, 64, 48)
	// 	rgba := image.NewRGBA(r)
	// 	draw.Draw(rgba, images[0].Bounds(), images[0], image.Point{0, 0}, draw.Src)
	// 	draw.Draw(rgba, image.Rect(32, 0, 32, 48), images[1], image.Point{32, 0}, draw.Src)
	// 	draw.Draw
	// 	out, _ := os.Create("out/test.png")
	// 	_ = png.Encode(out, croppedImg)
	// 	defer out.Close()
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
