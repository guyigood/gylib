package imgresize

import (
	"github.com/nfnt/resize"
	z "github.com/nutzam/zgo"
	"golang.org/x/image/bmp"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

func Img_resize(filename string) string {
	img := ""
	f_type := z.FileType(filename)
	switch f_type {
	// JPEG
	case "bmp":
		img, _ = bmp_resize(filename)
	case "jpeg":
		// ImageJPEG
		img, _ = jpg_resize(filename)
		// JPG
	case "jpg":
		// ImageJPEG
		img, _ = jpg_resize(filename)
		// PNG
	case "png":
		// ImagePNG
		img, _ = png_resize(filename)
	case "gif":
		img, _ = gif_resize(filename)
	}

	return img
}

func jpg_resize(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		return "", err
	}
	file.Close()
	m := resize.Resize(200, 0, img, resize.Lanczos3)
	//m := resize.Thumbnail(200, 0, img, resize.Lanczos3)
	out, err := os.Create(filename + ".thumb")
	if err != nil {
		return "", err
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	return filename + ".thumb", nil
}

func png_resize(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	img, err := png.Decode(file)
	if err != nil {
		return "", err
	}
	file.Close()
	m := resize.Resize(200, 0, img, resize.Lanczos3)
	//m := resize.Thumbnail(640, 0, img, resize.Lanczos3)
	out, err := os.Create(filename + ".thumb")
	if err != nil {
		return "", err
	}
	defer out.Close()
	png.Encode(out, m)
	return filename + ".thumb", nil
}

func gif_resize(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	img, err := gif.Decode(file)
	if err != nil {
		return "", err
	}
	file.Close()
	m := resize.Resize(200, 0, img, resize.Lanczos3)
	//m := resize.Thumbnail(640, 0, img, resize.Lanczos3)
	out, err := os.Create(filename + ".thumb")
	if err != nil {
		return "", err
	}
	defer out.Close()
	gif.Encode(out, m, nil)
	return filename + ".thumb", nil
}

func bmp_resize(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	img, err := bmp.Decode(file)
	if err != nil {
		return "", err
	}
	file.Close()
	m := resize.Resize(200, 0, img, resize.Lanczos3)
	//m := resize.Thumbnail(640, 0, img, resize.Lanczos3)
	out, err := os.Create(filename + ".thumb")
	if err != nil {
		return "", err
	}
	defer out.Close()
	bmp.Encode(out, m)
	return filename + ".thumb", nil
}
