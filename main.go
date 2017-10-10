package main

import (
	"flag"
	"fmt"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	outFile, lang, styleName string
	w, h, maxW, maxH         int
)

const (
	DEFAULT_OUT_FILE  = "./codeviz_gen.png"
	DEFAULT_STYLE     = "monokai"
	DEFAULT_CHAR_SIZE = 4
	DEFAULT_MAX_SIZE  = -1
)

func init() {
	flag.StringVar(&outFile, "o", DEFAULT_OUT_FILE, "Path of the output file.")
	flag.StringVar(&lang, "l", "", "Language to parse, determined automatically if not provided.")
	flag.StringVar(&styleName, "s", DEFAULT_STYLE, "Color scheme to use.")
	flag.IntVar(&w, "w", DEFAULT_CHAR_SIZE, "Width of a char in pixels.")
	flag.IntVar(&h, "h", DEFAULT_CHAR_SIZE, "Height of a char in pixels.")
	flag.IntVar(&maxW, "maxW", DEFAULT_MAX_SIZE, "Maximum width of the output image in pixels.")
	flag.IntVar(&maxH, "maxH", DEFAULT_MAX_SIZE, "Maximum height of the output image in pixels.")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("No input file provided ¯\\_(ツ)_/¯")
		return
	}

	inFile := flag.Args()[0]

	fmt.Println("Reading input file...")
	result, err := ioutil.ReadFile(inFile)

	if err != nil {
		fmt.Println("Error while reading input file: ", err)
		return
	}

	code := string(result)

	fmt.Println("Computing image boundaries...")
	rowCount, colCount := countRowsAndCols(strings.Split(code, "\n"))

	left, top := 0, 0
	right, bottom := (colCount+1)*w, rowCount*(h+1)

	if maxW != -1 && maxW < right {
		right = maxW
	}

	if maxH != -1 && maxH < bottom {
		bottom = maxH
	}

	fmt.Printf("Generated image dimensions will be %d x %d.\n", right, bottom)

	fmt.Println("Retrieving appropriate lexer...")
	lexer := getAppropriateLexer(lang, inFile)

	if lexer == nil {
		lexer = lexers.Fallback
		fmt.Printf("No lexer found for the language %s, we will use %s.\n", lang, lexer.Config().Name)
	}

	fmt.Println("Retrieving appropriate style...")
	style := styles.Get(styleName)

	if style == nil {
		style = styles.Fallback
		fmt.Printf("No style found with name %s, we will use %s\n", styleName, style.Name)
	}

	fmt.Println("Tokenising file content...")
	iterator, err := lexer.Tokenise(nil, code)

	if err != nil {
		fmt.Println("Error during tokenisation:", err)
	}

	img := image.NewRGBA(image.Rect(left, top, right, bottom))

	offsetV, offsetH := 0, 0

	bgColour := style.Get(chroma.Background).Background
	bgColor := getColor(&bgColour)

	tabRegex := regexp.MustCompile("\t")
	newLineEndRegex := regexp.MustCompile(".*\n$")
	spaceRegex := regexp.MustCompile(" ")

	fmt.Println("Generating output image...")
	for _, tk := range iterator.Tokens() {

		if tk.Value == "\n" {
			offsetV += h + 1
			offsetH = 0
			continue
		}

		spaceCount := len(spaceRegex.FindAllStringIndex(tk.Value, -1))
		tabCount := len(tabRegex.FindAllStringIndex(tk.Value, -1))
		charCount := len([]rune(tk.Value)) - tabCount

		if spaceCount == charCount {
			offsetH += w * spaceCount
			continue
		}

		if charCount < 1 {
			continue
		}

		se := style.Get(tk.Type)
		c := getColor(&se.Colour)

		offsetH += tabCount * w

		drawRect(img, offsetH, offsetV, charCount*w, h, c)

		if newLineEndRegex.Match([]byte(tk.Value)) {
			offsetV += h + 1
			offsetH = 0
			continue
		}

		offsetH += w * charCount
	}

	// Fill background
	for i := 0; i < right; i++ {
		for j := 0; j < bottom; j++ {
			if _, _, _, a := img.At(i, j).RGBA(); a == 0 {
				img.Set(i, j, bgColor)
			}
		}
	}

	fmt.Println("Saving image to output file...")
	f, _ := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()

	err = png.Encode(f, img)

	if err != nil {
		fmt.Println("Error while writing generated image: ", err)
	}

	fmt.Println("Done !!!")
}

func drawRect(rgba *image.RGBA, left, top, width, height int, c color.Color) {
	rect := rgba.Bounds()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if left+x > rect.Dx()-1 {
				break
			}

			rgba.Set(left+x, top+y, c)
		}

		if top+y > rect.Dy()-1 {
			break
		}
	}
}

func getColor(origin *chroma.Colour) color.Color {
	return color.RGBA{origin.Red(), origin.Green(), origin.Blue(), 255}
}

func countRowsAndCols(lines []string) (int, int) {
	maxCol := 0

	for _, line := range lines {
		colCount := len([]rune(line))

		if colCount > maxCol {
			maxCol = colCount
		}
	}

	return len(lines), maxCol
}

func getAppropriateLexer(lang, filename string) chroma.Lexer {
	if len(lang) > 0 {
		return lexers.Get(lang)
	}

	return lexers.Match(filename)
}
