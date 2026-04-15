package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return walk(out, path, "", printFiles)
}

func walk(out io.Writer, path, prefix string, printFiles bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	filtered := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || printFiles {
			filtered = append(filtered, e)
		}
	}

	for i, e := range filtered {
		last := i == len(filtered)-1
		branch := "├───"
		nextPrefix := prefix + "│   "
		if last {
			branch = "└───"
			nextPrefix = prefix + "    "
		}

		info, err := e.Info()
		if err != nil {
			return err
		}

		if e.IsDir() {
			fmt.Fprintf(out, "%s%s%s\n", prefix, branch, e.Name())
			if err := walk(out, filepath.Join(path, e.Name()), nextPrefix, printFiles); err != nil {
				return err
			}
			continue
		}

		size := info.Size()
		if size == 0 {
			fmt.Fprintf(out, "%s%s%s (empty)\n", prefix, branch, e.Name())
		} else {
			fmt.Fprintf(out, "%s%s%s (%db)\n", prefix, branch, e.Name(), size)
		}
	}

	return nil
}
