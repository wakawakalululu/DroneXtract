package parsing

import (
	"fmt"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/TwiN/go-color"
)

func ExecuteParser(index int) {
	filePath := helpers.FileInputString()
	switch in := index; in {
	case 1:
		suite := NewDJI_CSV_Parser(filePath)
		suite.ParseContents()
	case 2:
		suite := NewDJI_KML_Parser(filePath)
		suite.ParseContents()
	case 3:
		suite := NewDJI_GPX_Parser(filePath)
		suite.ParseContents()
	}
}

func GenTableHeader(name string, containBreak bool) {
	if containBreak {
		fmt.Println(color.Ize(color.Blue, "\n    ╔══════════════════════════════════════════════════════════════════════════════╗"))
	} else {
		fmt.Println(color.Ize(color.Blue, "    ╔══════════════════════════════════════════════════════════════════════════════╗"))
	}
	fmt.Println(color.Ize(color.Blue, "    ║"+helpers.CenterText(name, 78)+"║"))
	fmt.Println(color.Ize(color.Blue, "    ╠══════════════════════════════════════════════════════════════════════════════╣"))
}

func GenTableHeaderModified(name string) {
	fmt.Println(color.Ize(color.Blue, "    ╠══════════════════════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Blue, "    ║"+helpers.CenterText(name, 78)+"║"))
	fmt.Println(color.Ize(color.Blue, "    ╠══════════════════════════════════════════════════════════════════════════════╣"))
}
