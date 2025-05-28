package main

import (
	"strings"

	flag "github.com/spf13/pflag"
)

func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	from := []string{".", "_"}
	to := "-"
	for _, sep := range from {
		name = strings.ReplaceAll(name, sep, to)
	}
	return flag.NormalizedName(name)
}
