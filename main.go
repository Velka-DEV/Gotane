package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func checkProcess(args *CheckProcessArgs) CheckResult {
	time.Sleep(time.Second * 1)
	return CheckResultFree
}

func main() {

	lines := make([]string, 0)

	for i := 0; len(lines) < 10_000_000; i++ {
		lines = append(lines, strconv.FormatInt(int64(i), 10)+":"+RandStringBytesRmndr(4))
	}

	checker, err := NewCheckerBuilder().WithCombosFromLines(lines).SetCheckProcess(checkProcess).SetThreads(100).Build()

	if err != nil {
		panic(err)
	}

	start := time.Now()

	fmt.Printf("Started at %s \n", start.Format(time.RFC3339))

	go checker.Start()

	for {
		fmt.Printf("Progress: %d/%d | CPM: %d | ElapsedTime: %s \n", checker.Infos.GetChecked(), checker.Infos.GetTotal(), checker.Infos.GetCpm(), checker.Infos.GetElapsedTime().String())

		time.Sleep(time.Second)

		if checker.state == CheckerStateEnded {
			break
		}
	}

	end := time.Since(start)

	fmt.Printf("Free: %d, Hit: %d, Locked: %d, Invalid: %d, Error: %d, Total: %d, Time: %s", checker.Infos.GetFree(), checker.Infos.GetHit(), checker.Infos.GetLocked(), checker.Infos.GetInvalid(), checker.Infos.GetError(), checker.Infos.GetChecked(), end)
}
