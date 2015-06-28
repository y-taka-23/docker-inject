package main

var appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [OPTIONS] HOSTDIR CONTAINER:PATH

VERSION:
   {{.Version}}

OPTION:
  {{range .Flags}}{{.}}
  {{end}}
`
