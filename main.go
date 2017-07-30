package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var buffers = make(map[string]*beep.Buffer)
var buttons = make(map[string]pixel.Rect)
var sprites = make(map[string]*pixel.Sprite)

func main() {
	if err := load(); err != nil {
		log.Fatal("loading error:", err)
	}
	sr := beep.SampleRate(41000)
	speaker.Init(sr, sr.N(time.Second/30))

	pixelgl.Run(run)
}

func load() error {
	for _, name := range []string{"cow.mp3"} {
		f, err := os.Open("assets/sound/" + name)
		if err != nil {
			return err
		}
		s, format, err := mp3.Decode(f)
		if err != nil {
			return err
		}
		buf := beep.NewBuffer(format)
		buf.Append(s)
		buffers[name] = buf
	}
	return nil
}

func play(name string) error {
	if s, ok := buffers[name+".mp3"]; ok {
		speaker.Play(s.Streamer(0, s.Len()))
		return nil
	}
	return fmt.Errorf("%q not found", name)
}

func run() {
	for i, name := range []string{"cow"} {
		pic, err := loadPicture("assets/image/" + name + ".jpg")
		if err != nil {

			log.Fatal(err)
		}
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprites[name] = sprite
		buttons[name] = pic.Bounds().Moved(pixel.V(100*(float64(i)+1), 100*(float64(i)+1)))
	}
	cfg := pixelgl.WindowConfig{
		Title:     "aerth animals",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	last := time.Now()
	frames := 0
	second := time.Tick(time.Second)
	debug := false
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Blue)
		for k, sprite := range sprites {
			sprite.Draw(win, pixel.IM.Moved(buttons[k].Center()))
		}

		if win.JustPressed(pixelgl.KeyEqual) {
			debug = !debug
		}
		if debug {
			highlightbuttons(win)
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			button := getbutton(win.MousePosition())
			if button != "" {
				play(button)
			}
		}
		win.Update()
		frames++
		select {
		default:
		case <-second:
			log.Println(frames, dt)
			frames = 0
		}
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := Asset(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func getbutton(v pixel.Vec) string {
	for i, r := range buttons {
		if r.Contains(v) {
			return i
		}
	}
	return ""
}
func highlightbuttons(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	for _, r := range buttons {
		imd.Color = RandomColor()
		imd.Push(r.Min, r.Max)
		imd.Rectangle(2)
		imd.Draw(win)
	}

}

func RandomColor() pixel.RGBA {

	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	len := math.Sqrt(r*r + g*g + b*b)
	//if len == 0 {
	//	goto again
	//}
	return pixel.RGB(r/len, g/len, b/len)

}
