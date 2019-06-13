package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"

	//"strconv"
	"strings"
)

func main() {

	var file string

	flag.StringVar(&file, "file", "", "srt file")

	flag.Parse()

	if file == "" {
		fmt.Println("no file.")
		return
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	str := string(bytes)

	blockList := strings.Split(str, "\n\n")

	for k, block := range blockList {

		if k+1 == len(blockList) {
			break
		}

		lineList := strings.Split(block, "\n")

		if len(lineList) < 3 {
			fmt.Println("ERROR:", block)
			continue
		}

		lineListNext := strings.Split(blockList[k+1], "\n")

		if len(lineListNext) < 3 {
			fmt.Println("ERROR:", blockList[k+1])
			continue
		}

		timeLineList := strings.Split(lineList[1], " --> ")
		timeLineListNext := strings.Split(lineListNext[1], " --> ")

		if timeLineList[1] == timeLineListNext[0] {

			tmp := strings.Split(timeLineListNext[0], ",")

			if len(tmp) != 2 {
				fmt.Println("ERROR:", blockList[k+1])
				continue
			}

			tmpI, err := strconv.Atoi(tmp[1])
			if err != nil {
				// handle error
				fmt.Println(err)
				continue
			}

			if tmpI < 900 {

				tmpI = tmpI + 100

				tmp[1] = fmt.Sprintf("%d", tmpI)

				timeLineListNext[0] = strings.Join(tmp, ",")

			} else {

				tmp = strings.Split(tmp[0], ":")

				tmpI, err := strconv.Atoi(tmp[2])
				if err != nil {
					// handle error
					fmt.Println(err)
					continue
				}

				tmpI++
				tmp[2] = fmt.Sprintf("%d", tmpI)

				tmpStr := strings.Join(tmp, ":")

				timeLineListNext[0] = fmt.Sprintf("%s,000", tmpStr)
			}

			lineListNextNew := strings.Join(timeLineListNext, " --> ")

			str = strings.Replace(str, lineListNext[1], lineListNextNew, 1)
		}
	}

	blockList = strings.Split(str, "\n\n")

	d := make(map[int][]string)

	var i int
	i = 1

	for k, block := range blockList {

		d[i] = append(d[i], block)

		if (k+1)%50 == 0 {
			i++
		}
	}

	for fileNumber, list := range d {
		srtData := strings.Join(list, "\n\n")

		fileName := fmt.Sprintf("%s.new.%d.srt", file, fileNumber)

		err = ioutil.WriteFile(fileName, []byte(srtData), 0644)

		if err != nil {
			// handle error
			fmt.Println(err)
			continue
		} else {
			fmt.Println(fileName)
		}
	}
}
