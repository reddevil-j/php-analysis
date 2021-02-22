package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"github.com/z7zmey/php-parser/php7"
	"strings"
	l "medium/analyze/logger"
	visit "medium/analyze/visitor"

)

var fileList []string

func method_def_analysis(path string, basepath string) {

	l.Log(l.Info, "anlayzing file: %s", path)
	fileContents, _ := ioutil.ReadFile(path)
	parser := php7.NewParser(bytes.NewBufferString(string(fileContents)), path)
	parser.Parse()

	rootNode := parser.GetRootNode()

	nsresolver := visit.NewNamespaceResolver()
	rootNode.Walk(nsresolver)

	defvisit := visit.DefWalker {
		Writer: os.Stdout,
		Indent: "",
		NsResolver: nsresolver,
	}
	_ = defvisit
	visit.File = path
	visit.RelativePath = strings.TrimPrefix(path, basepath)
	rootNode.Walk(defvisit)
}


func main() {
	l.Level = l.Debug
	project_path := os.Args[1]

	err := filepath.Walk(project_path, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".php" ||
			filepath.Ext(path) == ".install" ||
			filepath.Ext(path) == ".engine" ||
			filepath.Ext(path) == ".module" ||
			filepath.Ext(path) == ".theme" ||
			filepath.Ext(path) == ".inc" {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}



	for _, file := range fileList {
		method_def_analysis(file, project_path)
	}

	// write functions to a file
	f, err := os.Create("functions.txt")
	if err == nil {
		for funcs := range visit.Functions {
			fmt.Fprintf(f,"%s\n",funcs)
		}
	}
	f.Close()

	// write methods to a file
	f, err = os.Create("methods.txt")
	if err == nil {
		for meths := range visit.Methods {
			fmt.Fprintf(f,"%s\n",meths)
		}
	}
	f.Close()



}
