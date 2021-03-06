package main

import (
	"errors"
	"github.com/abcum/lcp"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"text/template"
)

var displayNameRegExp = regexp.MustCompile(`display:\"(.*)\"`)
var transformUpperRegExp = regexp.MustCompile(`(?i)transform:\s*upper`)
var transformLowerRegExp = regexp.MustCompile(`(?i)transform:\s*lower`)
var transformNoneRegExp = regexp.MustCompile(`(?i)transform:\s*none`)

var enumHeaderTemplate *template.Template
var enumItemTemplate *template.Template

func init() {

	enumHeaderTemplate = template.Must(template.New("enumheader").Parse(enumHeaderTemplateText))
	enumItemTemplate = template.Must(template.New("enumitem").Parse(enumItemTemplateText))

}

type transform int

const (
	transformNone transform = iota
	transformUpper
	transformLower
)

type enum struct {
	PackageName string
	Values      []string
	//OverrideDisplayName contains a map of the Value string to override
	//value, if it exists. If it is in the map with value "" then it has been
	//overridden to have that value. If it is not in the map then it should be
	//default.
	OverrideDisplayName map[string]string
	Transform           map[string]transform
	DefaultTransform    transform
}

//findEnums processes the package at packageName and returns a list of enums
//that should be processed (that is, they have the magic comment)
func findEnums(inputPackageName string) (enums []*enum, err error) {

	packageASTs, err := parser.ParseDir(token.NewFileSet(), inputPackageName, nil, parser.ParseComments)

	if err != nil {
		return nil, errors.New("Parse error: " + err.Error())
	}

	for packageName, theAST := range packageASTs {
		for _, file := range theAST.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)

				if !ok {
					//Guess it wasn't a genDecl at all.
					continue
				}

				if genDecl.Tok != token.CONST {
					//We're only interested in Const decls.
					continue
				}

				if !enumConfig(genDecl.Doc.Text()) {
					//Must not have found the magic comment in the docs.
					continue
				}

				defaultTransform := configTransform(genDecl.Doc.Text(), transformNone)

				theEnum := &enum{
					PackageName:         packageName,
					OverrideDisplayName: make(map[string]string),
					Transform:           make(map[string]transform),
					DefaultTransform:    defaultTransform,
				}

				for _, spec := range genDecl.Specs {

					valueSpec, ok := spec.(*ast.ValueSpec)

					if !ok {
						//Guess it wasn't a valueSpec after all!
						continue
					}

					if len(valueSpec.Names) != 1 {
						return nil, errors.New("Found an enum that had more than one name on a line. That's not allowed for now.")
					}

					valueName := valueSpec.Names[0].Name

					theEnum.Values = append(theEnum.Values, valueName)

					if hasOverride, displayName := overrideDisplayname(valueSpec.Doc.Text()); hasOverride {
						theEnum.OverrideDisplayName[valueName] = displayName
					}

					theEnum.Transform[valueName] = configTransform(valueSpec.Doc.Text(), defaultTransform)

				}

				if len(theEnum.Values) > 0 {
					enums = append(enums, theEnum)
				}

			}
		}
	}

	return enums, nil
}

//outputForEnums takes the found enums and produces the output string
//representing the un-formatted go code to generate for those enums.
func outputForEnums(enums []*enum) (enumOutput string, err error) {
	for _, enum := range enums {

		if enumOutput == "" {
			enumOutput = enumHeaderForPackage(enum.PackageName)
		}

		var literals [][]byte

		for _, literal := range enum.Values {
			if !fieldNamePublic(literal) {
				continue
			}
			literals = append(literals, []byte(literal))
		}

		if len(literals) == 0 {
			return "", errors.New("No public literals in enum")
		}

		prefix := string(lcp.LCP(literals...))

		if len(prefix) == 0 {
			return "", errors.New("Enum with autoreader configured didn't have a common prefix.")
		}

		values := make(map[string]string, len(literals))

		i := 0

		for _, literal := range enum.Values {
			if !strings.HasPrefix(literal, prefix) {
				return "", errors.New("enum literal didn't have prefix we thought it did")
			}

			//If there's an override deisplay name, use that
			displayName, ok := enum.OverrideDisplayName[literal]

			//If there wasn't an override, do the default. Note that an
			//override "" that is in the map is legal.
			if !ok {
				displayName = titleCaseToWords(strings.Replace(literal, prefix, "", -1))

				switch enum.Transform[literal] {
				case transformLower:
					displayName = strings.ToLower(displayName)
				case transformUpper:
					displayName = strings.ToUpper(displayName)
				}
			}

			values[literal] = displayName
			i++
		}

		enumOutput += enumItem(prefix, values)

	}

	return enumOutput, nil
}

