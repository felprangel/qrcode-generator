package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/felipe/qrcode-generator/internal/config"
	"github.com/felipe/qrcode-generator/internal/qr"
)

const usage = `usage:
  qrcode-generator <number> [<number>...]
  qrcode-generator config set KEY=VALUE
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)

		return 1
	}

	if args[0] == "config" {
		return runConfig(args[1:], stderr)
	}

	return runGenerate(args, stdout, stderr)
}

func runGenerate(args []string, stdout, stderr io.Writer) int {
	for _, arg := range args {
		if !isDigits(arg) {
			fmt.Fprintf(stderr, "error: invalid argument %q (must be numeric)\n", arg)

			return 1
		}
	}

	cfg, err := config.Load()

	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)

		return 1
	}

	for _, number := range args {
		fmt.Fprintf(stdout, "--- %s ---\n", number)
		qr.Render(stdout, cfg.Prefix+number)
		fmt.Fprintln(stdout)
	}

	return 0
}

func runConfig(args []string, stderr io.Writer) int {
	if len(args) != 2 || args[0] != "set" {
		fmt.Fprint(stderr, usage)

		return 1
	}

	keyValue := args[1]

	equalIndex := strings.IndexByte(keyValue, '=')

	if equalIndex <= 0 {
		fmt.Fprintf(stderr, "error: expected KEY=VALUE, got %q\n", keyValue)

		return 1
	}

	key := keyValue[:equalIndex]
	value := keyValue[equalIndex+1:]

	err := config.Set(key, value)

	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)

		return 1
	}

	return 0
}

func isDigits(str string) bool {
	if str == "" {
		return false
	}

	for _, rune := range str {
		if rune < '0' || rune > '9' {
			return false
		}
	}

	return true
}
