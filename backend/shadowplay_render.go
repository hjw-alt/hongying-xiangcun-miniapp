package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"math"
)

const (
	playWidth   = 320
	playHeight  = 420
	playFrames  = 30
	playDelay   = 10
	puppetBoxW  = 180
	puppetBoxH  = 240
	screenLeft  = 40
	screenTop   = 54
	screenRight = 280
	screenFloor = 330
)

var (
	colorInk     = color.RGBA{R: 0x2A, G: 0x17, B: 0x12, A: 0xFF}
	colorInkSoft = color.RGBA{R: 0x6B, G: 0x33, B: 0x22, A: 0xFF}
	colorPaperHi = color.RGBA{R: 0xF6, G: 0xE3, B: 0xBA, A: 0xFF}
	colorFrame   = color.RGBA{R: 0x2E, G: 0x1B, B: 0x15, A: 0xFF}
	colorCurtain = color.RGBA{R: 0x7E, G: 0x24, B: 0x19, A: 0xFF}
	colorGold    = color.RGBA{R: 0xC3, G: 0x9A, B: 0x4B, A: 0xFF}
	colorGlow    = color.RGBA{R: 0xFF, G: 0xD9, B: 0x8A, A: 0xFF}
)

type puppetLayer struct {
	width  int
	height int
	alpha  []float64
	tone   []float64
	avgR   float64
	avgG   float64
	avgB   float64
}

type shadowPalette struct {
	primary   color.RGBA
	secondary color.RGBA
	highlight color.RGBA
	light     color.RGBA
}

var traditionalPalettes = []struct {
	hue   float64
	names shadowPalette
}{
	{8, shadowPalette{
		primary:   color.RGBA{R: 0xC0, G: 0x39, B: 0x2B, A: 0xFF},
		secondary: color.RGBA{R: 0xD2, G: 0x73, B: 0x8B, A: 0xFF},
		highlight: color.RGBA{R: 0xE8, G: 0xA3, B: 0x3D, A: 0xFF},
		light:     color.RGBA{R: 0xFF, G: 0xC9, B: 0x8A, A: 0xFF},
	}},
	{38, shadowPalette{
		primary:   color.RGBA{R: 0xE0, G: 0xA3, B: 0x3E, A: 0xFF},
		secondary: color.RGBA{R: 0xC0, G: 0x39, B: 0x2B, A: 0xFF},
		highlight: color.RGBA{R: 0xF6, G: 0xE3, B: 0xBA, A: 0xFF},
		light:     color.RGBA{R: 0xFF, G: 0xD9, B: 0x8A, A: 0xFF},
	}},
	{90, shadowPalette{
		primary:   color.RGBA{R: 0x4E, G: 0x7A, B: 0x57, A: 0xFF},
		secondary: color.RGBA{R: 0x8A, G: 0xA5, B: 0x62, A: 0xFF},
		highlight: color.RGBA{R: 0xE8, G: 0xA3, B: 0x3D, A: 0xFF},
		light:     color.RGBA{R: 0xD6, G: 0xE8, B: 0xB4, A: 0xFF},
	}},
	{190, shadowPalette{
		primary:   color.RGBA{R: 0x2F, G: 0x4A, B: 0x6B, A: 0xFF},
		secondary: color.RGBA{R: 0x4E, G: 0x7A, B: 0x57, A: 0xFF},
		highlight: color.RGBA{R: 0xE8, G: 0xA3, B: 0x3D, A: 0xFF},
		light:     color.RGBA{R: 0xA8, G: 0xC8, B: 0xE0, A: 0xFF},
	}},
	{285, shadowPalette{
		primary:   color.RGBA{R: 0xC2, G: 0x56, B: 0x6F, A: 0xFF},
		secondary: color.RGBA{R: 0x6B, G: 0x4B, B: 0x8A, A: 0xFF},
		highlight: color.RGBA{R: 0xE8, G: 0xA3, B: 0x3D, A: 0xFF},
		light:     color.RGBA{R: 0xF0, G: 0xB8, B: 0xC8, A: 0xFF},
	}},
}

func hueOf(r, g, b float64) float64 {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	if max == min {
		return 0
	}
	delta := max - min
	var hue float64
	switch max {
	case r:
		hue = math.Mod((g-b)/delta, 6)
	case g:
		hue = (b-r)/delta + 2
	default:
		hue = (r-g)/delta + 4
	}
	hue *= 60
	if hue < 0 {
		hue += 360
	}
	return hue
}

