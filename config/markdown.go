package config

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var MD = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
		extension.TaskList,
		extension.Strikethrough,
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)
