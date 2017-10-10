# codeviz

![codeviz](https://github.com/dancing-koala/codeviz/blob/master/codeviz_gen.png)

## What is it ?

Codeviz is a command-line application that takes a file containing code as input and generates an image by applying a color scheme.
It is heavily inspired by [minimap](https://github.com/Ivoah/minimap) and uses the great [chroma syntax highlighter](https://github.com/alecthomas/chroma) for Go.

## Why ?

Firstly, I loved the idea of [minimap](https://github.com/Ivoah/minimap): creating visual art from a codebase.
Secondly, I did it because it seemed fun and in order to improve my programming in Golang.

## How do I use it ?

### Build and install it

Execute the following command line in the projects folder:

`go install`

This will generate a binary file named codeviz in your $GOPATH/bin.

Make sure you have added $GOPATH/bin to your PATH for easier command-line usage.

### Usage

#### Quick usage

`codeviz <file>`

This command will read the file provided and generate an image using default settings.

#### Flags

`codeviz [flags] <file>`

`-o=<output_file>` specifies the name of the output file. Default is './codeviz_gen.png'.


`-w=12` specifies the width in pixels of a character. Default is 4.


`-h=12` specifies the height in pixels of a character. Default is 4.


`-maxW=768` specifies the maximum width in pixels of the generated image. Default is not limited.


`-maxH=1024` specifies the maximum height in pixels of the generated image. Default is not limited.


`-s=<style>` specifies a style to use. If not found, it will use chroma's defined fallback.


`-l=<lang>` specifies the language of the lexer. If not found, it will use chroma's defined fallback.

>Important: the path of the file to use **must come after** all flags!

#### Example

`codeviz -w=12 -h=10 -maxW=800 -maxH=600 -s=arduino -l=go -o=./example.png main.go`

This command will use the file **main.go** as a **go** file, applying the style **arduino** to it.
Each character will be **12** pixels wide and **10** pixels high.
The generated image's width will not be greater than **800** and its height will not be greater then **600**.
The image will be saved in a file named **example.png** in the current directory.

### FAQ

* What languages are supported ? Refer to [chroma's detailed list](https://github.com/alecthomas/chroma#supported-languages)