func paletteForLayer(layer puppetLayer) shadowPalette {
	hue := hueOf(layer.avgR, layer.avgG, layer.avgB)
	selected := traditionalPalettes[0].names
	best := 360.0
	for _, entry := range traditionalPalettes {
		diff := math.Abs(hue - entry.hue)
		if diff > 180 {
			diff = 360 - diff
		}
		if diff < best {
			best = diff
			selected = entry.names
		}
	}
	return selected
}

func clampUnit(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func luminance(r, g, b uint8) float64 {
	return (0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 255
}

func clamp255(value float64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

func centerCrop(src image.Image, aspect float64) image.Rectangle {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	cropW := width
	cropH := int(float64(width) / aspect)
	if cropH > height {
		cropH = height
		cropW = int(float64(height) * aspect)
	}
	offsetX := bounds.Min.X + (width-cropW)/2
	offsetY := bounds.Min.Y + (height-cropH)/2
	return image.Rect(offsetX, offsetY, offsetX+cropW, offsetY+cropH)
}

func resizeBilinear(src image.Image, width, height int, area image.Rectangle) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	srcW := float64(area.Dx())
	srcH := float64(area.Dy())
	for y := 0; y < height; y++ {
		sy := (float64(y) + 0.5) * srcH / float64(height)
		y0 := int(math.Floor(sy - 0.5))
		fy := sy - 0.5 - float64(y0)
		if y0 < 0 {
			y0 = 0
			fy = 0
		}
		y1 := y0 + 1
		if y1 >= area.Dy() {
			y1 = area.Dy() - 1
		}
		for x := 0; x < width; x++ {
			sx := (float64(x) + 0.5) * srcW / float64(width)
			x0 := int(math.Floor(sx - 0.5))
			fx := sx - 0.5 - float64(x0)
			if x0 < 0 {
				x0 = 0
				fx = 0
			}
			x1 := x0 + 1
			if x1 >= area.Dx() {
				x1 = area.Dx() - 1
			}
			r00, g00, b00, _ := src.At(area.Min.X+x0, area.Min.Y+y0).RGBA()
			r10, g10, b10, _ := src.At(area.Min.X+x1, area.Min.Y+y0).RGBA()
			r01, g01, b01, _ := src.At(area.Min.X+x0, area.Min.Y+y1).RGBA()
			r11, g11, b11, _ := src.At(area.Min.X+x1, area.Min.Y+y1).RGBA()
			mix := func(v00, v10, v01, v11 uint32) uint8 {
				top := float64(v00)*(1-fx) + float64(v10)*fx
				bottom := float64(v01)*(1-fx) + float64(v11)*fx
				return uint8(clampUnit((top*(1-fy)+bottom*fy)/65535) * 255)
			}
			dst.SetRGBA(x, y, color.RGBA{
				R: mix(r00, r10, r01, r11),
				G: mix(g00, g10, g01, g11),
				B: mix(b00, b10, b01, b11),
				A: 0xFF,
			})
		}
	}
	return dst
}

func buildPuppetLayer(src image.Image) puppetLayer {
	area := centerCrop(src, float64(puppetBoxW)/float64(puppetBoxH))
	small := resizeBilinear(src, puppetBoxW, puppetBoxH, area)

	lum := make([]float64, puppetBoxW*puppetBoxH)
	total := 0.0
	for y := 0; y < puppetBoxH; y++ {
		for x := 0; x < puppetBoxW; x++ {
			c := small.RGBAAt(x, y)
			value := luminance(c.R, c.G, c.B)
			lum[y*puppetBoxW+x] = value
			total += value
		}
	}
	mean := total / float64(len(lum))
	threshold := mean * 0.94
	if threshold < 0.18 {
		threshold = 0.18
	}
	if threshold > 0.72 {
		threshold = 0.72
	}
	lower := threshold * 0.74
	upper := threshold * 1.06
	if upper-lower < 0.05 {
		upper = lower + 0.05
	}

	layer := puppetLayer{width: puppetBoxW, height: puppetBoxH}
	layer.alpha = make([]float64, len(lum))
	layer.tone = make([]float64, len(lum))
	sumR := 0.0
	sumG := 0.0
	sumB := 0.0
	subject := 0
	for index, value := range lum {
		alpha := clampUnit((upper - value) / (upper - lower))
		tone := clampUnit(value / threshold)
		layer.alpha[index] = alpha
		layer.tone[index] = tone
		if alpha > 0.6 {
			c := small.RGBAAt(index%puppetBoxW, index/puppetBoxW)
			sumR += float64(c.R) / 255
			sumG += float64(c.G) / 255
			sumB += float64(c.B) / 255
			subject++
		}
	}
	if subject > 0 {
		layer.avgR = sumR / float64(subject)
		layer.avgG = sumG / float64(subject)
		layer.avgB = sumB / float64(subject)
	} else {
		layer.avgR, layer.avgG, layer.avgB = 0.6, 0.35, 0.25
	}
	return layer
}

func blendPixel(img *image.RGBA, x, y int, c color.RGBA, alpha float64) {
	if alpha <= 0 || x < 0 || y < 0 || x >= img.Rect.Dx() || y >= img.Rect.Dy() {
		return
	}
	alpha = clampUnit(alpha)
	base := img.RGBAAt(x, y)
	mix := func(dst, src uint8) uint8 {
		return uint8(clampUnit((float64(dst)*(1-alpha)+float64(src)*alpha)/255) * 255)
	}
	img.SetRGBA(x, y, color.RGBA{
		R: mix(base.R, c.R),
		G: mix(base.G, c.G),
		B: mix(base.B, c.B),
		A: 0xFF,
	})
}

func drawLine(img *image.RGBA, x0, y0, x1, y1, thickness float64, c color.RGBA, alpha float64) {
	length := math.Hypot(x1-x0, y1-y0)
	if length == 0 {
		return
	}
	steps := int(length * 2)
	radius := thickness / 2
	for step := 0; step <= steps; step++ {
		t := float64(step) / float64(steps)
		px := x0 + (x1-x0)*t
		py := y0 + (y1-y0)*t
		for dy := -int(radius) - 1; dy <= int(radius)+1; dy++ {
			for dx := -int(radius) - 1; dx <= int(radius)+1; dx++ {
				dist := math.Hypot(float64(dx), float64(dy))
				if dist > radius {
					continue
				}
				edge := 1.0
				if dist > radius-1 {
					edge = radius - dist
				}
				blendPixel(img, int(px)+dx, int(py)+dy, c, alpha*edge)
			}
		}
	}
}

func fillRect(img *image.RGBA, rect image.Rectangle, c color.RGBA, alpha float64) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			blendPixel(img, x, y, c, alpha)
		}
	}
}

