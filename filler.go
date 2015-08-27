package raster

import (
	"fmt"
	"image"
	"math"
)

type Intersection struct {
	x       int
	cover   uint8
	step    uint8
	winding int8
	next    *Intersection
}

type Rasterizer struct {
	table []*Intersection
}

func NewRasterizer() *Rasterizer {
	r := new(Rasterizer)
	return r
}

func (r *Rasterizer) Fill(img *image.Alpha, p Polygon, nonZeroWindingRule bool) {
	r.table = make([]*Intersection, img.Bounds().Dy())
	var xmin, ymin, xmax, ymax float64
	xmin = p[0]
	ymin = p[1]
	xmax = xmin
	ymax = ymin
	var x, y float64
	for i := 2; i < len(p); i += 2 {
		x, y = p[i], p[i+1]
		if x > xmax {
			xmax = x
		} else if x < xmin {
			xmin = x
		}
		if y > ymax {
			ymax = y
		} else if y < ymin {
			ymin = y
		}
	}
	prevX, prevY := p[0], p[1]
	for i := 2; i < len(p); i += 2 {
		x, y = p[i], p[i+1]
		r.edge(prevX, prevY, x, y)
		prevX, prevY = x, y
	}
	r.edge(x, y, p[0], p[1])

	if nonZeroWindingRule {
		r.scanNonZero(img, int(ymin), int(ymax+0.5))
	} else {
		r.scanEvenOdd(img, int(ymin), int(ymax+0.5))
	}
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func (r *Rasterizer) edge(x1, y1, x2, y2 float64) {
	var swap, dy float64
	var winding int8 = 1
	if y2 < y1 {
		swap = x1
		x1 = x2
		x2 = swap
		swap = y1
		y1 = y2
		y2 = swap
		winding = -1
	}
	dy = y2 - y1

	x := x1
	y := y1

	var cover uint8
	if dy == 0 {
		cover = uint8(math.Abs(x-float64(int(x))) * 256)
		r.insert(int(x), int(y), cover, 0, winding)
		return
	}

	dx := (x2 - x1) / dy
	yend := float64(int(y2-1) + 1)
	step := uint8(math.Abs(dx))
	for y < yend {
		cover = uint8(math.Abs(x-float64(int(x))) * 256)
		r.insert(int(x), int(y), cover, step, winding)
		x += dx
		y += 1
	}

}

func (r *Rasterizer) insert(x int, y int, cover, step uint8, winding int8) {
	i := &Intersection{x, cover, step, winding, nil}
	current := r.table[y]
	var prev *Intersection
	for current != nil {
		if x < current.x {
			i.next = current
			break
		}
		prev = current
		current = current.next
	}
	if prev != nil {
		prev.next = i
	} else {
		r.table[y] = i
	}

}

func printIntersection(i *Intersection) {
	if i == nil {
		fmt.Print("nil")
	} else {
		for i != nil {
			fmt.Print(i.x, " ")
			i = i.next
		}
	}
	fmt.Println()
}

func (r *Rasterizer) scanEvenOdd(img *image.Alpha, ymin, ymax int) {
	var ix1, ix2 int
	var i, j *Intersection
	fill := true
	pix := img.Pix[ymin*img.Stride:]
	for y := ymin; ; {
		i = r.table[y]
		if i != nil {
			fill = true
			j = i.next
			for j != nil {
				if fill {
					ix1 = i.x
					ix2 = j.x
					if ix1 != ix2 {
						alpha := 0xff - i.cover
						pix[ix1] = alpha
						ix1++
						for ix1 < ix2 {
							pix[ix1] = 0xff
							ix1++
						}
						pix[ix2] = j.cover
					} else {
						pix[ix1] = 0xff >> 3
					}

				}
				fill = !fill
				i = j
				j = i.next
			}
		}
		y++
		if y < ymax {
			pix = pix[img.Stride:]
		} else {
			break
		}
	}
}

func (r *Rasterizer) scanNonZero(img *image.Alpha, ymin, ymax int) {
	var ix1, ix2 int
	var i, j *Intersection
	pix := img.Pix[ymin*img.Stride:]
	var winding int8 = 0
	for y := ymin; ; {
		i = r.table[y]
		if i != nil {
			winding = i.winding
			j = i.next
			for j != nil {
				if winding != 0 {
					ix1 = i.x
					ix2 = j.x
					if ix1 != ix2 {
						pix[ix1] = 0xff - i.cover
						ix1++
						for ix1 < ix2 {
							pix[ix1] = 0xff
							ix1++
						}
						pix[ix2] = j.cover
					} else {
						pix[ix1] = 0xff >> 3
					}
				}
				winding = winding + j.winding
				i = j
				j = i.next
			}
		}
		y++
		if y < ymax {
			pix = pix[img.Stride:]
		} else {
			break
		}
	}
}
