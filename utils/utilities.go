package utils

import (
	"html/template"
	"strings"
)

// @todo make it handle select and other tags
func Html(n string, m map[string]string) template.HTML {
	var a string
	h := "<" + strings.TrimSpace(n)
	for k, v := range m {
		// create attributes
		switch {
		case (k == "checked" || k == "selected") && v == "true":
			a += " " + strings.TrimSpace(k) + " "
		case k != "text" && k != "checked" && k != "selected":
			a += " " + strings.TrimSpace(k) + "=\"" + template.HTMLEscapeString(strings.TrimSpace(v)) + "\""
		}
	}
	if t, ok := m["text"]; ok {
		h += a + ">" + template.HTMLEscapeString(t) + "</" + n + ">"
	} else {
		h += a + "/>"
	}

	return template.HTML(h)
}