func mixColor(a, b color.RGBA, t float64) color.RGBA {
	t = clampUnit(t)
	return color.RGBA{
		R: clamp255(float64(a.R)*(1-t) + float64(b.R)*t),
		G: clamp255(float64(a.G)*(1-t) + float64(b.G)*t),
		B: clamp255(float64(a.B)*(1-t) + float64(b.B)*t),
		A: 0xFF,
	}
}

func drawGlow(img *image.RGBA, cx, cy, radius, alpha float64, glow color.RGBA) {
	startX := int(cx - radius)
	endX := int(cx + radius)
	startY := int(cy - radius)
	endY := int(cy + radius)
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			dist := math.Hypot(float64(x)-cx, float64(y)-cy) / radius
			if dist >= 1 {
				continue
			}
			falloff := math.Pow(1-dist, 2.2)
			blendPixel(img, x, y, glow, alpha*falloff)
		}
	}
}

func drawStage(img *image.RGBA, open, flicker float64, pal *shadowPalette) {
	fillRect(img, image.Rect(0, 0, playWidth, playHeight), colorFrame, 1)

	for y := screenTop; y < screenFloor; y++ {
		t := float64(y-screenTop) / float64(screenFloor-screenTop)
		base := color.RGBA{
			R: uint8(clampUnit(0.96-0.18*t) * 255),
			G: uint8(clampUnit(0.87-0.24*t) * 255),
			B: uint8(clampUnit(0.70-0.28*t) * 255),
			A: 0xFF,
		}
		if pal != nil {
			base = mixColor(base, pal.light, 0.22)
		}
		fillRect(img, image.Rect(screenLeft, y, screenRight, y+1), base, 1)
	}

	glow := colorGlow
	if pal != nil {
		glow = pal.light
	}
	drawGlow(img, float64(playWidth)/2, 120, 132, 0.5*flicker, glow)
	drawGlow(img, float64(playWidth)/2, 300, 150, 0.18*flicker, glow)

	for index := 0; index < 4; index++ {
		y := screenTop + 14 + index*80
		blendPixel(img, screenLeft+8, y, colorGold, 0.18)
	}

	curtainWidth := 96 * (1 - 0.7*open)
	for y := 0; y < playHeight; y++ {
		fold := 0.82 + 0.18*math.Abs(math.Sin(float64(y)*0.06))
		fillRect(img, image.Rect(0, y, int(curtainWidth), y+1), colorCurtain, fold)
		fillRect(img, image.Rect(playWidth-int(curtainWidth), y, playWidth, y+1), colorCurtain, fold)
	}

	fillRect(img, image.Rect(screenLeft-6, screenTop-6, screenRight+6, screenTop), colorGold, 0.85)
	fillRect(img, image.Rect(screenLeft-6, screenFloor, screenRight+6, screenFloor+6), colorGold, 0.85)
	fillRect(img, image.Rect(screenLeft-6, screenTop, screenLeft, screenFloor), colorGold, 0.85)
	fillRect(img, image.Rect(screenRight, screenTop, screenRight+6, screenFloor), colorGold, 0.85)

	fillRect(img, image.Rect(0, screenFloor+6, playWidth, screenFloor+14), color.RGBA{R: 0x3A, G: 0x1E, B: 0x14, A: 0xFF}, 0.9)
}

