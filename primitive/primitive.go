package primitive

import (
	"fmt"
	"os/exec"
	"strings"
)

// (main.go)
// type cache struct {
// 	images     map[string][]*primitive.TransformedImage
// 	imgService *primitive.ImageService
// }

// func (c *cache) IsKeyExpired(key string) bool {
// 	for _, image := range c.images[key] {
// 		if c.imgService.GetTransformedImageByPath(image.GetImage().Path()) == nil {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (c *cache) UpdateCache(key string, images []*primitive.TransformedImage) {
// 	c.images[key] = make([]*primitive.TransformedImage, 0, len(images))
// 	c.images[key] = append(c.images[key], images...)
// }

// func (c *cache) IsEmpty(key string) bool {
// 	return len(c.images[key]) == 0
// }

// func (c *cache) GetCache(key string) []*primitive.TransformedImage {
// 	return c.images[key]
// }

// func (c *cache) ClearCache(key string) error {
// 	for _, image := range c.images[key] {
// 		if err := c.imgService.DeleteImage(image.GetImage()); err != nil {
// 			return err
// 		}
// 	}
// 	c.images[key] = make([]*primitive.TransformedImage, 0)
// 	return nil
// }

// type stepHandler func(step string, cache *cache, base *primitive.TransformedImage) (string, error)

// type TransformationOption struct {
// 	NextStep string
// 	ImageUrl string
// }

// type stepProcessor struct {
// 	handlers     map[string]stepHandler
// 	imageService *primitive.ImageService
// 	cache        *cache
// }

// func (sp *stepProcessor) ProcessStep(step string, baseImage *primitive.TransformedImage) (string, error) {
// 	handler, exists := sp.handlers[step]
// 	if !exists {
// 		return "", fmt.Errorf("invalid step: %s", step)
// 	}

// 	return handler(step, sp.cache, baseImage)
// }

// const resourcePath string = "/Users/juniorgreen/Documents/go_exercises/public"
// const parentImageName string = "parent"

// const htmlTemplate string = `
// <!DOCTYPE html>
// <html lang="en">

// <head>
//     <meta charset="UTF-8">
//     <meta name="viewport" content="width=device-width, initial-scale=1.0">
//     <title>Image Maker</title>
// </head>

// <body>
//     <h1>Pick One</h1>
//     {{range .TransformationOption}}   
//     <ul>

//         <li>
//             <a href="/transform?step={{.NextStep}}">
//                 <img src="{{.ImageUrl}}"/>
//             </a>
//         </li>
//     </ul>
//     {{end}}
// </body>
// </html>
// `

// func main() {
// 	imgService, err := primitive.NewImageService(resourcePath, true)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	cache := &cache{
// 		images:     make(map[string][]*primitive.TransformedImage),
// 		imgService: imgService,
// 	}
// 	processor := NewStepProcessor(imgService, cache)

// 	//Attach handlers
// 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(resourcePath))))
// 	http.Handle("/", getRootHandler())
// 	http.Handle("/transform", getStepHandler(processor))
// 	http.Handle("/upload_image", getImageUploadHandler(imgService))

// 	//Start server
// 	fmt.Println("Starting the server on :8080")
// 	http.ListenAndServe(":8080", nil)
// }

// func getImageUploadHandler(service *primitive.ImageService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		formFile, header, err := r.FormFile("file")
// 		if err != nil {
// 			http.Error(w, "no image file", http.StatusBadRequest)
// 			return
// 		}

// 		format, err := primitive.ParseFormat(header.Filename)
// 		if err != nil {
// 			http.Error(w, "invalid image format", http.StatusBadRequest)
// 			return
// 		}

// 		service.GetTransformedImageByName()
// 		img, err := service.NewImage(parentImageName, format)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		service.ToTransformedImage(img)

// 		file, err := os.Create(img.Path())
// 		if err != nil {
// 			http.Error(w, "error opening file", http.StatusInternalServerError)
// 			return
// 		}

