package generator

import (
	"fmt"
)

// every input type needs a persist type

// every method needs a query string method
// every method needs a generated interface satisfied by the input
// every method needs a default implementation
// every method needs a PersistImpl handler

// every service needs a PersistImpl declaration
// every service needs a PersistHelper struct declaration
// every service needs a Merge function
// every service needs a default opts function

// just a string, but represents the path that these persist
// files will live in
type PackagePath string

type PersistPackage struct {
	// the files seperated into packages they will generate into
	files map[PackagePath][]*FileStruct
}

func NewPersistPackage(fileList *FileList) *PersistPackage {
	grouped := make(map[PackagePath][]*FileStruct)

	for _, f := range *fileList {
		if !f.NeedsPersistLibDir() {
			continue
		}
		p := PackagePath(f.GetPersistLibFullFilepath().path)

		grouped[p] = append(grouped[p], f)
	}

	return &PersistPackage{files: grouped}
}

type PersistContent struct {
	Content string // the contents that belong in the code generator response
	Name    string // the filename for the response
}

func (p *PersistPackage) Generate() []PersistContent {
	var contents []PersistContent
	for pkg, files := range p.files {
		// files is all the files that need to be in this package.
		contents = append(contents, GeneratePkgLevelContent(pkg, files))
		for _, file := range files {
			contents = append(contents, GenerateFileQueryContent(pkg, file))
			contents = append(contents, GenerateFilePersistHandlerContent(pkg, file))
		}
	}
	return contents
}
func GenerateFileQueryContent(pkg PackagePath, file *FileStruct) PersistContent {
	stringer := PersistStringer{}
	p := file.GetPersistLibFullFilepath()
	content := PersistContent{
		Name:    fmt.Sprintf("%s/%s_queries.persist.go", pkg, p.filename),
		Content: "package persist_lib\n import \"cloud.google.com/go/spanner\"\n",
	}
	for _, s := range *file.ServiceList {
		if !s.IsSpanner() {
			continue
		}
		for _, m := range *s.Methods {
			content.Content += stringer.QueryFunction(m)
		}
	}
	for _, s := range *file.ServiceList {
		if !s.IsSpanner() {
			continue
		}
		for _, m := range *s.Methods {
			content.Content += stringer.QueryInterfaceDefinition(m)
		}
	}
	return content
}
func GenerateFilePersistHandlerContent(pkg PackagePath, file *FileStruct) PersistContent {
	stringer := PersistStringer{}
	p := file.GetPersistLibFullFilepath()
	content := PersistContent{
		Name: fmt.Sprintf("%s/%s_query_handlers.persist.go", pkg, p.filename),
		Content: fmt.Sprintf(
			"package persist_lib\nimport(\n\"%s\"\n\"%s\"\n)\n",
			"cloud.google.com/go/spanner",
			"golang.org/x/net/context",
		),
	}
	for _, s := range *file.ServiceList {
		if !s.IsSpanner() {
			continue
		}
		content.Content += stringer.HandlersStructDeclaration(s)
		content.Content += stringer.HelperFunctionImpl(s)
		content.Content += stringer.DefaultFunctionsImpl(s)
	}
	return content
}
func GeneratePkgLevelContent(pkg PackagePath, files []*FileStruct) PersistContent {
	stringer := PersistStringer{}
	content := PersistContent{
		Name: string(pkg) + "/pkg_level_definitions.persist.go",
	}
	// imports := Imports(make([]*Import, 0))
	haveProcessed := make(map[string]bool)
	// addToImports := func(m *Method) {
	// 	if opt := m.GetMethodOption(); opt != nil {
	// 		if mapping := opt.GetMapping(); mapping != nil {
	// 			for _, t := range mapping.GetTypes() {
	// 				imports.GetOrAddImport(GetGoPackage(t.GetGoPackage()), GetGoPath(t.GetGoPackage()))
	// 			}
	// 		}
	// 	}
	// }
	for _, f := range files {
		for _, s := range *f.ServiceList {
			for _, m := range *s.Methods {
				if !haveProcessed[NewPLInputName(m).String()] {
					haveProcessed[NewPLInputName(m).String()] = true
					content.Content += stringer.MessageInputDeclaration(m)
					// addToImports(m)
				}
			}
		}
	}
	p := &Printer{}
	p.P("package persist_lib\n")
	p.P("import \"cloud.google.com/go/spanner\"\n")
	p.P("%s", stringer.DeclareSpannerGetter())
	content.Content = p.String() + content.Content

	return content
}