func puppetTransform(preset string, progress float64) (float64, float64, float64, float64) {
	centerX := float64(playWidth) / 2
	centerY := 232.0
	rotation := 0.0
	scale := 1.0

	switch preset {
	case "舞剑":
		rotation = math.Sin(progress*math.Pi*4) * 0.10
		scale = 1 + math.Sin(progress*math.Pi*2)*0.035
		centerY = 236 + math.Sin(progress*math.Pi*4)*6
	case "赶路":
		centerX = 66 + progress*(float64(playWidth)-132)
		centerY = 238 + math.Sin(progress*math.Pi*6)*7
		rotation = math.Sin(progress*math.Pi*6) * 0.06
	default:
		if progress < 0.38 {
			t := progress / 0.38
			ease := 1 - math.Pow(1-t, 3)
			centerX = -30 + ease*(float64(playWidth)/2+30)
			rotation = (1 - ease) * 0.32
			centerY = 238 - (1-ease)*16
		} else {
			rotation = math.Sin((progress-0.38)*math.Pi*3.2) * 0.05
			centerY = 232 + math.Sin((progress-0.38)*math.Pi*2)*5
		}
	}
	return centerX, centerY, rotation, scale
}

func isSilhouetteEdge(layer puppetLayer, x, y int) bool {
	if layer.alpha[y*layer.width+x] < 0.35 {
		return false
	}
	offsets := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, offset := range offsets {
		nx := x + offset[0]
		ny := y + offset[1]
		if nx < 0 || ny < 0 || nx >= layer.width || ny >= layer.height {
			return true
		}
		if layer.alpha[ny*layer.width+nx] < 0.35 {
			return true
		}
	}
	return false
}

func drawPuppet(img *image.RGBA, layer puppetLayer, centerX, centerY, rotation, scale, opacity float64, pal *shadowPalette) {
	cosA := math.Cos(-rotation)
	sinA := math.Sin(-rotation)
	halfW := float64(layer.width) * scale / 2
	halfH := float64(layer.height) * scale / 2

	startX := int(centerX - halfW - 2)
	endX := int(centerX + halfW + 2)
	startY := int(centerY - halfH - 2)
	endY := int(centerY + halfH + 2)

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			srcX := (dx*cosA-sinA*dy)/scale + float64(layer.width)/2
			srcY := (dx*sinA+cosA*dy)/scale + float64(layer.height)/2
			sampleX := int(srcX)
			sampleY := int(srcY)
			if sampleX < 0 || sampleY < 0 || sampleX >= layer.width || sampleY >= layer.height {
				continue
			}
			index := sampleY*layer.width + sampleX
			alpha := layer.alpha[index] * opacity
			if alpha <= 0.01 {
				continue
			}
			tone := layer.tone[index]

			var paint color.RGBA
			if pal == nil {
				paint = mixColor(colorInk, colorInkSoft, tone)
				if tone > 0.74 {
					paint = colorPaperHi
					alpha *= 0.5
				}
			} else {
				paint = mixColor(pal.primary, pal.secondary, tone*0.65)
				if tone > 0.78 {
					paint = mixColor(pal.highlight, colorPaperHi, (tone-0.78)/0.22)
					alpha *= 0.42
				} else {
					alpha *= 0.82
				}
				if isSilhouetteEdge(layer, sampleX, sampleY) {
					paint = colorInk
					alpha = 0.9
				}
			}
			blendPixel(img, x, y, paint, alpha)
		}
	}
}

