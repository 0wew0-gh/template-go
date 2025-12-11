package tool

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

/**
* 入参： JPG 图片文件的二进制数据
* 出参：JPG 图片的宽和高
**/
func GetJPGWidthHeight(imgBytes []byte) (int, int) {
	var offset int
	imgByteLen := len(imgBytes)
	for i := 0; i < imgByteLen-1; i++ {
		if imgBytes[i] != 0xff {
			continue
		}
		if imgBytes[i+1] == 0xC0 || imgBytes[i+1] == 0xC1 || imgBytes[i+1] == 0xC2 {
			offset = i
			break
		}
	}
	offset += 5
	if offset >= imgByteLen {
		return 0, 0
	}
	height := int(imgBytes[offset])<<8 + int(imgBytes[offset+1])
	width := int(imgBytes[offset+2])<<8 + int(imgBytes[offset+3])
	return width, height
}

// 获取 PNG 图片的宽高
func GetPNGWidthHeight(imgBytes []byte) (int, int) {
	pngHeader := "\x89PNG\r\n\x1a\n"
	if string(imgBytes[:len(pngHeader)]) != pngHeader {
		return 0, 0
	}
	offset := 12
	if string(imgBytes[offset:offset+4]) != "IHDR" {
		return 0, 0
	}
	offset += 4
	width := int(binary.BigEndian.Uint32(imgBytes[offset : offset+4]))
	height := int(binary.BigEndian.Uint32(imgBytes[offset+4 : offset+8]))
	return width, height
}

// 获取 GIF 图片的宽高
func GetGifWidthHeight(imgBytes []byte) (int, int) {
	version := string(imgBytes[:6])
	if version != "GIF87a" && version != "GIF89a" {
		return 0, 0
	}
	width := int(imgBytes[6]) + int(imgBytes[7])<<8
	height := int(imgBytes[8]) + int(imgBytes[9])<<8
	return width, height
}

// 获取 WEBP 图片的宽高
func GetWEBPWidthHeight(imgBytes []byte) (int, int) {
	offset := 26
	width := int(imgBytes[offset+1]&0x3f)<<8 | int(imgBytes[offset])
	height := int(imgBytes[offset+3]&0x3f)<<8 | int(imgBytes[offset+2])
	return width, height
}

// 获取 BMP 图片的宽高
func GetBmpWidthHeight(imgBytes []byte) (int, int) {
	if string(imgBytes[:2]) != "BM" {
		return 0, 0
	}
	width := int(binary.LittleEndian.Uint32(imgBytes[18:22]))
	height := int(int32(binary.LittleEndian.Uint32(imgBytes[22:26])))
	if height < 0 {
		height = -height
	}
	return width, height
}

// mediaDim image width and height
type mediaDim struct {
	Width  float64 `json:"width"`  //宽
	Height float64 `json:"height"` //高
}

// 解析图片的宽高信息
func DecodeImageWidthHeight(imgBytes []byte, fileType string) (*mediaDim, error) {
	var (
		imgConf image.Config
		err     error
	)
	switch strings.ToLower(fileType) {
	case "jpg", "jpeg":
		imgConf, err = jpeg.DecodeConfig(bytes.NewReader(imgBytes))
	case "webp":
		imgConf, err = webp.DecodeConfig(bytes.NewReader(imgBytes))
	case "png":
		imgConf, err = png.DecodeConfig(bytes.NewReader(imgBytes))
	case "tif", "tiff":
		imgConf, err = tiff.DecodeConfig(bytes.NewReader(imgBytes))
	case "gif":
		imgConf, err = gif.DecodeConfig(bytes.NewReader(imgBytes))
	case "bmp":
		imgConf, err = bmp.DecodeConfig(bytes.NewReader(imgBytes))
	default:
		return nil, errors.New("unknown file type")
	}
	if err != nil {
		return nil, err
	}
	return &mediaDim{
		Width:  float64(imgConf.Width),
		Height: float64(imgConf.Height),
	}, nil
}

// 计算图片中已存在的图片数量
func CalcExist(imageExist []bool) int {
	Exists := 0
	for _, v := range imageExist {
		if !v {
			Exists++
		}
	}
	return Exists
}
