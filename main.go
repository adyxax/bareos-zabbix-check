//+build !test

package main

import (
	"bareos-zabbix-check/zabbix"
	"fmt"
)

func main() {
	fmt.Print(zabbix.Main())
}
