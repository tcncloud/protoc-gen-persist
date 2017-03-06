package files

var (
	TemplateList = map[string]string{
		// import template
		"import_template": `{{define "import_template"}}
import(
	{{range $import := .}}
	{{$import.GoPackageName}} "{{$import.GoImportPath}}"
	{{end}}
)
{{end}}`,
		// Implement structs
		"implement_structs": `{{define "implement_structs"}}
{{end}}`,
		// implement services
		"implement_services": `{{define "implement_services"}}
{{range $srv := .}}
{{if $srv.IsServiceEnabled}}
type {{$srv.GetName}}Impl struct {
	{{if $srv.IsSQL}}SQLDB *sql.DB{{end}}
	{{if $srv.IsMongo}}MongoDB *mgo.DB{{end}}
}
{{range $method := $srv.Methods}}
{{template "implement_method" $method}}
{{end}}
{{end}}
{{end}}
{{end}}`,
		// Implement method
		"implement_method": `{{define "implement_method"}}
{{if .IsUnary}} {{template "unary_method" .}} {{end}}
{{if .IsClientStreaming}} {{template "client_streaming_method" .}} {{end}}
{{if .IsServerStreaming}} {{template "server_streaming_method" .}} {{end}}
{{if .IsBidiStreaming}} {{template "bidi_method" .}} {{end}}
{{end}}`,
		"unary_method": `{{define "unary_method"}}
{{if .IsSQL}}{{template "sql_unary_method" .}}{{end}}
{{if .IsMongo}}{{template "mongo_unary_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_unary_method" .}}{{end}}
{{end}}`,
		"client_streaming_method": `{{define "client_streaming_method"}}
{{if .IsSQL}}{{template "sql_client_streaming_method" .}}{{end}}
{{if .IsMongo}}{{template "mongo_client_streaming_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_client_streaming_method" .}}{{end}}
{{end}}`,
		"server_streaming_method": `{{define "server_streaming_method"}}
{{if .IsSQL}}{{template "sql_server_streaming_method" .}}{{end}}
{{if .IsMongo}}{{template "mongo_server_streaming_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_server_streaming_method" .}}{{end}}
{{end}}`,
		"bidi_method": `{{define "bidi_method"}}
{{if .IsSQL}}{{template "sql_bidi_streaming_method" .}}{{end}}
{{if .IsMongo}}{{template "mongo_bidi_streaming_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_bidi_streaming_method" .}}{{end}}
{{end}}`,
		"sql_unary_method":                `{{define "sql_unary_method"}}// sql unary {{.GetName}} unimplemented{{end}}`,
		"sql_client_streaming_method":     `{{define "sql_client_streaming_method"}}// sql client streaming {{.GetName}} unimplemented{{end}}`,
		"sql_server_streaming_method":     `{{define "sql_server_streaming_method"}}// sql server streaming {{.GetName}} unimplemented{{end}}`,
		"sql_bidi_streaming_method":       `{{define "sql_bidi_streaming_method"}}// sql bidi streaming {{.GetName}} unimplemented{{end}}`,
		"mongo_unary_method":              `{{define "mongo_unary_method"}}// mongo unary {{.GetName}} unimplemented{{end}}`,
		"mongo_client_streaming_method":   `{{define "mongo_client_streaming_method"}}// mongo client streaming {{.GetName}} unimplemented{{end}}`,
		"mongo_server_streaming_method":   `{{define "mongo_server_streaming_method"}}// mongo server streaming {{.GetName}} unimplemented{{end}}`,
		"mongo_bidi_streaming_method":     `{{define "mongo_bidi_streaming_method"}}// mongo bidi streaming {{.GetName}} unimplemented{{end}}`,
		"spanner_unary_method":            `{{define "spanner_unary_method"}}// spanner unary {{.GetName}} unimplemented{{end}}`,
		"spanner_client_streaming_method": `{{define "spanner_client_streaming_method"}}// spanner client streaming {{.GetName}} unimplemented{{end}}`,
		"spanner_server_streaming_method": `{{define "spanner_server_streaming_method"}}// spanner server streaming {{.GetName}} unimplemented{{end}}`,
		"spanner_bidi_streaming_method":   `{{define "spanner_bidi_streaming_method"}}// spanner bidi streaming {{.GetName}} unimplemented{{end}}`,
	}
)
