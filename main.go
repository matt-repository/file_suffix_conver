package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fsc
// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o fsc.exe
// CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o fsc
func main() {
	path := flag.String("path", "", "the path")
	srcSuffix := flag.String("src", "", "the source Suffix")
	destSuffix := flag.String("dest", "", "the destination Suffix")
	subContains := flag.Bool("sub", false, "the sub files is contains")
	autoNumName := flag.Bool("auto", false, "the auto num name")

	flag.Parse()
	if *path == "" {
		fmt.Println("- Please enter the path")
		return
	}
	if *srcSuffix == "" {
		fmt.Println("- Please enter the source suffix")
		return
	}
	// 源后缀 to 新后缀
	if !*autoNumName {
		if *destSuffix == "" {
			fmt.Println("- Please enter the dest suffix")
			return
		}
		translateSuffix(*path, *srcSuffix, *destSuffix, *subContains)
	} else {
		//先加个下划线，防止文件名字相同
		translateSuffix(*path, *srcSuffix, "_"+*srcSuffix, *subContains)
		index = 0
		translateAutoNumName(*path, *srcSuffix, *subContains)
	}

}

var index int32 = 0

// translateSuffix 转换后缀
func translateSuffix(path, srcSuffix, destSuffix string, subTrans bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("%s the path read dir failed err:%s\n", path, err)
		return
	}
	if len(files) == 0 {
		fmt.Printf("%s the path files count is zero\n", path)
		return
	}
	// 按文件修改时间倒序
	sort.Slice(files, func(i, j int) bool {
		f1, _ := files[i].Info()
		f2, _ := files[j].Info()
		return f1.ModTime().After(f2.ModTime())
	})
	for _, f := range files {
		fileName := f.Name()
		src := filepath.Join(path, fileName)
		if f.IsDir() && subTrans {
			translateSuffix(src, srcSuffix, destSuffix, subTrans)
		}
		if len(fileName) <= len(srcSuffix) {
			continue
		}
		indexStart := len(fileName) - len(srcSuffix)
		fileFront := fileName[:indexStart]  //文件名后 oldSuffix外的前面名字
		fileSuffix := fileName[indexStart:] //文件名后缀（长度与oldSuffix相同）

		if fileSuffix == srcSuffix {
			index++
			newFileName := fileFront + destSuffix
			dest := filepath.Join(path, newFileName)
			fmt.Printf("%d: start converting file:[%s]\n", index, src)
			err = os.Rename(src, dest)
			if err != nil {
				fmt.Printf("%d: start converting file:[%s] failed err:%s\n", index, src, err)
			}
			fmt.Printf("%d:end converted file:[%s] to [%s]\n", index, src, dest)
		}
	}
	fmt.Printf("%s:%d convert successful!\n", path, index)
}

// translateAutoNumName 自动编号
func translateAutoNumName(path, srcSuffix string, subTrans bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("%s the path read dir failed err:%s\n", path, err)
		return
	}
	if len(files) == 0 {
		fmt.Printf("%s the path files count is zero\n", path)
		return
	}
	sort.Slice(files, func(i, j int) bool {
		f1, _ := files[i].Info()
		f2, _ := files[j].Info()
		return f1.ModTime().After(f2.ModTime())
	})

	for _, f := range files {
		fileName := f.Name()
		src := filepath.Join(path, fileName)
		if f.IsDir() && subTrans {
			translateAutoNumName(src, srcSuffix, subTrans)
		}
		if len(fileName) <= len(srcSuffix) {
			continue
		}
		indexStart := len(fileName) - len(srcSuffix)
		fileSuffix := fileName[indexStart:]

		if fileSuffix == srcSuffix {
			index++
			newFileName := fmt.Sprintf("%d%s", index, fileSuffix)
			dest := filepath.Join(path, newFileName)
			fmt.Printf("%d: start converting file:[%s]\n", index, src)
			err = os.Rename(src, dest)
			if err != nil {
				fmt.Printf("%d: start converting file:[%s] failed err:%s\n", index, src, err)
			}
			fmt.Printf("%d:end converted file:[%s] to [%s]\n", index, src, dest)
		}
	}
	fmt.Printf("%s:%d convert successful!\n", path, index)
}
