package main

import (
	"fmt"
	"github.com/bugra/kmeans"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func main() {

	// todo replace
	rdr, err := os.Open("test.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer rdr.Close()

	// decode img
	img, _, err := image.Decode(rdr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}

	// reduce size to make computation easier
	img = resize.Resize(100, 0, img, resize.Bilinear)

	// iterate the pixels
	bounds := img.Bounds()
	pxData := make([][]float64, bounds.Max.X*bounds.Max.Y)
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			c := make([]float64, 3)
			pxl := img.At(i, j)
			r, g, b, _ := pxl.RGBA()
			c[0], c[1], c[2] = float64(r), float64(g), float64(b)
			idx := i*bounds.Max.Y + j
			pxData[idx] = c
		}
	}

	count := 100
	// todo
	labels, err := kmeans.Kmeans(pxData, count, kmeans.EuclideanDistance, 100)
	if err != nil {
		log.Fatal(err)
	}

	// map pixels to clusters
	clusters := make([][][]float64, count)

	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			idx := i*bounds.Max.Y + j
			id := labels[idx]
			clusters[id] = append(clusters[id], pxData[idx])
		}
	}

	// calculate the cluster centers
	centers := make([]color.RGBA, count)

	for id, obs := range clusters {
		var r, g, b uint

		for i := 0; i < len(obs); i++ {
			r += uint(obs[i][0])
			g += uint(obs[i][1])
			b += uint(obs[i][2])
		}

		l := uint(len(obs))
		centers[id] = color.RGBA{uint8(r / l), uint8(g / l), uint8(b / l), 255}
	}

	fmt.Println(centers)
}
