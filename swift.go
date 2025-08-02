package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

func generateSwiftFiles() (files map[string][]byte) {
	files = make(map[string][]byte)
	un := make(map[string]bool)
	for _, x := range fields {
		for _, k := range x.Swift.Files {
			un[k] = true
		}
	}

	switch len(un) == 0 {
	case true:
		files[settings.File.swift] = generateSwiftCode()
	default:
		files[settings.File.swift] = generateSwiftCode()
		for k := range un {
			fn := filepath.Join(settings.File.path, k+".swift")
			files[fn] = generateSwiftCode(k)
		}
	}
	return

}

// insert sql query
func generateSwiftCode(filename ...string) []byte {

	var structname string
	if filename == nil {
		structname = settings.Go.StructName + "Model"
	} else {
		structname = filename[0] + "Model"
	}

	/*
		struct channelModel: Identifiable, Decodable, Encodable {
		   var id: Int
		   var premium: Bool? = false
		   var soon: Bool? = false
		   var level: Int
		   var title: String
		   var line: String? = ""
		   var lang: String
		   var category: String? = ""
		   var sub: String? = ""
		   var followers: Int? = 0
		   var news: Int? = 0
		}
	*/

	var list []string

	list = append(list, fmt.Sprintf("struct %s: Identifiable, Decodable, Encodable {", structname))

	var max int
	for _, x := range fields {
		n := "var " + x.Name
		if max < len(n) {
			max = len(n)
		}
	}
	max = max + 2

	var c int
	for _, x := range fields {
		if x.Swift.Skip {
			continue
		}
		if filename != nil && !exist(x.Swift.Files, filename[0]) {
			continue
		}
		c++
		var tsType string
		switch x.Type {
		case "int", "int64", "uint", "uint64", "time.Duration", "int16", "int32", "uint32", "int8", "uint8":
			tsType = "Int" + reqswift(x.Swift.Must, "= 0")
		case "float64", "float32":
			tsType = "Double" + reqswift(x.Swift.Must, "= 0")
		case "bool":
			tsType = "Bool" + reqswift(x.Swift.Must, "= false")
		case "string", "byte":
			tsType = `String` + reqswift(x.Swift.Must, `= ""`)
		case "map[string]string":
			tsType = "[String:String]" + reqswift(x.Swift.Must, "= [String:String]()")
		case "map[int]string":
			tsType = "[Int:String]" + reqswift(x.Swift.Must, "= [Int:String]()")
		case "map[int]int":
			tsType = "[Int:Int]" + reqswift(x.Swift.Must, "= [Int:Int]()")
		case "map[string]int":
			tsType = "[String:Int]" + reqswift(x.Swift.Must, "= [String:Int]()")
		case "map[string]bool":
			tsType = "[String:Bool]" + reqswift(x.Swift.Must, "= [String:Bool]()")
		case "[]string":
			tsType = "[String]" + reqswift(x.Swift.Must, "= [String]()")
		case "[]int", "[]int64", "[]uint", "[]uint64", "[]time.Duration", "[]int16", "[]int32", "[]uint32", "[]int8":
			tsType = "[Int]" + reqswift(x.Swift.Must, "= [Int]()")
		case "[]uint8", "[]byte":
			tsType = "Data" + reqswift(x.Swift.Must, "= Data()")
		case "map[string]interface{}", "map[string]any":
			tsType = "[String:Any]" + reqswift(x.Swift.Must, "= [String:Any]()")
		default:
			tsType = "Data" + reqswift(x.Swift.Must, "= Data()")
		}

		if x.Swift.Type != "" {
			tsType = x.Swift.Type
		}

		//fmt.Println("swift", x.Type, ":", tsType)

		line := "\t" + tools.FormatLineSwift("var "+x.Name, ": "+tsType, max)
		line += " // " + x.Type

		if x.Comment != "" {
			line += " " + x.Comment
		}
		list = append(list, line)

	}

	res := strings.Join(list, "\n") + "\n}"
	if c == 0 {
		return nil
	}
	return []byte(res)
}

func reqswift(rec bool, defaults string) string {
	if !rec {
		return "? " + defaults
	}
	return " " + defaults
}

func exist(list []string, id string) bool {
	for _, x := range list {
		if x == id {
			return true
		}
	}
	return false
}

/*
//exercise

	enum logo: Int {
	   case me, link, zoom, skype, meet, teams, telegram, slack, youtube
	   var icon: Image {
	        get {
	           switch self{
	            case .me:         return Image("me").resizable()
	            case .link:       return Image("link").resizable()
	            case .zoom:       return Image("zoom").resizable()
	            case .skype:      return Image("skype").resizable()
	            case .meet:       return Image("meet").resizable()
	            case .teams:      return Image("teams").resizable()
	            case .telegram:   return Image("telegram").resizable()
	            case .slack:      return Image("slack").resizable()
	            case .youtube:    return Image("youtube").resizable()
	           }
	        }
	   }
	}
*/
func generateSwiftEnum() (v []byte) {

	var list []string
	list = append(list, "//enum list")
	list = append(list, fmt.Sprintf("enum %s: String {\n", settings.Swift.Enum))

	// case
	var caselist []string //id
	var namelist []string
	var titlelist []string

	var max int
	for _, x := range fields {
		var name string
		switch x.Json.Name == "" {
		case true:
			name = x.Name
		case false:
			name = x.Json.Name
		}
		name = fmt.Sprintf("case .%s:", name)
		if max < len(name) {
			max = len(name)
			fmt.Println(max)
		}
	}

	max = max + 2

	for _, x := range fields {

		if x.Swift.Enum.Skip {
			continue
		}

		var name string
		switch x.Json.Name == "" {
		case true:
			name = x.Name
		case false:
			name = x.Json.Name
		}

		var title string
		switch x.Go.Title != "" {
		case true:
			title = x.Go.Title
		case false:
			title = x.Go.Name
		}

		caselist = append(caselist, name)
		namelist = append(namelist, "\t\t\t\t"+tools.FormatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return "%s"`, name), max))
		titlelist = append(titlelist, "\t\t\t\t"+tools.FormatLineSwift(fmt.Sprintf("case .%s:", name), fmt.Sprintf(`return "%s"`, title), max))
	}

	// case
	list = append(list, "\t//case")
	list = append(list, fmt.Sprintf("\tcase %s", strings.Join(caselist, ", ")))
	list = append(list, "")

	// name
	list = append(list, "\t//name")
	list = append(list, "\tvar string: String {")
	list = append(list, "\t\tget {")
	list = append(list, "\t\t\tswitch self {")
	list = append(list, fmt.Sprintf("%s", strings.Join(namelist, "\n")))
	list = append(list, "\t\t\t}")
	list = append(list, "\t\t}")
	list = append(list, "\t}")
	list = append(list, "\t")

	// title
	list = append(list, "\t//title")
	list = append(list, "\tvar title: String {")
	list = append(list, "\t\tget {")
	list = append(list, "\t\t\tswitch self {")
	list = append(list, fmt.Sprintf("%s", strings.Join(titlelist, "\n")))
	list = append(list, "\t\t\t}")
	list = append(list, "\t\t}")
	list = append(list, "\t}")
	list = append(list, "\t")

	// end
	list = append(list, "}\n")
	return []byte(strings.Join(list, "\n"))
}
