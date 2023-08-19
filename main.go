package main

import (
	"flag"
	"fmt"
	"os"
)

//CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fsc_linux
//CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o fsc_windows.exe
//CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o fsc_mac
func main() {
	path := flag.String("path", "", "the path")
	oldSuffix := flag.String("old", "", "the old Suffix")
	newSuffix := flag.String("new", "", "the new Suffix")
	subContains := flag.Bool("sub", true, "the sub files is contains")

	flag.Parse()
	if *path == "" {
		fmt.Println("- Please enter the path")
		return
	}
	if *oldSuffix == "" {
		fmt.Println("- Please enter the old suffix")
		return
	}
	if *newSuffix == "" {
		fmt.Println("- Please enter the new suffix")
		return
	}
	translate(*path, *oldSuffix, *newSuffix, *subContains)
}

func translate(path, oldSuffix, newSuffix string, subTrans bool) {
	transNum := 0
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	if len(files) == 0 {
		fmt.Println(path, "the path files count is zero,")
	}
	for _, f := range files {
		fileName := f.Name()
		if f.IsDir() && subTrans {
			translate(path+"\\"+fileName, oldSuffix, newSuffix, subTrans)
		}
		if len(fileName) <= len(oldSuffix) {
			continue
		}
		indexStart := len(fileName) - len(oldSuffix)
		fileFront := fileName[:indexStart]  //文件名后 oldSuffix外的前面名字
		fileSuffix := fileName[indexStart:] //文件名后缀（长度与oldSuffix相同）

		if fileSuffix == oldSuffix {
			transNum++
			fmt.Println("start converting file:", path+"\\"+fileName)
			newFileName := fileFront + newSuffix
			os.Rename(path+"\\"+fileName, path+"\\"+newFileName)
			fmt.Println("successful conversion:", path+"\\"+newFileName)
		}
	}
	fmt.Println(path, ":", transNum, "convert successful!")
}
