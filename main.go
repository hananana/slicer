package main

import (
	"fmt"
	"github.com/oliamb/cutter"
	"gopkg.in/urfave/cli.v2"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const baseWidth = 32
const baseHeight = 48

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

	croppedImages := make([]image.Image, 4)

	for i := 0; i < 4; i++ {
		croppedImg, _ := cutter.Crop(img, cutter.Config{
			Width:  baseWidth,
			Height: baseHeight * 4,
			Anchor: image.Point{baseWidth * i, 0},
			Mode:   cutter.TopLeft,
		})
		croppedImages[i] = croppedImg
	}

	return croppedImages
}

func combine(images []image.Image) {
	if !exists("out") {
		if err := os.Mkdir("out", 0777); err != nil {
			fmt.Println(err)
		}
	}

	tempDir := "slicerTemp"
	if !exists(tempDir) {
		if err := os.Mkdir(tempDir, 0777); err != nil {
			fmt.Println(err)
		}
	}

	for index, image := range images {
		out, err := os.Create(tempDir + "/" + strconv.Itoa(index) + ".png")
		if err != nil {
			println(err)
		}
		err = png.Encode(out, image)
		if err != nil {
			println(err)
		}
	}

	newImages := make([]image.Image, 4)
	for i := 0; i < 4; i++ {
		file, _ := os.Open(tempDir + "/" + strconv.Itoa(i) + ".png")
		img, _, _ := image.Decode(file)
		newImages[i] = img
	}

	baseRect := image.Rect(0, 0, baseWidth*4, baseHeight*4)
	rgba := image.NewRGBA(baseRect)
	stay := newImages[0]
	right := newImages[1]
	left := newImages[3]

	draw.Draw(rgba, left.Bounds(), left, image.ZP, draw.Src)
	draw.Draw(rgba, stay.Bounds().Add(image.Pt(baseWidth, 0)), stay, image.ZP, draw.Src)
	draw.Draw(rgba, right.Bounds().Add(image.Pt(baseWidth*2, 0)), right, image.ZP, draw.Src)

	out, err := os.Create("out/test.png")
	defer out.Close()
	if err != nil {
		println(err)
	}
	err = png.Encode(out, rgba)
	if err != nil {
		println(err)
	}

	os.RemoveAll(tempDir)
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
