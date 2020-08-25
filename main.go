package main

import (
	"github.com/jchenriquez/laundromat/cmd"
	"github.com/jchenriquez/laundromat/conf"
)

func main() {
	conf.SetDefaults()
	cmd.Execute()
}
