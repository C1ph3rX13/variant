package main

import (
	"fmt"
	"variant/hook"
)

func main() {
	hashes, err := hook.AutoDumpHashes()
	if err != nil {
		return
	}

	for _, h := range hashes {
		fmt.Println(h.Format())
	}
}
