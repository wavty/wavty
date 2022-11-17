package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateText = `---
title: {{printf "%q" .Title}}
---

{{print "{{% pageinfo %}}"}}
[{{.Title}}]({{.URL}})
{{print "{{% /pageinfo %}}"}}
`

	postStr = `
## 方法一：

时间复杂度 $O()$，空间复杂度 $O()$。

{{< tabpane persistLang=false >}}
{{< tab header="Solutions:" disabled=true />}}

{{< tab header="solution.go" lang="go" >}}

{{< /tab >}}

{{< /tabpane >}}
`
)

type Meta struct {
	Title string
	URL   string
}

func main() {
	var no = flag.Int("n", 1, "problem number.")
	var title = flag.String("t", "例子", "题目中文名.")
	var url = flag.String("u", "https://leetcode.cn/problems/two-sum/", "the problem internet link.")
	flag.Parse()

	meta := Meta{
		Title: fmt.Sprintf("%04d.%s", *no, *title),
		URL:   *url,
	}
	fileName := getFileFullName(*no, *url)
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("algorithm").Parse(templateText))
	// Run the template to verify the output.
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var render bytes.Buffer
	err = t.Execute(&render, meta)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	_, err = io.WriteString(f, render.String()+postStr)
	if err != nil {
		log.Fatal(err)
	}
}

func getFileFullName(no int, url string) string {
	dirs := []string{"0000-0499", "0500-0999", "1000-1499", "1500-1999", "2000-2499", "2500-2999"}
	err := os.MkdirAll(dirs[no/500], os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	patterns := strings.Split(strings.TrimSuffix(url, "/"), "/")
	filename := fmt.Sprintf("%04d.%s.md", no, patterns[len(patterns)-1])
	return filepath.Join("content/docs/Algorithm",dirs[no/500], filename)
}
