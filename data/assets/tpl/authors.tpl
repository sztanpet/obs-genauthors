{{block "git" .Contributors}}Contributors:
{{range .}}{{.Name}}{{if .Email}} ({{.Email}}){{end}}
{{end}}
{{end}}
{{block "translator" .Translations}}Translations:
{{range .}}{{.Language}}
{{range .Translators}}    {{.Name}}{{if .Nick}} ({{.Nick}}){{end}}
{{end}}
{{end}}
{{end}}
