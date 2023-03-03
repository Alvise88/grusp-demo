package util

import "strings"

func ToCommand(cmd string) []string {
	return strings.Split(cmd, " ")
}
