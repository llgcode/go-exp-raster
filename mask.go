package raster

import (
	"image"
	"image/color"
	"unsafe"
)

// Renders the mask to the canvas
func DrawSolidRGBA(dest *image.RGBA, mask *image.Alpha, rgba color.RGBA) {
	rect := dest.Bounds().Intersect(mask.Bounds())
	minX := uint32(rect.Min.X)
	minY := uint32(rect.Min.Y)
	maxX := uint32(rect.Max.X)
	maxY := uint32(rect.Max.Y)

	pixColor := *(*uint32)(unsafe.Pointer(&rgba))

	cs1 := pixColor & 0xff00ff
	cs2 := (pixColor >> 8) & 0xff00ff

	stride1 := uint32(dest.Stride)
	stride2 := uint32(mask.Stride)

	maxY *= stride1
	var pix, pixm []uint8
	var pixelx uint32
	var x, y1, y2 uint32
	for y1, y2 = minY*stride1, minY*stride2; y1 < maxY; y1, y2 = y1+stride1, y2+stride2 {
		pix = dest.Pix[y1:]
		pixm = mask.Pix[y2:]
		pixelx = minX * 4
		for x = minX; x < maxX; x++ {
			alpha := uint32(pixm[x])
			p := (*uint32)(unsafe.Pointer(&pix[pixelx]))
			if alpha != 0 {
				invAlpha := 0xff - alpha
				ct1 := (*p & 0xff00ff) * invAlpha
				ct2 := ((*p >> 8) & 0xff00ff) * invAlpha

				ct1 = ((ct1 + cs1*alpha) >> 8) & 0xff00ff
				ct2 = (ct2 + cs2*alpha) & 0xff00ff00
				*p = ct1 + ct2
			}
			pixelx += 4
		}
	}
}
