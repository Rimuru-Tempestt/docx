package docx

import (
	"github.com/beevik/etree"
	"strings"
)

type ReplaceContent struct {
	OldString string
	NewString string
}

type DocReplaceContents struct {
	Contents []ReplaceContent
}

func (d *Docx) RepeatingReplace(docReplace []DocReplaceContents) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(d.content); err != nil {
		return err
	}
	document := doc.SelectElement("w:document")
	// 构建替换用字符串
	baseStr := creatBase(document.SelectElement("w:body").ChildElements())

	d.content = ""
	for _, contents := range docReplace {
		d.content += baseStr
		for _, content := range contents.Contents {
			d.Replace(content.OldString, content.NewString, -1)
		}
	}

	// 构建空内容
	document.RemoveChildAt(0)
	body := etree.NewElement("w:body")
	document.AddChild(body)

	tmp := etree.NewDocument()
	tmp.ReadFromString(d.content)
	addChildren(body, tmp.ChildElements())
	d.content, _ = doc.WriteToString()
	return nil
}

func addChildren(body *etree.Element, children []*etree.Element) {
	for _, child := range children {
		body.AddChild(child)
	}
}

func creatBase(element []*etree.Element) string {
	writeSettings := etree.WriteSettings{}
	var builder strings.Builder
	for _, ele := range element {
		ele.WriteTo(&builder, &writeSettings)
	}
	return builder.String()
}
