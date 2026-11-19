package steganography

import (
	"fmt"
	"github.com/ANG13T/DroneXtract/assets"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/TwiN/go-color"
)

func ExecuteSteganography(index int) {
	printfStgMenu(index)
	output := helpers.OptionWithBack(0, 5, false)
	if output == -1 {
		return
	}
	inputFile := helpers.FileInputString()
	switch in := index; in {
	case 1:
		suite := NewDJI_EXIF_Parser(inputFile)
		suite.ExecuteEXIFAnalysis(output)
	case 2:
		suite := NewDJI_DNG_Parser(inputFile)
		if suite == nil {
			return
		}
		suite.ExecuteDNGAnalysis(output)
	case 3:
		suite := NewDJI_SRT_Parser(inputFile)
		if suite == nil {
			return
		}
		suite.ExecuteSRTAnalysis(output)
	case 4:
		suite := NewDJI_XMP_Parser(inputFile)
		if suite == nil {
			return
		}
		suite.ExecuteXMPAnalysis(output)
	}
}

// printStgMenu 打印对应子分析的双语菜单。
func printfStgMenu(index int) {
	var parsingBanners = []string{"exif.txt", "dng.txt", "srt.txt", "xmp.txt"}
	fmt.Println(color.Ize(color.Cyan, assets.Text("steganography/"+parsingBanners[index-1])))
}
