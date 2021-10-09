package main

import (
	"fmt"
	"os"
	"prashamhtrivedi/link"
	"strings"
)

var html = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to <span>another page</span></a>
</body>
</html>
`

func main() {
	r := strings.NewReader(html)
	links, error := link.Parse(r)
	if error != nil {
		panic(error)
	}
	fmt.Printf("%+v\n", links)
	fmt.Println()
	readFile("ex2.html")
	fmt.Println()
	readFile("ex3.html")
	fmt.Println()
	readFile("ex4.html")

}

func readFile(fileName string) {

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	linksFromFile, error := link.Parse(file)
	if error != nil {
		panic(error)
	}
	fmt.Printf("%+v\n", linksFromFile)
}
