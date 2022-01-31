package plasticator

import (
	"errors"
	"image"
	"image/draw"
	"math"
)

type srcImage struct {
	img    *image.NRGBA
	width  float64
	height float64
}

type IPlastic interface {
	Width() int
	Height() int
	Lens(centerX, centerY, radius, effectIntensity int) error
	Swirl(centerX, centerY, radius, step int) error
	Image() *image.NRGBA
}

func NewImage(img *image.NRGBA) IPlastic {
	b := img.Bounds()
	tmpImg := (*image.NRGBA)(image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy())))
	draw.Draw(tmpImg, tmpImg.Bounds(), img, b.Min, draw.Src)

	return &srcImage{
		img:    tmpImg,
		width:  float64(img.Bounds().Max.X),
		height: float64(img.Bounds().Max.Y),
	}
}

func (m *srcImage) Width() int {
	return int(m.width)
}

func (m *srcImage) Height() int {
	return int(m.height)
}

func (m *srcImage) Lens(centerX, centerY, radius, effectIntensity int) error {
	if effectIntensity > 100 || effectIntensity < 1 {
		return errors.New("effectIntensity must be range 1 .. 100")
	}

	effectIntensityF := float64(effectIntensity) / 10
	centerXF := float64(centerX)
	centerYF := float64(centerY)
	offsetX := float64(0)
	offsetY := float64(0)

	radiusF := float64(radius)
	radiusSquared := float64(radius * radius)

	srcPixels := m.img.Pix
	dstPixels := make([]uint8, len(srcPixels))
	copy(dstPixels, srcPixels)

	for y := -radiusF; y < radiusF; y++ {
		for x := -radiusF; x < radiusF; x++ {
			if x*x+y*y <= radiusSquared {
				offsetX = x + centerXF
				offsetY = y + centerYF

				if offsetX < 0 || offsetX >= m.width || offsetY < 0 || offsetY >= m.height {
					continue
				}

				destPosition := int(offsetY*m.width + offsetX)
				destPosition *= 4
				r := math.Sqrt(x*x + y*y)
				alpha := math.Atan2(y, x)
				interpolationFactor := r / radiusF
				r = interpolationFactor*r + (1.0-interpolationFactor)*effectIntensityF*math.Sqrt(r)
				newY := r * math.Sin(alpha)
				newX := r * math.Cos(alpha)
				offsetX = newX + centerXF
				offsetY = newY + centerYF

				if offsetX < 0 || offsetX >= m.width || offsetY < 0 || offsetY >= m.height {
					continue
				}

				x0 := math.Floor(newX)
				xf := x0 + 1
				y0 := math.Floor(newY)
				yf := y0 + 1
				deltaX := int(newX - x0)
				deltaY := int(newY - y0)

				pos0 := int(((y0+centerYF)*m.width + x0 + centerXF) * 4)
				pos1 := int(((y0+centerYF)*m.width + xf + centerXF) * 4)
				pos2 := int(((yf+centerYF)*m.width + x0 + centerXF) * 4)
				pos3 := int(((yf+centerYF)*m.width + xf + centerXF) * 4)

				for k := 0; k < 4; k++ {
					componentX0 := int(srcPixels[pos1+k]-srcPixels[pos0+k])*deltaX + int(srcPixels[pos0+k])
					componentX1 := int(srcPixels[pos3+k]-srcPixels[pos2+k])*deltaX + int(srcPixels[pos2+k])
					finalPixelComponent := (componentX1-componentX0)*deltaY + componentX0

					if finalPixelComponent > 255 {
						dstPixels[destPosition+k] = 255
					} else {
						if finalPixelComponent < 0 {
							dstPixels[destPosition+k] = 0
						} else {
							dstPixels[destPosition+k] = uint8(finalPixelComponent)
						}

					}
				}
			}
		}
	}

	m.img.Pix = dstPixels

	return nil
}

func (m *srcImage) Swirl(centerX, centerY, radius, step int) error {
	radiusF := float64(radius)
	radiusSquared := float64(radius * radius)
	centerXF := float64(centerX)
	centerYF := float64(centerY)
	stepF := float64(step) / -100

	srcPixels := m.img.Pix
	dstPixels := make([]uint8, len(srcPixels))
	copy(dstPixels, srcPixels)

	for y := -radiusF; y < radiusF; y++ {
		for x := -radiusF; x < radiusF; x++ {
			if x*x+y*y >= radiusSquared {
				continue
			}

			r := math.Sqrt(x*x + y*y)
			alpha := math.Atan2(y, x)

			destPosition := (y+centerYF)*m.width + x + centerXF
			destPosition *= 4

			degrees := (alpha * 180.0) / math.Pi
			degrees += r * 10 * stepF

			alpha = (degrees * math.Pi) / 180.0
			newY := math.Floor(r * math.Sin(alpha))
			newX := math.Floor(r * math.Cos(alpha))
			sourcePosition := (newY+centerYF)*m.width + newX + centerXF
			sourcePosition *= 4

			dstPixels[int(destPosition+0)] = srcPixels[int(sourcePosition+0)]
			dstPixels[int(destPosition+1)] = srcPixels[int(sourcePosition+1)]
			dstPixels[int(destPosition+2)] = srcPixels[int(sourcePosition+2)]
			dstPixels[int(destPosition+3)] = srcPixels[int(sourcePosition+3)]
		}
	}

	m.img.Pix = dstPixels

	return nil
}

func (m *srcImage) Image() *image.NRGBA {
	return m.img
}
