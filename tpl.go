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
	templateIndex = `---
title: {{printf "%q" .Title}}
linkTitle: {{printf "%q" .Title}}
description: >
	Solutions to LeetCode Problems {{.Title}}.
---`

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
	writeFile(fileName, templateText, postStr, meta)
}

func getFileFullName(no int, url string) string {
	dirs := []int{(no / 50) * 50, (no/50+1)*50 - 1}
	dirName := filepath.Join("content/docs/Algorithm", fmt.Sprintf("%04d-%04d", dirs[0], dirs[1]))
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		meta := Meta{Title: fmt.Sprintf("%04d-%04d", dirs[0], dirs[1])}
		filename := filepath.Join(dirName, "_index.md")
		writeFile(filename, templateIndex, "", meta)
	}
	patterns := strings.Split(strings.TrimSuffix(url, "/"), "/")
	filename := fmt.Sprintf("%04d.%s.md", no, patterns[len(patterns)-1])
	return filepath.Join(dirName, filename)
}

func writeFile(fileName, tpl, help string, meta Meta) {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("algorithm").Parse(tpl))
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

	_, err = io.WriteString(f, render.String()+help)
	if err != nil {
		log.Fatal(err)
	}
}
