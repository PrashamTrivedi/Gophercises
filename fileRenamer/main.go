package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type fileR struct {
	name string
	path string
}

func main() {

	// fileName := "birthday_001.txt"
	// newName, err := match(fileName, 5)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
	// fmt.Println(newName)

	dir := "sample"

	toRename := make([]fileR, 0)
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		name := info.Name()
		if _, err := match(name); err == nil {
			toRename = append(toRename, fileR{path: path, name: name})
		}

		return nil
	})
	for _, f := range toRename {
		fmt.Printf("%q\n", f)
		var n fileR
		var err error
		n.name, err = match(f.name)
		if err != nil {
			fmt.Println("Error In matching:", f.path, err.Error())
		}

		n.path = filepath.Join(dir, n.name)
		err = os.Rename(f.path, n.path)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("moving %s -> %s\n", f.path, n.path)
	}

	// files, err := ioutil.ReadDir(dir)
	// if err != nil {
	// 	panic(err)
	// }
	// count := 0
	// var toRename []string
	// for _, file := range files {
	// 	if !file.IsDir() {
	// 		_, err := match(file.Name(), 4)
	// 		if err == nil {
	// 			count++
	// 			toRename = append(toRename, file.Name())
	// 		}
	// 	}
	// }
	// for _, orig := range toRename {
	// 	origPath := filepath.Join(dir, orig)
	// 	newFileName, err := match(orig, count)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	newPath := filepath.Join(dir, newFileName)
	// 	err = os.Rename(origPath, newPath)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	fmt.Printf("moving %s -> %s\n", origPath, newPath)
	// }
}

func match(fileName string) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}

	return fmt.Sprintf("%s - %d.%s", strings.Title(name), number, ext), nil
}