var titleCaseReplacer *strings.Replacer

//titleCaseToWords writes "ATitleCaseString" to "A Title Case String"
func titleCaseToWords(in string) string {

	if titleCaseReplacer == nil {

		var replacements []string

		for r := 'A'; r <= 'Z'; r++ {
			str := string(r)
			replacements = append(replacements, str)
			replacements = append(replacements, " "+str)
		}

		titleCaseReplacer = strings.NewReplacer(replacements...)

	}

	return strings.TrimSpace(titleCaseReplacer.Replace(in))

}

func processEnums(packageName string) (enumOutput string, err error) {
	enums, err := findEnums(packageName)

	if err != nil {
		return "", errors.New("Couldn't parse for enums: " + err.Error())
	}

	if len(enums) == 0 {
		//No enums. That's totally legit.
		return "", nil
	}

	output, err := outputForEnums(enums)

	if err != nil {
		return "", errors.New("Couldn't generate output for enums: " + err.Error())
	}

	return output, nil

}

func enumConfig(docLines string) bool {

	for _, docLine := range strings.Split(docLines, "\n") {
		docLine = strings.ToLower(docLine)
		docLine = strings.TrimPrefix(docLine, "//")
		docLine = strings.TrimSpace(docLine)
		if strings.HasPrefix(docLine, magicDocLinePrefix) {
			return true
		}
	}

	return false
}

func configTransform(docLines string, defaultTransform transform) transform {
	for _, line := range strings.Split(docLines, "\n") {
		if transformLowerRegExp.MatchString(line) {
			return transformLower
		}
		if transformUpperRegExp.MatchString(line) {
			return transformUpper
		}
		if transformNoneRegExp.MatchString(line) {
			return transformNone
		}
	}

	return defaultTransform
}

func overrideDisplayname(docLines string) (hasOverride bool, displayName string) {
	for _, line := range strings.Split(docLines, "\n") {
		result := displayNameRegExp.FindStringSubmatch(line)

		if len(result) == 0 {
			continue
		}

		if len(result[0]) == 0 {
			continue
		}
		if len(result) != 2 {
			continue
		}

		//Found it! Even if the matched expression is "", that's fine. if
		//there are quoted strings that's fine, because that's exactly how
		//they should be output at the end.
		return true, result[1]

	}

	return false, ""
}

func enumHeaderForPackage(packageName string) string {
	return templateOutput(enumHeaderTemplate, map[string]string{
		"packageName": packageName,
	})
}

func enumItem(prefix string, values map[string]string) string {
	return templateOutput(enumItemTemplate, map[string]interface{}{
		"prefix": prefix,
		"values": values,
	})
}

const enumHeaderTemplateText = `/************************************
 *
 * This file contains auto-generated methods to help configure enums. 
 * It was generated by autoreader.
 *
 * DO NOT EDIT by hand.
 *
 ************************************/

package {{.packageName}}

import (
	"github.com/jkomoros/boardgame/enum"
)

var Enums = enum.NewSet()

`

const enumItemTemplateText = `var {{.prefix}}Enum = Enums.MustAdd("{{.prefix}}", map[int]string{
	{{ $prefix := .prefix -}}
	{{range $name, $value := .values -}}
	{{$name}}: "{{$value}}",
	{{end}}
})

`
