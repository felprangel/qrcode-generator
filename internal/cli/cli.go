// Package cli wires argument parsing to config and qr rendering.
package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/felipe/qrcode-generator/internal/config"
	"github.com/felipe/qrcode-generator/internal/qr"
)

const usage = `usage:
  qrcode-generator [-c|--clear] [-p|--preset NAME | --no-preset] <text> [<text>...]
  qrcode-generator config set NAME=PREFIX
  qrcode-generator config list

Set a default preset with: qrcode-generator config default NAME
aliases: qr
`

const clearScreen = "\033[H\033[2J\033[3J"

// Run dispatches CLI arguments. Returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)

		return 1
	}

	if args[0] == "config" {
		return runConfig(args[1:], stdout, stderr)
	}

	return runGenerate(args, stdout, stderr)
}

func runGenerate(args []string, stdout, stderr io.Writer) int {
	shouldClear := false
	noPreset := false
	preset := ""
	contents := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		switch arg := args[i]; arg {
		case "-c", "--clear":
			shouldClear = true
		case "--no-preset":
			noPreset = true
		case "-p", "--preset":
			if i+1 >= len(args) {
				fmt.Fprintf(stderr, "error: %s requires a preset name\n", arg)

				return 1
			}

			i++
			preset = args[i]
		default:
			contents = append(contents, arg)
		}
	}

	if len(contents) == 0 {
		fmt.Fprint(stderr, usage)

		return 1
	}

	presets, err := config.Load()

	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)

		return 1
	}

	prefix := ""

	if !noPreset {
		resolved, ok := config.Resolve(presets, preset)

		if preset != "" && !ok {
			fmt.Fprintf(stderr, "error: unknown preset %q\n", preset)

			return 1
		}

		prefix = resolved
	}

	if shouldClear {
		fmt.Fprint(stdout, clearScreen)
	}

	for _, content := range contents {
		fmt.Fprintf(stdout, "--- %s ---\n", content)
		qr.Render(stdout, prefix+content)
		fmt.Fprintln(stdout)
	}

	return 0
}

func runConfig(args []string, stdout, stderr io.Writer) int {
	switch {
	case len(args) == 2 && args[0] == "set":
		return configSet(args[1], stderr)
	case len(args) == 2 && args[0] == "default":
		return configErr(stderr, config.Set(config.DefaultKey, args[1]))
	case len(args) == 1 && args[0] == "list":
		return configList(stdout, stderr)
	default:
		fmt.Fprint(stderr, usage)

		return 1
	}
}

func configSet(nameValue string, stderr io.Writer) int {
	equalIndex := strings.IndexByte(nameValue, '=')

	if equalIndex <= 0 {
		fmt.Fprintf(stderr, "error: expected NAME=PREFIX, got %q\n", nameValue)

		return 1
	}

	return configErr(stderr, config.Set(nameValue[:equalIndex], nameValue[equalIndex+1:]))
}

func configList(stdout, stderr io.Writer) int {
	presets, err := config.Load()

	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)

		return 1
	}

	names := make([]string, 0, len(presets))

	for name := range presets {
		if name != config.DefaultKey {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	defaultName := config.DefaultName(presets)

	for _, name := range names {
		marker := ""

		if name == defaultName {
			marker = "  (default)"
		}

		fmt.Fprintf(stdout, "%s=%s%s\n", name, presets[name], marker)
	}

	return 0
}

func configErr(stderr io.Writer, err error) int {
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)

		return 1
	}

	return 0
}
