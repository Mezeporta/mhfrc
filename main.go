package main

import (
	"fmt"
	"hash/crc32"
	"mhfrc/jmp"
	"mhfrc/pac"
	"os"
)

const (
	crcJmp = 0xB894E6DF
	crcPac = 0x6FB2B035
)

func main() {
	fmt.Printf("Starting MHFReCompiler...\n")
	command := os.Args[1]
	switch command {
	case "decompile", "d":
		decompile()
	case "compile", "c":
		compile()
	default:
		fmt.Printf("Invalid command. Use 'decompile' or 'compile'.\n")
	}
}

func decompile() {
	fmt.Printf("Decompiling...\n")

	binJmp, err := os.ReadFile("mhfjmp.bin")
	if err != nil {
		fmt.Printf("Error reading mhfjmp.bin: %v\n", err)
	} else {
		crc := crc32.ChecksumIEEE(binJmp)
		if crc != crcJmp {
			fmt.Printf("Detected non-vanilla mhfjmp.bin\n")
		} else {
			fmt.Printf("Detected vanilla mhfjmp.bin\n")
		}

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

	binPac, err := os.ReadFile("mhfpac.bin")
	if err != nil {
		fmt.Printf("Error reading mhfpac.bin: %v\n", err)
	} else {
		crc := crc32.ChecksumIEEE(binPac)
		if crc != crcPac {
			fmt.Printf("Detected non-vanilla mhfpac.bin\n")
		} else {
			fmt.Printf("Detected vanilla mhfpac.bin\n")
		}

		json, err := pac.DecompilePac(binPac)
		if err != nil {
			fmt.Printf("Error decompiling mhfpac.bin: %v\n", err)
		} else {
			err = os.WriteFile("mhfpac.json", json, 0644)
			if err != nil {
				fmt.Printf("Error writing mhfpac.json: %v\n", err)
			} else {
				fmt.Printf("Decompiled mhfpac.bin to mhfpac.json\n")
			}
		}
	}
}

func compile() {
	fmt.Printf("Compiling...\n")

	jsonData, err := os.ReadFile("mhfjmp.json")
	if err != nil {
		fmt.Printf("Error reading mhfjmp.json: %v\n", err)
	} else {
		jmpData, err := jmp.CompileJmp(jsonData)
		if err != nil {
			fmt.Printf("Error compiling mhfjmp.json: %v\n", err)
		} else {
			err = os.WriteFile("mhfjmp_recomp.bin.", jmpData, 0644)
			if err != nil {
				fmt.Printf("Error writing mhfjmp_recomp.bin: %v\n", err)
			} else {
				fmt.Printf("Compiled mhfjmp.json to mhfjmp_recomp.bin\n")
			}
		}
	}
}
