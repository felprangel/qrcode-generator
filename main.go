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
	for _, a := range args {
		if !isDigits(a) {
			fmt.Fprintf(stderr, "error: invalid argument %q (must be numeric)\n", a)
			return 1
		}
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	for _, n := range args {
		fmt.Fprintf(stdout, "--- %s ---\n", n)
		qr.Render(stdout, cfg.URLPrefix+n)
		fmt.Fprintln(stdout)
	}
	return 0
}

func runConfig(args []string, stderr io.Writer) int {
	if len(args) != 2 || args[0] != "set" {
		fmt.Fprint(stderr, usage)
		return 1
	}
	kv := args[1]
	eq := strings.IndexByte(kv, '=')
	if eq <= 0 {
		fmt.Fprintf(stderr, "error: expected KEY=VALUE, got %q\n", kv)
		return 1
	}
	key := kv[:eq]
	value := kv[eq+1:]
	if err := config.Set(key, value); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
