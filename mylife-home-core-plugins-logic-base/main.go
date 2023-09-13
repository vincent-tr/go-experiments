//go:generate go run ../mylife-home-core-generator/cmd/main.go .

package plugin_entry

import (
	"fmt"
	_ "mylife-home-core-plugins-logic-base/plugin"
)

func init() {
	fmt.Println("mod load logic-base")
}
