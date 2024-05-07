package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/1boombacks1/zipViewer/internal/app"
)

func main() {
	var (
		port           int
		pathToZip, ext string
		err            error
	)

	for i := 1; i < len(os.Args); i = i + 2 {
		switch os.Args[i] {
		case "-p", "--port":
			port, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Println("Пропишите только число без ':' у аргумента '--port'")
				return
			}
		case "-e", "--ext":
			ext = os.Args[i+1]
		case "--path":
			pathToZip = os.Args[i+1]
		default:
			pathToZip = os.Args[i]
		}
	}

	if pathToZip == "" {
		fmt.Println("Пропишите путь к архиву")
		return
	}

	a := app.New(pathToZip)
	if port != 0 {
		a.SetPort(port)
	}
	if ext != "" {
		a.SetExt(ext)
	}

	a.Start()
}
