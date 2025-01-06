package primitive

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

type ImageFormat string

const (
	PNG  ImageFormat = "png"
	JPG  ImageFormat = "jpg"
	JPEG ImageFormat = "jpeg"
	SVG  ImageFormat = "svg"
)

type Image struct {
	path     string
	created  time.Time
	lifetime time.Duration
}

func (img *Image) ExtendLife(by time.Duration) {
	img.lifetime = img.lifetime + by
}

func (img *Image) Path() string {
	return img.path
}

type ImageService struct {
	mu          sync.Mutex
	root        string
	images      map[string]*Image
	transformed map[string]*TransformedImage
	imageFormat ImageFormat
}

func (s *ImageService) GetImageByPath(path string) *Image {
	return s.images[path]
}

func (s *ImageService) GetImageByName(name string) *Image {
	path := fmt.Sprintf("%s/%s", s.root, name)
	return s.images[path]
}

func (s *ImageService) GetTransformedImageByName(name string) *TransformedImage {
	path := fmt.Sprintf("%s/%s", s.root, name)
	return s.transformed[path]
}

func (s *ImageService) GetTransformedImageByPath(path string) *TransformedImage {
	return s.transformed[path]
}

func (s *ImageService) NewImage(name string, ext ImageFormat) (*Image, error) {
	path := fmt.Sprintf("%s/%s.%s", s.root, name, ext)

	if _, ok := s.images[path]; ok {
		return nil, fmt.Errorf("image %s.%s already exists", name, ext)
	}

	switch ext {
	case PNG, JPG, JPEG, SVG:
		img := &Image{
			path:     path,
			created:  time.Now(),
			lifetime: time.Minute * 3,
		}
		s.images[path] = img

		return img, nil
	}

	return nil, fmt.Errorf("%s is an invalid image type", ext)
}

func (img *ImageService) Has(path string) bool {
	_, exists := img.images[path]
	return exists
}

func (s *ImageService) ToTransformedImage(base *Image) *TransformedImage {
	tImg := &TransformedImage{base: base, image: base}
	s.transformed[base.path] = tImg
	return tImg

}

func (s *ImageService) DeleteImage(img *Image) error {
	cmd := exec.Command("rm", img.path)
	if err := cmd.Run(); err != nil {
		return err
	}

	delete(s.images, img.path)
	delete(s.transformed, img.path)

	return nil
}

func (s *ImageService) TransformImage(base *TransformedImage, transform Transformer) (*TransformedImage, error) {
	img, err := s.NewImage(generateImageName(), s.imageFormat)
	if err != nil {
		return nil, err
	}

	dst := &TransformedImage{base: base.base, image: img, mods: base.mods}

	if err := transform.Transform(base, dst); err != nil {
		return nil, err
	}

	s.transformed[dst.image.path] = dst
	return dst, nil
}

func NewImageService(rootFolder string, useGc bool) (*ImageService, error) {
	fi, err := os.Stat(rootFolder)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", rootFolder)
	}

	if err = createEmptyDir(rootFolder); err != nil {
		return nil, err
	}

	service := &ImageService{
		root:        rootFolder,
		images:      make(map[string]*Image),
		transformed: make(map[string]*TransformedImage),
	}

	if useGc {
		go func() {
			for range time.Tick(time.Minute * 5) {
				service.mu.Lock()
				for _, img := range service.images {
					if !img.created.Add(img.lifetime).After(time.Now()) {
						service.DeleteImage(img)
					}
				}
				service.mu.Unlock()
			}
		}()
	}
	return service, nil
}

func ParseFormat(filename string) (ImageFormat, error) {

	switch path.Ext(filename) {
	case ".png":
		return PNG, nil
	case ".jpeg":
		return JPEG, nil
	case ".jpg":
		return JPG, nil
	case ".svg":
		return SVG, nil
	}

	return "", fmt.Errorf("invalid format: %s", path.Ext(filename))
}

func generateImageName() string {
	bytes := make([]byte, 16)
	for i := range bytes {
		lower := 'a' + rand.Intn('z'-'a')
		upper := 'A' + rand.Intn('Z'-'A')
		bit := rand.Intn(2)

		if bit == 0 {
			bytes[i] = byte(lower)
		} else {
			bytes[i] = byte(upper)
		}
	}

	return string(bytes)
}

func createEmptyDir(rootFolder string) error {
	// Check if the directory exists
	fi, err := os.Stat(rootFolder)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist, create it
			return os.Mkdir(rootFolder, 0755)
		}
		// Other errors (e.g., permission issues)
		return fmt.Errorf("error accessing %q: %w", rootFolder, err)
	}

	// Check if the path is a directory
	if !fi.IsDir() {
		return fmt.Errorf("%q exists but is not a directory", rootFolder)
	}

	// Clean the directory by removing all its contents
	args := fmt.Sprintf("-rf %s", rootFolder)
	cmd := exec.Command("rm", strings.Fields(args)...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clear directory %q: %w", rootFolder, err)
	}

	// Recreate the directory
	if err := os.Mkdir(rootFolder, 0755); err != nil {
		return fmt.Errorf("failed to recreate directory %q: %w", rootFolder, err)
	}

	return nil
}