// 		if _, err = io.Copy(file, formFile); err != nil {
// 			http.Error(w, "error saving image", http.StatusInternalServerError)
// 			return
// 		}

// 		img.ExtendLife(time.Minute * 60)

// 		url := fmt.Sprintf("transform?step=n&image=parent.%s", format)
// 		http.Redirect(w, r, url, http.StatusSeeOther)
// 	}
// }

// func getStepHandler(sp *stepProcessor) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		filename := r.URL.Query().Get("image")
// 		html, err := sp.ProcessStep(r.URL.Query().Get("step"), sp.imageService.GetTransformedImageByName(filename))
// 		if err != nil {
// 			http.Redirect(w, r, "/", http.StatusSeeOther)
// 			return
// 		}
// 		fmt.Fprint(w, html)
// 	}
// }

// func getRootHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		html := `
// 		<html>
// 		<form action="/upload_image" method="post" enctype="multipart/form-data">
// 			<div>
// 			<label for="file">Choose file to upload</label>
// 			<input type="file" id="file" name="file" multiple />
// 			</div>
// 			<div>
// 				<button type="submit">Submit</button>
// 			</div>
// 		</form>
// 		</html>`
// 		fmt.Fprint(w, html)
// 	}
// }

// func NewStepProcessor(service *primitive.ImageService, cache *cache) *stepProcessor {
// 	return &stepProcessor{
// 		handlers: map[string]stepHandler{
// 			"n":   handleStepN,
// 			"m":   handleStepM,
// 			"rep": handleStepRep,
// 			"a":   handleStepA,
// 			"bg":  handleStepBg,
// 		},
// 		imageService: service,
// 		cache:        cache,
// 	}
// }

// func handleStepN(step string, cache *cache, base *primitive.TransformedImage) (string, error) {
// 	//1. Look over cached images to check if any are expired
// 	if cache.IsKeyExpired(step) || cache.IsEmpty(step) {
// 		if err := cache.ClearCache(step); err != nil {
// 			return "", err
// 		}
// 		options := [5]uint{50, 100, 250, 500, 1000}
// 		images := make([]*primitive.TransformedImage, 0, len(options))

// 		for _, option := range options {
// 			transformer := primitive.ShapeCountTransformer{ShapeCount: option}
// 			tImg, err := cache.imgService.TransformImage(base, &transformer)
// 			if err != nil {
// 				return "", err
// 			}
// 			images = append(images, tImg)
// 		}
// 		cache.UpdateCache(step, images)
// 	}
// 	//1a. If expired images exist:
// 	//	I - delete all images currently in cache
// 	// 	II - reprocess all transfomrations and store in cache

// 	//1b. Otherwise, use cached images

// 	//2. inject all images in html template and return its string
// 	images := cache.GetCache(step)
// 	tImages := make([]TransformationOption, 0, len(images))

// 	for _, img := range images {
// 		tImages = append(tImages, TransformationOption{
// 			NextStep: "m",
// 			ImageUrl: img.GetImage().Path(),
// 		})
// 	}
// 	var buf bytes.Buffer

// 	tmpl := template.Must(template.New("").Parse(htmlTemplate))
// 	if err := tmpl.Execute(&buf, tImages); err != nil {
// 		return "", err
// 	}

// 	return buf.String(), nil
// }

// func handleStepM(step string, cache *cache, base *primitive.TransformedImage) (string, error) {
// 	panic("not implemented")
// }

// func handleStepRep(step string, cache *cache, base *primitive.TransformedImage) (string, error) {
// 	panic("not implemented")
// }

// func handleStepA(step string, cache *cache, base *primitive.TransformedImage) (string, error) {
// 	panic("not implemented")
// }

// func handleStepBg(step string, cache *cache, base *primitive.TransformedImage) (string, error) {
// 	panic("not implemented")
// }


type Mode uint8
type Color string

type ModifierOption func(*modifiers)

