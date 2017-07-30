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

var animals = make(map[string]animal)

type animal struct {
	buffer *beep.Buffer
	rect   pixel.Rect
	sprite *pixel.Sprite
}

func main() {
	if err := load(); err != nil {
		log.Fatal("loading error:", err)
	}
	sr := beep.SampleRate(41000)
	speaker.Init(sr, sr.N(time.Second/30))
	pixelgl.Run(run)

}

func load() error {
	// create animal buttons
	for i, name := range []string{"cow", "horse", "chicken", "sheep", "cat", "dog"} {
		log.Println("loading", name)
		pic, err := loadPicture("assets/image/" + name + ".jpg")
		if err != nil {
			return err
		}
		var a animal
		a.sprite = pixel.NewSprite(pic, pic.Bounds())

		x := 100.0
		y := -200 + pic.Bounds().Max.Y*(float64(i)+1)
		if i > 2 {
			x = 420
			y = -200 + pic.Bounds().Max.Y*(float64(i)+1) - (pic.Bounds().Max.Y * 3)
		}

		a.rect = pic.Bounds().Moved(pixel.V(float64(x), y))
		f, err := os.Open("assets/sound/" + name + ".mp3")
		if err != nil {
			return err
		}
		s, format, err := mp3.Decode(f)
		if err != nil {
			return err
		}
		buf := beep.NewBuffer(format)
		buf.Append(s)
		a.buffer = buf
		animals[name] = a
	}

	if len(animals) == 0 {
		return fmt.Errorf("no animals found")
	}

	return nil
}

func play(name string) error {
	if animal, ok := animals[name]; ok {
		log.Printf("The %s goes ____", name)
		speaker.Play(animal.buffer.Streamer(0, animal.buffer.Len()))
		return nil
	}
	return fmt.Errorf("%q not found", name)
}

func run() {

	// create window
	cfg := pixelgl.WindowConfig{
		Title:     "aerth animals",
		Bounds:    pixel.R(0, 0, 800, 600),
		Resizable: false,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	frames := 0
	second := time.Tick(time.Second)
	debug := false // toggle debug with '='
	fs := false    // toggle fullscreen with 'f'
	for !win.Closed() && !win.JustPressed(pixelgl.KeyEscape) {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Black)
		for _, animal := range animals {
			animal.sprite.Draw(win, pixel.IM.Moved(animal.rect.Center()))
		}

		if win.JustPressed(pixelgl.KeyEqual) {
			debug = !debug
		}
		if debug {
			highlightbuttons(win)
		}
		if win.JustPressed(pixelgl.KeyF) {
			fs = !fs
			if fs {
				win.SetMonitor(pixelgl.PrimaryMonitor())
			} else {
				win.SetMonitor(nil)
			}

		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			button := getbutton(win.MousePosition())
			if button != "" {
				if err := play(button); err != nil {
					log.Println(err)
				}
			}
		}
		win.Update()
		frames++
		select {
		default:
		case <-second:
			win.SetTitle(fmt.Sprintf("aerth animals [ fps: %v | delta: %v ]", frames, dt))
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
	for name, a := range animals {
		if a.rect.Contains(v) {
			return name
		}
	}
	return ""
}
func highlightbuttons(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	for _, animal := range animals {
		r := animal.rect
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
