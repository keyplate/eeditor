package main

import (
	"fmt"
	"os"
)

func (g *Game)getFile() error {
    if len(os.Args) < 2 {
        return fmt.Errorf("Too few arguments")
    }

    filePath := os.Args[1]
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    
    g.file = file
    return nil
}

func (g *Game)loadFile() error {
    b, err := os.ReadFile(g.file.Name())
    if err != nil {
        return err
    }
    
    g.gapBuffer.Insert(string(b))
    return  nil
}

func (g *Game)saveFile() error {
    err := os.WriteFile(g.file.Name(), []byte(g.gapBuffer.String()), 0666)
    
    if err != nil {
        return err
    }

    return nil
}
