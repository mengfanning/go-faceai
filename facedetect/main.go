package main
 
import (
	"fmt"
	"os"
	"image"
	"log"
	"image/color"
	"path/filepath"
	"gocv.io/x/gocv"
	"github.com/Kagami/go-face"
)

const dataDir = "./data"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("参数不全")
		return
	}

	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	defer rec.Close()

	// Test image with 10 faces.
	testImagePristin := filepath.Join(dataDir, "mfn.jpg")
	// Recognize faces on that image.
	faces, err := rec.RecognizeFile(testImagePristin)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	
	
	var samples []face.Descriptor
	var cats []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		cats = append(cats, int32(i))
	}
	labels := []string{
		"mengfaning", 
	}
	rec.SetSamples(samples, cats)

	// parse args
	deviceID := os.Args[1]
	xmlFile := os.Args[2]
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return
	}
	defer  webcam.Close()

	// create display window
	window := gocv.NewWindow("face recognition")
	defer window.Close()

	// prepare image matrixz
	img := gocv.NewMat()
	defer img.Close()

	// color for then rect when detected
	green := color.RGBA{0, 255, 0, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Println("分类器 文件加载失败")
		return
	}

	for {
		if ok := webcam.Read(&img);  !ok {
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)

		for _, r := range rects {
			faceimg := img.Region(r)
			gocv.IMWrite("1.jpg", faceimg)


			// rec.recognize([]byte, 0)
			// fmt.Printf("识别的头像", faceimg)
			var who string 

			testImageNayoung := filepath.Join("./", "1.jpg")
			nayoungFace, err := rec.RecognizeSingleFile(testImageNayoung)

			if err != nil {
				fmt.Printf("识别出错")
			}

			if nayoungFace == nil {
				who = "dont't know"
			} else {
				catID := rec.Classify(nayoungFace.Descriptor)

				if catID < 0 {
					who = "dont"
				} else  {
					who = labels[catID]
				}
				
			}
			

			gocv.Rectangle(&img, r, green, 3)
			size := gocv.GetTextSize("", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, who, pt, gocv.FontHersheyPlain, 2, green, 2)

		}
		
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}

	}
}
	