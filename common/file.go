package common

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func saveToFile(filePath string, info string) {
	file, e := os.OpenFile("./data/"+filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if e != nil {
		Error.Println(fmt.Sprintf("Save To File :%s Error！", filePath))
		return
	}
	defer file.Close()

	_, e = file.WriteString(info)
	if e != nil {
		Error.Println(fmt.Sprintf("Save To File :%s Error！", filePath))
		return
	}
}

func getFileBuff(filePath string) {
	file, err := os.OpenFile(filePath+"file.log", os.O_RDONLY, 0600)
	if err != nil {
		Error.Println(err)
		return
	}
	defer file.Close()

	buff := bufio.NewReader(file)
	for i := 1; ; i++ {
		line, err := buff.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		fmt.Printf("%d line: %s", i, string(line))
		// 文件已经到达结尾
		if err == io.EOF {
			break
		}
	}
}
