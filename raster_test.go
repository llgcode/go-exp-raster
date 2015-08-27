package raster

import (
	"image"
	"image/color"
	"testing"

	"code.google.com/p/freetype-go/freetype/raster"
	"github.com/llgcode/draw2d/draw2dbase"
	"github.com/llgcode/draw2d/draw2dimg"
)

type Path struct {
	points []float64
}

func (p *Path) LineTo(x, y float64) {
	if len(p.points)+2 > cap(p.points) {
		points := make([]float64, len(p.points)+2, len(p.points)+32)
		copy(points, p.points)
		p.points = points
	} else {
		p.points = p.points[0 : len(p.points)+2]
	}
	p.points[len(p.points)-2] = x
	p.points[len(p.points)-1] = y
}

func TestFreetypeRasterizer(t *testing.T) {
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	rasterizer := raster.NewRasterizer(200, 200)
	rasterizer.UseNonZeroWinding = false
	rasterizer.Start(raster.Point{
		X: raster.Fix32(10 * 256),
		Y: raster.Fix32(190 * 256)})
	for j := 0; j < len(poly); j = j + 2 {
		rasterizer.Add1(raster.Point{
			X: raster.Fix32(poly[j] * 256),
			Y: raster.Fix32(poly[j+1] * 256)})
	}
	painter := raster.NewRGBAPainter(img)
	painter.SetColor(color)
	rasterizer.Rasterize(painter)

	draw2dimg.SaveToPngFile("output/raster/TestFreetype.png", img)
}

func TestFreetypeRasterizerNonZeroWinding(t *testing.T) {
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	rasterizer := raster.NewRasterizer(200, 200)
	rasterizer.UseNonZeroWinding = true
	rasterizer.Start(raster.Point{
		X: raster.Fix32(10 * 256),
		Y: raster.Fix32(190 * 256)})
	for j := 0; j < len(poly); j = j + 2 {
		rasterizer.Add1(raster.Point{
			X: raster.Fix32(poly[j] * 256),
			Y: raster.Fix32(poly[j+1] * 256)})
	}
	painter := raster.NewRGBAPainter(img)
	painter.SetColor(color)
	rasterizer.Rasterize(painter)

	draw2dimg.SaveToPngFile("output/TestFreetypeRasterizerNonZeroWinding.png", img)
}

func TestRasterizer(t *testing.T) {
	bounds := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(bounds)
	mask := image.NewAlpha(bounds)
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)

	poly := Polygon(p.points)
	rgba := color.RGBA{0, 0, 0, 0xff}
	r := NewRasterizer()
	r.Fill(mask, poly, false)
	DrawSolidRGBA(img, mask, rgba)
	draw2dimg.SaveToPngFile("output/TestRasterizer.png", img)
}

func TestRasterizerNonZeroWinding(t *testing.T) {
	bounds := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(bounds)
	mask := image.NewAlpha(bounds)
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)
	poly := Polygon(p.points)
	rgba := color.RGBA{0, 0, 0, 0xff}
	r := NewRasterizer()
	r.Fill(mask, poly, true)
	DrawSolidRGBA(img, mask, rgba)
	draw2dimg.SaveToPngFile("output/TestRasterizerNonZeroWinding.png", img)
}

func BenchmarkFreetype(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)

	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		rasterizer := raster.NewRasterizer(200, 200)
		rasterizer.UseNonZeroWinding = false
		rasterizer.Start(raster.Point{
			X: raster.Fix32(10 * 256),
			Y: raster.Fix32(190 * 256)})
		for j := 0; j < len(poly); j = j + 2 {
			rasterizer.Add1(raster.Point{
				X: raster.Fix32(poly[j] * 256),
				Y: raster.Fix32(poly[j+1] * 256)})
		}
		painter := raster.NewRGBAPainter(img)
		painter.SetColor(color)
		rasterizer.Rasterize(painter)
	}
}

func BenchmarkFreetypeNonZeroWinding(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)

	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		rasterizer := raster.NewRasterizer(200, 200)
		rasterizer.UseNonZeroWinding = true
		rasterizer.Start(raster.Point{
			X: raster.Fix32(10 * 256),
			Y: raster.Fix32(190 * 256)})
		for j := 0; j < len(poly); j = j + 2 {
			rasterizer.Add1(raster.Point{
				X: raster.Fix32(poly[j] * 256),
				Y: raster.Fix32(poly[j+1] * 256)})
		}
		painter := raster.NewRGBAPainter(img)
		painter.SetColor(color)
		rasterizer.Rasterize(painter)
	}
}

func BenchmarkRasterizer(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)
	poly := Polygon(p.points)
	rgba := color.RGBA{0, 0, 0, 0xff}
	rasterizer := NewRasterizer()
	for i := 0; i < b.N; i++ {
		bounds := image.Rect(0, 0, 200, 200)
		img := image.NewRGBA(bounds)
		mask := image.NewAlpha(bounds)
		rasterizer.Fill(mask, poly, false)
		DrawSolidRGBA(img, mask, rgba)
	}
}

func BenchmarkRasterizerNonZeroWinding(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	draw2dbase.TraceCubic(&p, []float64{10, 190, 10, 10, 190, 10, 190, 190}, 0.5)
	poly := Polygon(p.points)
	rgba := color.RGBA{0, 0, 0, 0xff}
	rasterizer := NewRasterizer()
	for i := 0; i < b.N; i++ {
		bounds := image.Rect(0, 0, 200, 200)
		img := image.NewRGBA(bounds)
		mask := image.NewAlpha(bounds)
		rasterizer.Fill(mask, poly, true)
		DrawSolidRGBA(img, mask, rgba)
	}
}
