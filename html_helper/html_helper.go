package html_helper

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func GetHtmlTitleAndDescription(html []byte) (title, description string, err error) {
	htmlReader := bytes.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return title, description, errors.New(fmt.Sprintf("Failde to parse html: %v", err))
	}
	metaTag := doc.Find(`meta[name$=escription]`)
	description, _ = metaTag.Attr("content")
	if description == "" {
		metaTag = doc.Find(`meta[property$=escription]`)
		description, _ = metaTag.Attr("content")
	}
	spaceRegExp := regexp.MustCompile(`\s+`)

	description = spaceRegExp.ReplaceAllString(description, " ")
	description = strings.Replace(description, " ", " ", -1)

	titleTag := doc.Find("title")
	title = titleTag.Text()
	title = spaceRegExp.ReplaceAllString(title, " ")
	title = strings.Replace(title, " ", " ", -1)
	return title, description, nil
}