func WithMode(mode Mode) ModifierOption {
	return func(mods *modifiers) {
		mods.mode = mode
	}
}

func WithShapeCount(count uint) ModifierOption {
	return func(mods *modifiers) {
		mods.shapeCount = count
	}
}

func WithRepetitions(reps uint) ModifierOption {
	return func(mods *modifiers) {
		mods.repetitions = reps
	}
}

func WithAlpha(alpha uint8) ModifierOption {
	return func(mods *modifiers) {
		mods.alpha = alpha
	}
}

func WithBackgroundColor(hex Color) ModifierOption {
	return func(mods *modifiers) {
		mods.backgroundColor = hex
	}
}

const (
	Combo Mode = iota
	Triangle
	Rectangle
	Ellipse
	Circle
	RotatedRectangle
	Beziers
	RotatedEllipse
	Polygon
)

type Transformer interface {
	Transform(src *TransformedImage, dst *TransformedImage) error
}

type modifiers struct {
	mode            Mode
	shapeCount      uint
	repetitions     uint
	alpha           uint8
	backgroundColor Color
}

func (m modifiers) Set(options ...ModifierOption) modifiers {
	mods := modifiers{
		mode:            m.mode,
		shapeCount:      m.shapeCount,
		repetitions:     m.repetitions,
		alpha:           m.alpha,
		backgroundColor: m.backgroundColor,
	}

	for _, option := range options {
		option(&mods)
	}

	return mods
}

type TransformedImage struct {
	image *Image
	base  *Image
	mods  modifiers
}

func (i *TransformedImage) GetImage() *Image {
	return i.image
}

func (i *TransformedImage) GetBaseImage() *Image {
	return i.base
}

type ModeTransformer struct {
	Mode Mode
}

func (t *ModeTransformer) Transform(src *TransformedImage, dst *TransformedImage) error {
	dst.mods = src.mods.Set(WithMode(t.Mode))
	return execPrimitive(dst.mods, src, dst)
}

type ShapeCountTransformer struct {
	ShapeCount uint
}

func (t *ShapeCountTransformer) Transform(src *TransformedImage, dst *TransformedImage) error {
	dst.mods = src.mods.Set(WithShapeCount(t.ShapeCount))
	return execPrimitive(dst.mods, src, dst)
}

type RepititionsTransformer struct {
	Count uint
}

func (t *RepititionsTransformer) Transform(src *TransformedImage, dst *TransformedImage) error {
	dst.mods = src.mods.Set(WithRepetitions(t.Count))
	return execPrimitive(dst.mods, src, dst)
}

type AlphaTransformer struct {
	Alpha uint8
}

func (t *AlphaTransformer) Transform(src *TransformedImage, dst *TransformedImage) error {
	dst.mods = src.mods.Set(WithAlpha(t.Alpha))
	return execPrimitive(dst.mods, src, dst)
}

type BgColorTransformer struct {
	Hex string
}

func (t *BgColorTransformer) Transform(src *TransformedImage, dst *TransformedImage) error {
	dst.mods = src.mods.Set(WithBackgroundColor(Color(t.Hex)))
	return execPrimitive(dst.mods, src, dst)
}

func execPrimitive(mods modifiers, in, out *TransformedImage) error {
	args := fmt.Sprintf("-i %s -o %s -n %d -m %d -rep %d -r 256 -s 1024 -a %d -bg %s", in.base.path, out.image.path, mods.shapeCount, mods.mode, mods.repetitions, mods.alpha, mods.backgroundColor)
	cmd := exec.Command("primitive", strings.Fields(args)...)
	return cmd.Run()
}

func TransformImage(t Transformer, base, dst *TransformedImage) {
}

func NewModifier(m ...ModifierOption) modifiers {
	mods := modifiers{
		mode:            Triangle,
		shapeCount:      100,
		repetitions:     0,
		alpha:           128,
		backgroundColor: "avg",
	}

	for _, option := range m {
		option(&mods)
	}

	return mods
}
