package utils

import (
	"fmt"
	"os"
	"strings"
)

func NewLine(params_optional ...string) {
	newLine := "1"
	if len(params_optional) > 0 {
		newLine = params_optional[0]
	}
	var str string
	if newLine == "1" {
		str = "\n"
	} else if newLine == "0" {
		str = ""
	}
	fmt.Printf(str)

}

func Color(params_optional ...string) {
	typeValue := ""
	option := "w"
	newLine := "1"
	if len(params_optional) > 0 {
		typeValue = params_optional[0]
	}
	if len(params_optional) > 1 {
		option = params_optional[1]
	}
	if len(params_optional) > 2 {
		newLine = params_optional[2]
	}
	var str string
	if newLine == "1" {
		str = "\n"
	} else if newLine == "0" {
		str = ""
	}
	option = strings.ToLower(option)
	if option == "e" || option == "d" {
		fmt.Print("\033[1;31m", typeValue, "\033[0m"+str)
	} else if option == "s" {
		fmt.Print("\033[1;32m", typeValue, " \033[0m"+str)
	} else if option == "w" {
		fmt.Print("\033[1;33m", typeValue, " \033[0m"+str)
	} else if option == "i" {
		fmt.Print("\033[1;36m", typeValue, " \033[0m"+str)
	} else if option == "nc" {
		fmt.Print("\033[0;37m", typeValue, " \033[0m"+str)
	}
}

func ColorTab(params_optional ...string) {
	typeValue := ""
	option := "w"
	newLine := "1"
	if len(params_optional) > 0 {
		typeValue = params_optional[0]
	}
	if len(params_optional) > 1 {
		option = params_optional[1]
	}
	if len(params_optional) > 2 {
		newLine = params_optional[2]
	}
	Color(typeValue, option, newLine)
}

func ColorTabExit(params_optional ...string) {
	typeValue := ""
	option := "w"
	if len(params_optional) > 0 {
		typeValue = params_optional[0]
	}
	if len(params_optional) > 1 {
		option = params_optional[1]
	}
	result := fmt.Sprintf("%s\t\n", typeValue)
	Color(result, option)
	os.Exit(0)
}

func Line(params_optional ...string) {
	option := "i"
	if len(params_optional) > 1 {
		option = params_optional[0]
	}
	ColorTab("=================================================================", option)
}
