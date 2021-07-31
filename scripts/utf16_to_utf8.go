package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dollarkillerx/Quant/utils"
)

func main() {
	path16 := "./data_utf16"
	path8 := "./data_utf8"
	dir, err := ioutil.ReadDir(path16)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range dir {
		if v.IsDir() {
			continue
		}

		pt := fmt.Sprintf("%s/%s", path16, v.Name())
		file, err := ioutil.ReadFile(pt)
		if err != nil {
			log.Println(err)
			continue
		}

		utf8, err := utils.UTF16ToUTF8(file)
		if err != nil {
			log.Println(err)
			continue
		}

		pt = fmt.Sprintf("%s/%s", path8, v.Name())
		err = ioutil.WriteFile(pt, []byte(strings.ReplaceAll(string(utf8), "\ufeff", "")), 00666)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
