package controllers

import (
	"bytes"
	"math"
	"math/rand"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

// this is in its own file given that the workspace controller file is already quite long

func (self *FilesController) resetPreview(s style.TextStyle, cmd func() string) string {
	if self.c.Git() == nil {
		return ""
	}
	return s.Sprint(cmd())
}

func (self *FilesController) animateExplosion() {
	self.Explode(self.c.Views().Files, func() {
		self.c.PostRefreshUpdate(self.c.Contexts().Files)
	})
}

// Animates an explosion within the view by drawing a bunch of flamey characters
func (self *FilesController) Explode(v *gocui.View, onDone func()) {
	width := v.InnerWidth()
	height := v.InnerHeight()
	styles := []style.TextStyle{
		style.FgLightWhite.SetBold(),
		style.FgYellow.SetBold(),
		style.FgRed.SetBold(),
		style.FgBlue.SetBold(),
		style.FgBlack.SetBold(),
	}

	self.c.OnWorker(func(_ gocui.Task) error {
		max := 25
		for i := range max {
			image := getExplodeImage(width, height, i, max)
			style := styles[(i*len(styles)/max)%len(styles)]
			coloredImage := style.Sprint(image)
			self.c.OnUIThread(func() error {
				v.SetOrigin(0, 0)
				v.SetContent(coloredImage)
				return nil
			})
			time.Sleep(time.Millisecond * 20)
		}
		self.c.OnUIThread(func() error {
			v.Clear()
			onDone()
			return nil
		})
		return nil
	})
}

// Render an explosion in the given bounds.
func getExplodeImage(width int, height int, frame int, max int) string {
	// Predefine the explosion symbols
	explosionChars := []rune{'*', '.', '@', '#', '&', '+', '%'}

	// Initialize a buffer to build our string
	var buf bytes.Buffer

	// Initialize RNG seed
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// calculate the center of explosion
	centerX, centerY := width/2, height/2

	// calculate the max radius (hypotenuse of the view)
	maxRadius := math.Hypot(float64(centerX), float64(centerY))

	// calculate frame as a proportion of max, apply square root to create the non-linear effect
	progress := math.Sqrt(float64(frame) / float64(max))

	// calculate radius of explosion according to frame and max
	radius := progress * maxRadius * 2

	// introduce a new radius for the inner boundary of the explosion (the shockwave effect)
	var innerRadius float64
	if progress > 0.5 {
		innerRadius = (progress - 0.5) * 2 * maxRadius
	}

	for y := range height {
		for x := range width {
			// calculate distance from center, scale x by 2 to compensate for character aspect ratio
			distance := math.Hypot(float64(x-centerX), float64(y-centerY)*2)

			// if distance is less than radius and greater than innerRadius, draw explosion char
			if distance <= radius && distance >= innerRadius {
				// Make placement random and less likely as explosion progresses
				if random.Float64() > progress {
					// Pick a random explosion char
					char := explosionChars[random.Intn(len(explosionChars))]
					buf.WriteRune(char)
				} else {
					buf.WriteRune(' ')
				}
			} else {
				// If not explosion, then it's empty space
				buf.WriteRune(' ')
			}
		}
		// End of line
		if y < height-1 {
			buf.WriteRune('\n')
		}
	}

	return buf.String()
}
