//+build !test

package main

import (
	"fmt"

	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/zabbix"
)

func main() {
	fmt.Print(zabbix.Main())
}
