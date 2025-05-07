package main

import (
	"fmt"
	"mhfrc/jmp"
	"os"
)

func main() {
	fmt.Printf("Starting MHFReCompiler...")
	decompile()
}

func decompile() {
	fmt.Printf("Decompiling...\n")

	binJmp, err := os.ReadFile("mhfjmp.bin")
	if err != nil {
		fmt.Printf("Error reading mhfjmp.bin: %v\n", err)
	} else {
		json, err := jmp.DecompileJmp(binJmp)
		if err != nil {
			fmt.Printf("Error decompiling mhfjmp.bin: %v\n", err)
		} else {
			err = os.WriteFile("mhfjmp.json", json, 0644)
			if err != nil {
				fmt.Printf("Error writing mhfjmp.json: %v\n", err)
			} else {
				fmt.Printf("Decompiled mhfjmp.bin to mhfjmp.json\n")
			}
		}
	}
}