func drawRods(img *image.RGBA, centerX, centerY, progress float64, pal *shadowPalette) {
	sway := math.Sin(progress*math.Pi*3) * 8
	rod := colorInk
	if pal != nil {
		rod = color.RGBA{R: 0xA9, G: 0x7C, B: 0x42, A: 0xFF}
	}
	drawLine(img, centerX-96, centerY-190, centerX-26, centerY-34, 4, rod, 0.72)
	drawLine(img, centerX+96, centerY-166, centerX+34, centerY+8, 4, rod, 0.72)
	drawLine(img, centerX-118+sway, centerY-196, centerX-96, centerY-190, 3, rod, 0.5)
}

func drawDust(img *image.RGBA, frame int) {
	for index := 0; index < 14; index++ {
		x := float64((index*53+frame*7)%playWidth)
		y := float64(screenTop + (index*37+frame*3)%(screenFloor-screenTop))
		alpha := 0.10 + 0.16*math.Abs(math.Sin(float64(index)+float64(frame)*0.3))
		blendPixel(img, int(x), int(y), colorGlow, alpha)
	}
}

func renderPlayFrame(layer puppetLayer, preset string, frame, total int, pal *shadowPalette) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, playWidth, playHeight))
	progress := float64(frame) / float64(total)
	open := clampUnit(progress / 0.24)
	open = 1 - math.Pow(1-open, 2)
	flicker := 0.86 + 0.14*math.Abs(math.Sin(float64(frame)*0.7))

	drawStage(img, open, flicker, pal)
	centerX, centerY, rotation, scale := puppetTransform(preset, progress)
	drawRods(img, centerX, centerY, progress, pal)
	drawPuppet(img, layer, centerX, centerY, rotation, scale, 0.97, pal)
	drawDust(img, frame)
	return img
}

func encodePlayGIF(frames []*image.RGBA) ([]byte, error) {
	palette := color.Palette{
		color.RGBA{R: 0x1E, G: 0x14, B: 0x11, A: 0xFF},
		color.RGBA{R: 0x2E, G: 0x1B, B: 0x15, A: 0xFF},
		color.RGBA{R: 0x4A, G: 0x2A, B: 0x1F, A: 0xFF},
		color.RGBA{R: 0x7E, G: 0x24, B: 0x19, A: 0xFF},
		color.RGBA{R: 0xA9, G: 0x36, B: 0x28, A: 0xFF},
		color.RGBA{R: 0xC3, G: 0x9A, B: 0x4B, A: 0xFF},
		color.RGBA{R: 0xE4, G: 0xC8, B: 0x8E, A: 0xFF},
		color.RGBA{R: 0xF6, G: 0xE3, B: 0xBA, A: 0xFF},
		color.RGBA{R: 0xFF, G: 0xF3, B: 0xD4, A: 0xFF},
		color.RGBA{R: 0x2A, G: 0x17, B: 0x12, A: 0xFF},
		color.RGBA{R: 0x6B, G: 0x33, B: 0x22, A: 0xFF},
		color.RGBA{R: 0xFF, G: 0xD9, B: 0x8A, A: 0xFF},
	}
	anim := &gif.GIF{LoopCount: 0}
	for _, frame := range frames {
		paletted := image.NewPaletted(frame.Bounds(), palette)
		draw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})
		anim.Image = append(anim.Image, paletted)
		anim.Delay = append(anim.Delay, playDelay)
	}
	buffer := &bytes.Buffer{}
	if err := gif.EncodeAll(buffer, anim); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func encodePoster(frame *image.RGBA) ([]byte, error) {
	buffer := &bytes.Buffer{}
	if err := jpeg.Encode(buffer, frame, &jpeg.Options{Quality: 88}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func GenerateShadowPlay(src image.Image, preset string, colored bool) ([]byte, []byte, int, error) {
	layer := buildPuppetLayer(src)
	var pal *shadowPalette
	if colored {
		selected := paletteForLayer(layer)
		pal = &selected
	}
	frames := make([]*image.RGBA, 0, playFrames)
	for frame := 0; frame < playFrames; frame++ {
		frames = append(frames, renderPlayFrame(layer, preset, frame, playFrames, pal))
	}
	gifData, err := encodePlayGIF(frames)
	if err != nil {
		return nil, nil, 0, err
	}
	poster, err := encodePoster(frames[int(float64(len(frames))*0.72)])
	if err != nil {
		return nil, nil, 0, err
	}
	return gifData, poster, len(frames), nil
}
