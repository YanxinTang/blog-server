package utils

import "fmt"

func ColumnPlaceholder(args ...string) string {
	var column string
	length := len(args)
	for i := 0; i < length; i++ {
		if i == length-1 {
			column += fmt.Sprintf("`%s`", args[i])
		} else {
			column += fmt.Sprintf("`%s`, ", args[i])
		}
	}
	return column
}
