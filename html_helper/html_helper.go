package html_helper

import "strings"

const (
	closeAngelBracketChar = byte('<')
	backSlashChar         = byte('\\')
	closeBracket          = byte('"')
)

func GetHtmlTitle(htmlString string) string {
	var title string
	var isCloseChar bool
	substr := "<title>"
	i := strings.Index(htmlString, substr)
	if i == -1 {
		return ""
	}

	i += 7
	startIndex := i

	for !isCloseChar {
		i++
		if htmlString[i] == closeAngelBracketChar {
			isCloseChar = true
		}
	}

	title = htmlString[startIndex:i]
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Replace(title, " ", " ", -1)
	title = strings.TrimSpace(title)
	title = strings.Replace(title, "\t", " ", -1)


	return title
}

func GetHtmlDescription(htmlStr string) string {
	var description string
	var isCloseChar bool
	var lastChar byte
	substrLower := "name=\"description\" content=\""
	substrUpper := "name=\"Description\" content=\""

	i := strings.Index(htmlStr, substrLower)
	if i == -1 {
		i = strings.Index(htmlStr, substrUpper)
		if i == -1 {
			return ""
		}
	}

	i += len(substrLower)
	startIndex := i

	lastChar = htmlStr[i]

	for !isCloseChar {
		i++
		curChar := htmlStr[i]
		if curChar == closeBracket && lastChar != backSlashChar {
			isCloseChar = true
		} else {
			lastChar = curChar
		}
	}

	description = htmlStr[startIndex:i]
	description = strings.Replace(description, "\n", "", -1)
	description = strings.TrimSpace(description)
	description = strings.Replace(description, "\t", " ", -1)
	description = strings.Replace(description, " ", " ", -1)
	return description
}
