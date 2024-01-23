package build

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"variant/log"
)

func ModifyIconHash(inputFile, outputFile string, maxColorChange, count int) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	// 计算原始Hash
	originalHash := CalculateHash(img)
	log.Infof("Original PNG Hash: %s", originalHash)

	rand.NewSource(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		rgba := addAlphaChannel(img, color.RGBA{R: 255, G: 255, B: 255, A: 255})

		for y := rgba.Bounds().Min.Y; y < rgba.Bounds().Max.Y; y++ {
			for x := rgba.Bounds().Min.X; x < rgba.Bounds().Max.X; x++ {
				r, g, b, a := rgba.At(x, y).RGBA()

				rChange := rand.Intn(maxColorChange*2) - maxColorChange
				gChange := rand.Intn(maxColorChange*2) - maxColorChange
				bChange := rand.Intn(maxColorChange*2) - maxColorChange

				r = uint32(clamp(int(r>>8)+rChange, 0, 255))
				g = uint32(clamp(int(g>>8)+gChange, 0, 255))
				b = uint32(clamp(int(b>>8)+bChange, 0, 255))

				rgba.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a >> 8)})
			}
		}

		// 计算修改后的Hash
		modifiedHash := CalculateHash(rgba)
		log.Infof("Modified PNG Hash: %s", modifiedHash)

		// 生成文件名
		dir, iconFile := filepath.Split(outputFile)
		baseName := strings.TrimSuffix(iconFile, ".png")
		ext := ".png"
		fileName := filepath.Join(dir, fmt.Sprintf("%s-%d%s", baseName, i+1, ext))

		outFile, err := os.Create(fileName)
		if err != nil {
			return err
		}

		err = png.Encode(outFile, rgba)
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func addAlphaChannel(img image.Image, color color.RGBA) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, &image.Uniform{C: color}, image.Pt(0, 0), draw.Src)
	draw.Draw(rgba, b, img, b.Min, draw.Over)
	return rgba
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func CalculateHash(img image.Image) string {
	imageBytes := getImageBytes(img)
	hash := md5.Sum(imageBytes)
	return hex.EncodeToString(hash[:])
}

func getImageBytes(img image.Image) []byte {
	resizedImg := resize.Resize(100, 0, img, resize.Lanczos3)
	buf := new(strings.Builder)
	_ = png.Encode(buf, resizedImg)
	return []byte(buf.String())
}
