package main

var appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [OPTIONS] HOSTPATH CONTAINER:PATH

VERSION:
   {{.Version}}

OPTION:
  {{range .Flags}}{{.}}
  {{end}}
`
