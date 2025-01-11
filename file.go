package main

import (
	"fmt"
	"os"
)

func getFile() (*os.File, error) {
    if len(os.Args) < 2 {
        return nil, fmt.Errorf("Too few arguments")
    }

    filePath := os.Args[1]
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }

    return file, nil
}

func loadFile(file *os.File) (string, error) {
    b, err := os.ReadFile(file.Name())
    if err != nil {
        return "", err
    }

    return string(b), nil
}

