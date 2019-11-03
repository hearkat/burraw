package Burraw

import (
	"os"
	"path/filepath"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	burraw := newBurraw(dir)

	burraw.run()
}
