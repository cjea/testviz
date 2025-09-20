package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

func main() {
	in := os.Stdin
	if len(os.Args) > 1 && os.Args[1] != "-" {
		f, err := os.Open(os.Args[1])
		if err != nil {
			die(err)
		}
		defer f.Close()
		in = f
	}

	doc, err := html.Parse(in)
	if err != nil {
		die(err)
	}

	injectSortButtons(doc)
	transform(doc, false)
	injectFilterInput(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		die(err)
	}
	_, _ = io.Copy(os.Stdout, &buf)
}

func transform(n *html.Node, inFiles bool) {
	jsBytes, err := os.ReadFile("script.js")
	if err != nil {
		panic("could not read script.js")
	}
	scriptJS := string(jsBytes)

	if n.Type == html.ElementNode {
		if hasAttrKV(n, "id", "topbar") {
			n.Attr = []html.Attribute{{Key: "id", Val: "files-container"}}
		}

		switch n.Data {
		case "script":
			for c := n.FirstChild; c != nil; {
				next := c.NextSibling
				n.RemoveChild(c)
				c = next
			}
			n.Attr = filterOut(n.Attr, "src")
			n.AppendChild(&html.Node{Type: html.TextNode, Data: scriptJS})

		case "select":
			if hasAttrKV(n, "id", "files") {
				n.Data = "div"
				n.Attr = []html.Attribute{{Key: "id", Val: "files"}}
				inFiles = true
			}

		case "option":
			if inFiles {
				val := getAttr(n, "value")
				n.Data = "div"
				n.Attr = []html.Attribute{
					{Key: "class", Val: "row"},
					{Key: "value", Val: val},
					{Key: "style", Val: "cursor: pointer; margin-bottom: 5px; color: rgba(220,220,220,1);"},
					{Key: "onmouseenter", Val: "this.style.background='green';"},
					{Key: "onmouseleave", Val: "this.style.background='';"},
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childInFiles := inFiles
		if n.Type == html.ElementNode && n.Data == "div" && hasAttrKV(n, "id", "files") {
			childInFiles = true
		}
		transform(c, childInFiles)
	}
}

func hasAttrKV(n *html.Node, k, v string) bool {
	for _, a := range n.Attr {
		if a.Key == k && a.Val == v {
			return true
		}
	}
	return false
}
func getAttr(n *html.Node, k string) string {
	for _, a := range n.Attr {
		if a.Key == k {
			return a.Val
		}
	}
	return ""
}
func filterOut(attrs []html.Attribute, key string) []html.Attribute {
	out := attrs[:0]
	for _, a := range attrs {
		if a.Key != key {
			out = append(out, a)
		}
	}
	return out
}
func die(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

func injectSortButtons(doc *html.Node) {
	nav := findByID(doc, "nav")
	if nav == nil {
		return
	}
	mk := func(id, label string) *html.Node {
		n := &html.Node{
			Type: html.ElementNode,
			Data: "button",
			Attr: []html.Attribute{
				{Key: "id", Val: id},
				{Key: "type", Val: "button"},
				{Key: "style", Val: "margin-right:10px; margin-bottom: 10px; padding:2px 8px; cursor:pointer;"},
			},
		}
		n.AppendChild(&html.Node{Type: html.TextNode, Data: label})
		return n
	}

	btnDesc := mk("sortDesc", "Sort by coverage (DESC)")
	btnAsc := mk("sortAsc", "Sort by coverage (ASC)")
	btnName := mk("sortName", "Sort by name")

	if nav.FirstChild != nil {
		nav.InsertBefore(btnName, nav.FirstChild)
		nav.InsertBefore(btnAsc, nav.FirstChild)
		nav.InsertBefore(btnDesc, nav.FirstChild)
	} else {
		nav.AppendChild(btnDesc)
		nav.AppendChild(btnAsc)
		nav.AppendChild(btnName)
	}
}

func findByID(n *html.Node, want string) *html.Node {
	var res *html.Node
	var walk func(*html.Node)
	walk = func(x *html.Node) {
		if res != nil {
			return
		}
		if x.Type == html.ElementNode {
			for _, a := range x.Attr {
				if a.Key == "id" && a.Val == want {
					res = x
					return
				}
			}
		}
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return res
}

func injectFilterInput(doc *html.Node) {
	nav := findByID(doc, "nav")
	if nav == nil {
		return
	}

	input := &html.Node{
		Type: html.ElementNode,
		Data: "input",
		Attr: []html.Attribute{
			{Key: "id", Val: "filter"},
			{Key: "type", Val: "text"},
			{Key: "placeholder", Val: "filter (regex)…"},
			{Key: "style", Val: "margin-right:10px; padding:2px 6px; width:220px;"},
		},
	}

	// place AFTER sort buttons, BEFORE the files list
	filesDiv := findByID(doc, "files")
	if filesDiv != nil && filesDiv.Parent == nav {
		nav.InsertBefore(input, filesDiv)
	} else {
		nav.AppendChild(input)
	}
}
