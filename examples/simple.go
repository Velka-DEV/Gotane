package main

import (
	"fmt"

	"github.com/Velka-DEV/Gotane/pkg/core"
)

func checkProcess(args *core.CheckProcessArgs) *core.CheckResult {
	fmt.Println(args.Combo.String())

	return &core.CheckResult{
		Combo:  args.Combo,
		Status: core.CheckStatusFree,
		Captures: map[string]interface{}{
			"points": "49",
		},
	}
}

func main() {

	combos, err := core.LoadCombosFromFile("combos.txt")

	if err != nil {
		panic(err)
	}

	checker := core.NewCheckerBuilder().WithConfig(&core.CheckerConfig{
		Threads:       300,
		OutputFree:    true,
		OutputInvalid: false,
		OutputLocked:  true,
		OutputUnknown: true,
	}).WithCheckProcess(checkProcess).WithCombos(combos).Build()

	checker.StartAndWait()
}
