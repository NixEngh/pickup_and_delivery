package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/run"
)

func main() {
	problems, err := problem.LoadProblems("./data/input/")
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(problems, func(i, j int) bool {
		name1 := strings.Split(problems[i].Name, "_")[1]
		name2 := strings.Split(problems[j].Name, "_")[1]

		return name1 < name2
	})

	run.RunFinalExam(problems)
}
