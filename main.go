package main

import (
	"image"
	"log"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var buffers = make(map[string]*beep.Buffer)

func main() {
	if err := load(); err != nil {
		log.Fatal("loading error:", err)
	}
	sr := beep.SampleRate(41000)
	speaker.Init(sr, sr.N(time.Second/30))

	pixelgl.Run(run)
	//	play("cow.mp3")
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
	if s, ok := buffers[name]; ok {
		speaker.Play(s.Streamer(0, s.Len()))
	}
	return nil
}

func run() {
	cowpic, err := loadPicture("assets/image/cow.jpg")
	if err != nil {

		log.Fatal(err)
	}
	cowsprite := pixel.NewSprite(cowpic, cowpic.Bounds())

	cfg := pixelgl.WindowConfig{
		Title:     "aerth soundboard",
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

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Blue)
		cowsprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center()))
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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func getbutton(pixel.Vec) string {
	return "cow.mp3"
}
