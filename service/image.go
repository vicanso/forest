// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"time"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/util"
)

var (
	fontPath string
)

const (
	captchaKeyPrefix = "captcha-"
)

type (
	// CaptchaInfo captcha info
	CaptchaInfo struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id,omitempty"`
	}
)

func init() {
	fontPath = config.GetString("resources.font")

}

// createCaptcha create captcha image
func createCaptcha(width, height int, text string) (img *image.RGBA, err error) {
	img = image.NewRGBA(image.Rect(0, 0, width, height))
	gc := draw2dimg.NewGraphicContext(img)
	draw2dkit.RoundedRectangle(gc, 0, 0, float64(width), float64(height), 0, 0)
	// 背景色设置为透明
	gc.SetFillColor(color.RGBA{255, 255, 255, 0})
	gc.Fill()

	gc.FillStroke()
	draw2d.SetFontFolder(fontPath)

	// Set the font luximbi.ttf
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})

	gc.SetFillColor(color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 120,
	})
	fontCount := len(text)
	offset := 10
	eachFontWidth := (width - 2*offset) / fontCount
	fontSize := float64(eachFontWidth) * 1.2
	for index, ch := range text {
		newFontSize := float64(rand.Int63n(40)+80) / 100 * fontSize
		gc.SetFontSize(newFontSize)
		angle := float64(rand.Int63n(20))/100 - 0.1
		offsetX := float64(eachFontWidth + index*eachFontWidth + int(rand.Int63n(10)) - 10)
		offsetY := float64(height) + float64(rand.Int63n(10)) - float64(15)
		if offsetY > float64(height) || offsetX < float64(height)-newFontSize {
			offsetY = float64(height)
		}
		gc.Rotate(angle)
		gc.FillStringAt(string(ch), offsetX, offsetY)
	}

	return
}

// GetCaptcha get captcha
func GetCaptcha() (info *CaptchaInfo, err error) {
	value := util.RandomDigit(4)
	img, err := createCaptcha(80, 40, value)
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, img)
	if err != nil {
		return
	}
	id := util.GenUlid()
	err = redisSrv.Set(captchaKeyPrefix+id, value, 2*time.Minute)
	if err != nil {
		return
	}
	info = &CaptchaInfo{
		Data: buffer.Bytes(),
		ID:   id,
	}
	return
}

// ValidateCaptcha validate the captch
func ValidateCaptcha(id, value string) (valid bool, err error) {
	data, err := redisSrv.GetAndDel(captchaKeyPrefix + id)
	if err != nil {
		return
	}
	valid = data == value
	return
}
