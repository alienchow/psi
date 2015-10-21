package main

import (
	"fmt"
	"os"

	"github.com/alienchow/psi/psi"
	"github.com/alienchow/psi/region"
)

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	r := region.GetByArg(arg)

	reading := psi.NewReading()
	if err := reading.Refresh(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reading.Get(r))
}
