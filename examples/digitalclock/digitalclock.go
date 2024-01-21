// This program displays the current time on the screen.
package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/Solace7/go-smartwatch"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var watch *smartwatch.Watch

var showSeconds bool

var timeFormat string
var timePadding string

func updateTimeFormat(tf string, tp string) {
	timeFormat = tf
	timePadding = tp
}

func main() {
	showSeconds := true
	timeFormat = "00:00"
	timePadding = "%02d:%02d"
	watch, _ = smartwatch.Open()
	width, height := watch.Size()

	// Pick an appropriate font.
	font := &tinyfont.Org01 // fallback font
	fonts := []*tinyfont.Font{&freemono.Bold9pt7b, &freemono.Bold12pt7b, &freemono.Bold18pt7b, &freemono.Bold24pt7b}

	if showSeconds {
		updateTimeFormat("00:00:00", "%02d:%02d:%02d")
	}

	for _, f := range fonts {
		// If the font fits on this screen, use it.
		lineWidth, _ := tinyfont.LineWidth(f, timeFormat)
		if int16(lineWidth) <= width {
			font = f
		}
	}
	fontHeight := int16(font.Glyphs[0].Height)

	// Draw the current time.
	fmt.Printf("Time settings: %v, %s, %s", showSeconds, timeFormat, timePadding)
	for {
		// Clear the screen.
		watch.FillScreen(color.RGBA{0, 0, 0, 255})
		// Draw the current time (with second precision).
		now := time.Now()
		// Quick-and-dirty hack to get the current time (roughly) without
		// relying on locale support, which is not yet supported in TinyGo.
		hour := (now.Unix() / 60 / 60) % 24
		minute := now.Unix() / 60 % 60
		second := now.Unix() % 60
		msg := (fmt.Sprintf(timePadding, hour, minute))
		if showSeconds {
			msg = (fmt.Sprintf(timePadding, hour, minute, second))
		}
		textWidth, _ := tinyfont.LineWidth(font, msg)
		tinyfont.WriteLine(watch, font, width/2-int16(textWidth/2), height/2+fontHeight/2, msg, color.RGBA{255, 255, 255, 255})

		watch.Display()

		// Sleep until the next second/minute #TODO find a less expensive implementation.
		time.Sleep(time.Second - time.Duration(now.Nanosecond())*time.Nanosecond)
	}
}
