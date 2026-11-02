package main

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"testing"
)

func syntheticPortrait() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 400, 520))
	for y := 0; y < 520; y++ {
		for x := 0; x < 400; x++ {
			shade := uint8(232 - y/12)
			img.SetRGBA(x, y, color.RGBA{R: shade, G: shade - 6, B: uint8(int(shade) - 20), A: 0xFF})
		}
	}
	for y := 110; y < 300; y++ {
		for x := 130; x < 270; x++ {
			dx := float64(x-200) / 70
			dy := float64(y-190) / 95
			if dx*dx+dy*dy <= 1 {
				img.SetRGBA(x, y, color.RGBA{R: 42, G: 33, B: 30, A: 0xFF})
			}
		}
	}
	for y := 300; y < 470; y++ {
		width := 40 + (y-300)/2
		for x := 200 - width; x < 200+width; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 58, G: 44, B: 38, A: 0xFF})
		}
	}
	return img
}

func TestGenerateShadowPlayProducesAnimatedGIF(t *testing.T) {
	gifData, posterData, frames, err := GenerateShadowPlay(syntheticPortrait(), "亮相", true)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if frames != playFrames {
		t.Fatalf("expected %d frames, got %d", playFrames, frames)
	}
	if len(posterData) == 0 {
		t.Fatal("poster data is empty")
	}

	anim, err := gif.DecodeAll(bytes.NewReader(gifData))
	if err != nil {
		t.Fatalf("decode gif failed: %v", err)
	}
	if len(anim.Image) != playFrames {
		t.Fatalf("expected %d gif frames, got %d", playFrames, len(anim.Image))
	}
	bounds := anim.Image[0].Bounds()
	if bounds.Dx() != playWidth || bounds.Dy() != playHeight {
		t.Fatalf("unexpected frame size %dx%d", bounds.Dx(), bounds.Dy())
	}
	if anim.LoopCount != 0 {
		t.Fatalf("expected looping animation, got loop count %d", anim.LoopCount)
	}

	first := anim.Image[0].Pix
	last := anim.Image[len(anim.Image)-1].Pix
	differences := 0
	for index := 0; index < len(first) && index < len(last); index += 97 {
		if first[index] != last[index] {
			differences++
		}
	}
	if differences < 20 {
		t.Fatalf("animation appears static, only %d sampled differences", differences)
	}
	checker := anim.Image[len(anim.Image)/2]
	if len(checker.Pix) == 0 {
		t.Fatal("middle frame is empty")
	}

	lit := 0
	ink := 0
	screenPixels := 0
	for y := screenTop + 10; y < screenFloor-10; y++ {
		for x := screenLeft + 10; x < screenRight-10; x++ {
			r, g, b, _ := checker.At(x, y).RGBA()
			lum := luminance(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			screenPixels++
			if lum > 0.6 {
				lit++
			}
			if lum < 0.35 {
				ink++
			}
		}
	}
	if lit*100/screenPixels < 55 {
		t.Fatalf("screen looks unlit, only %d%% bright pixels", lit*100/screenPixels)
	}
	if ink*100/screenPixels < 3 {
		t.Fatalf("puppet silhouette missing, only %d%% dark pixels", ink*100/screenPixels)
	}
}

func TestSeedImageGeneratesPreviewPlay(t *testing.T) {
	source, err := os.Open("static/seed/art-puppet.png")
	if err != nil {
		t.Skipf("seed image unavailable: %v", err)
	}
	defer source.Close()
	src, _, err := image.Decode(source)
	if err != nil {
		t.Fatalf("decode seed image failed: %v", err)
	}
	gifData, _, _, err := GenerateShadowPlay(src, "亮相", true)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if err := os.WriteFile("static/seed/demo-shadow-play.gif", gifData, 0o644); err != nil {
		t.Fatalf("write preview failed: %v", err)
	}
	inkData, _, _, err := GenerateShadowPlay(src, "亮相", false)
	if err != nil {
		t.Fatalf("ink generate failed: %v", err)
	}
	if err := os.WriteFile("static/seed/demo-shadow-play-ink.gif", inkData, 0o644); err != nil {
		t.Fatalf("write ink preview failed: %v", err)
	}
}

func middleFrame(t *testing.T, colored bool) image.Image {
	t.Helper()
	gifData, _, _, err := GenerateShadowPlay(syntheticPortrait(), "亮相", colored)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	anim, err := gif.DecodeAll(bytes.NewReader(gifData))
	if err != nil {
		t.Fatalf("decode gif failed: %v", err)
	}
	return anim.Image[len(anim.Image)/2]
}

func TestPaletteFollowsSourceColor(t *testing.T) {
	cases := []struct {
		name        string
		r, g, b     float64
		dominant    string
		expectAbove func(color.RGBA) bool
	}{
		{"blue source picks indigo", 0.20, 0.35, 0.75, "blue", func(c color.RGBA) bool { return c.B > c.R }},
		{"green source picks jade", 0.30, 0.60, 0.30, "green", func(c color.RGBA) bool { return c.G > c.R && c.G > c.B }},
		{"red source picks cinnabar", 0.70, 0.25, 0.20, "red", func(c color.RGBA) bool { return c.R > c.B }},
	}
	for _, item := range cases {
		pal := paletteForLayer(puppetLayer{avgR: item.r, avgG: item.g, avgB: item.b})
		if !item.expectAbove(pal.primary) {
			t.Fatalf("%s: got primary %#v", item.name, pal.primary)
		}
	}
}

func TestColorModeDiffersFromInkMode(t *testing.T) {
	colored := middleFrame(t, true)
	ink := middleFrame(t, false)
	bounds := colored.Bounds()
	different := 0
	sampled := 0
	for y := screenTop + 10; y < screenFloor-10; y += 3 {
		for x := screenLeft + 10; x < screenRight-10; x += 3 {
			cr, cg, cb, _ := colored.At(x, y).RGBA()
			ir, ig, ib, _ := ink.At(x, y).RGBA()
			sampled++
			delta := math.Abs(float64(cr>>8)-float64(ir>>8)) +
				math.Abs(float64(cg>>8)-float64(ig>>8)) +
				math.Abs(float64(cb>>8)-float64(ib>>8))
			if delta > 30 {
				different++
			}
		}
	}
	if bounds.Dx() == 0 {
		t.Fatal("empty frame")
	}
	if different*100/sampled < 15 {
		t.Fatalf("color mode barely changes the frame: %d%% of sampled pixels differ", different*100/sampled)
	}
}
