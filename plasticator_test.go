package plasticator_test

import (
	"gopkg.in/check.v1"
	"image"
	"image/color"
	"image/draw"
	"plasticator"
	"testing"
)

type suite struct {
	img    *image.NRGBA
	tmpImg *image.NRGBA
}

var _ = check.Suite(&suite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

func HLine(x1, x2, y int, img *image.NRGBA, col color.RGBA) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

func VLine(y1, y2, x int, img *image.NRGBA, col color.RGBA) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

func (s *suite) SetUpSuite(c *check.C) {
	width := 201
	height := 131
	s.img = image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(s.img, s.img.Bounds(), &image.Uniform{C: color.RGBA{B: 255, A: 255}}, image.Point{}, draw.Src)
	col := color.RGBA{B: 255, G: 255, A: 255}

	for y := 0; y <= height; y += 10 {
		HLine(0, width, y, s.img, col)
	}

	for x := 0; x <= width; x += 10 {
		VLine(0, height, x, s.img, col)
	}
}

func (s *suite) SetUpTest(c *check.C) {
	b := s.img.Bounds()
	s.tmpImg = (*image.NRGBA)(image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy())))
	draw.Draw(s.tmpImg, s.tmpImg.Bounds(), s.img, b.Min, draw.Src)
}

func (s *suite) TestLens(c *check.C) {
	var err error
	var img plasticator.IPlastic

	img = plasticator.NewImage(s.tmpImg)
	err = img.Lens(img.Width()/2, img.Height()/2, img.Height()/2, 0)
	c.Assert(err, check.NotNil)

	img = plasticator.NewImage(s.tmpImg)
	err = img.Lens(img.Width()/2, img.Height()/2, img.Height()/2, 200)
	c.Assert(err, check.NotNil)

	for i := 1; i <= 100; i++ {
		img = plasticator.NewImage(s.tmpImg)
		err = img.Lens(img.Width()/2, img.Height()/2, img.Height()/2, i)
		c.Assert(err, check.IsNil)

		////TODO Compare with testdata
		//buf := new(bytes.Buffer)
		//err = png.Encode(buf, img.Image())
		//c.Assert(err, check.IsNil)
		//f, _ := os.Create(fmt.Sprintf("lens_%03d.png", i))
		//_, _ = io.Copy(f, bytes.NewReader(buf.Bytes()))
		//_ = f.Close()
	}
}

func (s *suite) TestSwirl(c *check.C) {
	var err error
	var img plasticator.IPlastic

	for i := -10; i <= 10; i++ {
		img = plasticator.NewImage(s.tmpImg)
		err = img.Swirl(img.Width()/2, img.Height()/2, img.Height()/2, i)
		c.Assert(err, check.IsNil)

		////TODO Compare with testdata
		//buf := new(bytes.Buffer)
		//err = png.Encode(buf, img.Image())
		//c.Assert(err, check.IsNil)
		//f, _ := os.Create(fmt.Sprintf("swirl_%03d.png", i))
		//_, _ = io.Copy(f, bytes.NewReader(buf.Bytes()))
		//_ = f.Close()
	}
}
