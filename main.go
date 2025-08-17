package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	INS_INVALID uint8 = iota
	INS_WRITE
	INS_REPLACE
	INS_GOTO
	INS_OPEN
	INS_CLOSE
	INS_NEXT_PARAM uint8 = 0xFD
	INS_FINALIZE   uint8 = 0xFF
)

var instructionMap = map[string]uint8{
	"open":  INS_OPEN,
	"close": INS_CLOSE,
	"wri":   INS_WRITE,
	"rep":   INS_REPLACE,
	"goto":  INS_GOTO,
}

const (
	STATE_INS int16 = iota
)

func startsWith(str string, toTest string) bool {
	if str[:len(toTest)] == toTest {
		return true
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Wrong argument count! %d passed, 1 required!\n", len(os.Args)-1)
		return
	}

	if _, err := os.Stat(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	compiledBytes := []byte{0xAE, 'p'}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}
		tokens := strings.Split(line, " ")

		instruction := instructionMap[tokens[0]]
		parameters := tokens[1:]
		compiledBytes = append(compiledBytes, byte(instruction))
		for i, param := range parameters {
			paramBytes := []byte(param)
			if startsWith(param, "0x") {
				paramBytes, err = hex.DecodeString(param[2:])
				if err != nil {
					fmt.Println(err)
					return
				}
				slices.Reverse(paramBytes)
			}
			compiledBytes = append(compiledBytes, byte(len(paramBytes)))
			compiledBytes = append(compiledBytes, paramBytes...)
			if i < len(parameters)-1 {
				compiledBytes = append(compiledBytes, INS_NEXT_PARAM)
			}
		}
		compiledBytes = append(compiledBytes, INS_FINALIZE)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile("result.bin", compiledBytes, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
