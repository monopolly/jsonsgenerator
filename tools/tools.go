package tools

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var Padding = 100

func FindTags(line string) (res map[string]string) {
	res = make(map[string]string)
	re := regexp.MustCompile(`.*?\{(.*?)\}`)
	rev := regexp.MustCompile(`\{(.*?)\}`)
	for _, x := range re.FindAllString(line, -1) {
		x = strings.TrimSpace(x)
		i := strings.Index(x, "{")
		if i == -1 {
			continue
		}
		k := x[:i]
		v := rev.FindString(x)
		v = strings.ReplaceAll(v, "{", "")
		v = strings.ReplaceAll(v, "}", "")
		res[k] = v
	}
	return
}

func ExistList(list []string, index string) bool {
	for _, x := range list {
		if x == index {
			return true
		}
	}
	return false
}

func Replace(v, tag, value string) string {
	return strings.ReplaceAll(v, tag, value)
}

func Title(s string) string {
	if s == "" {
		return ""
	}
	up := strings.ToUpper(string(s[0]))
	return up + string(s[1:])
}

// any text from start until \t\n\space
func Value(raw string, key string, endtag ...string) (value string) {
	i := strings.Index(raw, key)
	if i == -1 {
		return
	}
	for _, x := range raw[i+len(key):] {
		switch {
		case unicode.IsSpace(x):
			if endtag != nil {
				value = value + string(x)
			} else {
				return
			}
		default:
			if endtag != nil && string(x) == endtag[0] {
				return
			}
			value = value + string(x)
		}
	}

	return
}

// any {{tag value}}
func Between(raw string, startTag, endTag string) (value string) {
	indexStart := strings.Index(raw, startTag)
	if indexStart == -1 {
		return
	}

	from := indexStart + len(startTag)
	indexEnd := strings.Index(raw[from:], endTag)
	if indexEnd == -1 {
		return
	}
	to := from + indexEnd //+ len(endTag)

	return strings.TrimSpace(raw[from:to])
}

// any {{tag value}}
// func BetweenAndDelete(comment *string, startTag, endTag string) (value string) {
// 	raw := *comment
// 	indexStart := strings.Index(raw, startTag)
// 	if indexStart == -1 {
// 		return
// 	}

// 	from := indexStart + len(startTag)
// 	indexEnd := strings.Index(raw[from:], endTag)
// 	if indexEnd == -1 {
// 		return
// 	}
// 	to := from + indexEnd //+ len(endTag)
// 	value = strings.TrimSpace(raw[from:to])
// 	strings.ReplaceAll(*comment, value, "")
// 	return
// }

// any {{tag value}}
func Betweens(raw string, startTag, endTag string) (value []string) {
	un := map[string]bool{}
	for {
		indexStart := strings.Index(raw, startTag)
		if indexStart == -1 {
			return
		}

		from := indexStart + len(startTag)
		indexEnd := strings.Index(raw[from:], endTag)
		if indexEnd == -1 {
			return
		}
		to := from + indexEnd //+ len(endTag)

		v := raw[from:to]
		if un[v] {
			raw = raw[to:]
			continue
		}
		value = append(value, v)
		raw = raw[to:]
	}

}

func Formatline(k string, v any) (res string) {
	switch Padding / 10 {
	case 0:
		return fmt.Sprintf("%-15s %v", k, v)
	case 1:
		return fmt.Sprintf("%-15s %v", k, v)
	case 2:
		return fmt.Sprintf("%-25s %v", k, v)
	case 3:
		return fmt.Sprintf("%-35s %v", k, v)
	default:
		return fmt.Sprintf("%-45s %v", k, v)
	}
}

func FormatLineFix(pad int, k string, v any) (res string) {
	p := strings.Repeat(" ", pad)
	return fmt.Sprintf("%s%s%v", p, k, v)
}

func FormatLineSwift(k string, v any, lens int) (res string) {
	return fmt.Sprintf("%s%s%v", k, strings.Repeat(" ", lens-len(k)), v)
}

func SqlWord(s string) bool {
	return sqlkeywords[s]
}

// нельзя использовать в качестве sql полей
var sqlkeywords = map[string]bool{
	"all":               true,
	"analyse":           true,
	"analyze":           true,
	"and":               true,
	"any":               true,
	"array":             true,
	"as":                true,
	"asc":               true,
	"asymmetric":        true,
	"authorization":     true,
	"binary":            true,
	"both":              true,
	"case":              true,
	"cast":              true,
	"check":             true,
	"collate":           true,
	"collation":         true,
	"column":            true,
	"concurrently":      true,
	"constraint":        true,
	"create":            true,
	"cross":             true,
	"current_catalog":   true,
	"current_date":      true,
	"current_role":      true,
	"current_schema":    true,
	"current_time":      true,
	"current_timestamp": true,
	"current_user":      true,
	"day":               true,
	"default":           true,
	"deferrable":        true,
	"desc":              true,
	"distinct":          true,
	"do":                true,
	"else":              true,
	"end":               true,
	"except":            true,
	"false":             true,
	"fetch":             true,
	"filter":            true,
	"for":               true,
	"foreign":           true,
	"freeze":            true,
	"from":              true,
	"full":              true,
	"grant":             true,
	"group":             true,
	"having":            true,
	"hour":              true,
	"ilike":             true,
	"in":                true,
	"initially":         true,
	"inner":             true,
	"intersect":         true,
	"into":              true,
	"is":                true,
	"isnull":            true,
	"join":              true,
	"lateral":           true,
	"leading":           true,
	"left":              true,
	"like":              true,
	"limit":             true,
	"localtime":         true,
	"localtimestamp":    true,
	"minute":            true,
	"month":             true,
	"natural":           true,
	"not":               true,
	"notnull":           true,
	"null":              true,
	"offset":            true,
	"on":                true,
	"only":              true,
	"or":                true,
	"order":             true,
	"outer":             true,
	"over":              true,
	"overlaps":          true,
	"placing":           true,
	"primary":           true,
	"references":        true,
	"returning":         true,
	"right":             true,
	"second":            true,
	"select":            true,
	"session_user":      true,
	"similar":           true,
	"some":              true,
	"symmetric":         true,
	"table":             true,
	"tablesample":       true,
	"then":              true,
	"to":                true,
	"trailing":          true,
	"true":              true,
	"union":             true,
	"unique":            true,
	"user":              true,
	"using":             true,
	"variadic":          true,
	"varying":           true,
	"verbose":           true,
	"when":              true,
	"where":             true,
	"window":            true,
	"with":              true,
	"within":            true,
	"without":           true,
	"year":              true,
}
