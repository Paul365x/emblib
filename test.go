package main

import (
	"path/filepath"
	"strings"

	"github.com/emblib/adapters/jef"
	"github.com/emblib/adapters/pes_pec"
	"github.com/emblib/adapters/shared"
	"github.com/emblib/engine"
)

var file string = "designs/ATG12847.pes"

//"designs/D1124.pes"

//"designs/ATG12847.jef"

func main() {
	var pay *shared.Payload

	file_type := strings.ToLower(filepath.Ext(file))
	switch file_type {
	case ".pes":
		pay = pes_pec.Read_pes(file)
	case ".jef":
		pay = jef.Read_jef(file)
	}

	render := engine.NewEngine()
	render.Setup(engine.Fyne, pay)
	render.Run()
	render.Display()

}
