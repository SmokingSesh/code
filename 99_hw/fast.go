package main

import (
	"bufio"
	//"encoding/json"
	"fmt"
	"hw3/model"

	// "hw3/model"
	"io"
	//"io/ioutil"
	// "log"
	"os"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	seenBrowsers := make(map[string]bool)
	var foundUsers strings.Builder

	lineNum := 0
	for scanner.Scan() {
		var user model.User
		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			panic(err)
		}
		isAndroid := false
		isMSIE := false
		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				seenBrowsers[browser] = true
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				seenBrowsers[browser] = true
			}
		}
		if isAndroid && isMSIE {
			email := strings.ReplaceAll(user.Email, "@", " [at] ")
			foundUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", lineNum, user.Name, email))

		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
