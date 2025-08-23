package main

import (
	"fmt"
	"reflect"
	"strings"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/charmbracelet/lipgloss"
)

var (
	grey  = lipgloss.Color("#7d7d7d")
	green = lipgloss.Color("#a6e3a1")
	blue  = lipgloss.Color("#7287fd")

	greyStyle  = lipgloss.NewStyle().Foreground(grey)
	greenStyle = lipgloss.NewStyle().Foreground(green)
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(blue)
)

func printStruct(settings github_recon_settings.Settings, s any, indent int) {
	if settings.Silent {
		return
	}

	prefix := strings.Repeat("  ", indent)

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		if v.NumField() == 0 {
			fmt.Println(prefix + greyStyle.Render("No data found"))
			fmt.Println("")
			return
		}

		printed := 0

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i).Name
			value := v.Field(i)

			if !value.IsValid() || (value.Kind() == reflect.String && value.String() == "") {
				continue
			}
			if value.Kind() == reflect.String && value.String() == "0001-01-01 00:00:00 +0000 UTC" {
				continue
			}
			if (field == "FirstFoundIn" || field == "FoundIn") && !settings.ShowSource {
				continue
			}
			printed++

			switch value.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Ptr:
				fmt.Println(prefix + greyStyle.Render(field+":"))
				printStruct(settings, value.Interface(), indent+1)
			case reflect.String:
				fmt.Printf("%s%s %s\n", prefix, greyStyle.Render(field+":"), greenStyle.Render(fmt.Sprintf("%q", value.Interface())))
			default:
				fmt.Printf("%s%s %s\n", prefix, greyStyle.Render(field+":"), greenStyle.Render(fmt.Sprintf("%v", value.Interface())))
			}
		}
		if printed == 0 {
			fmt.Println(prefix + greyStyle.Render("No data found"))
		}
		fmt.Println("")

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			fmt.Println(prefix + greyStyle.Render("No data found"))
			fmt.Println("")
			return
		}
		for i := 0; i < v.Len(); i++ {
			printStruct(settings, v.Index(i).Interface(), indent)
		}

	default:
		fmt.Println(prefix + greenStyle.Render(fmt.Sprintf("%v", v.Interface())))
		fmt.Println("")
	}
}

func header() {
	asciiArt := "        __                       \n  ___ _/ /  _______ _______  ___ \n / _ `/ _ \\/ __/ -_) __/ _ \\/ _ \\\n \\_, /_//_/_/  \\__/\\__/\\___/_//_/\n/___/                            "

	grey := lipgloss.Color("#7d7d7d")

	greyStyle := lipgloss.NewStyle().Foreground(grey)
	fmt.Println(
		greyStyle.Render(lipgloss.JoinVertical(lipgloss.Right, asciiArt, "@anotherhadi\n")),
	)
}

func printTitle(silent bool, title string) {
	if silent {
		return
	}
	fmt.Println(titleStyle.Render(title) + "\n")
}
