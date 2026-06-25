package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kristyancarvalho/tux-letter/internal/app"
	"github.com/kristyancarvalho/tux-letter/internal/config"
	"github.com/kristyancarvalho/tux-letter/internal/version"
)

const usage = `tux-letter - AI-assisted Linux and open-source newsletter

Usage:
  tux-letter <command> [flags]

Commands:
  once              Run one collection/summarization/delivery job and exit
  serve             Run as a scheduled background service
  validate-config   Load and validate the configuration file
  sources test      Fetch sources and show discovery results (no AI, no email)

Flags:
  --config <path>   Path to a TOML or JSON config file
  --version         Print version metadata
  -h, --help        Show this help
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("tux-letter", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	configPath := fs.String("config", "", "path to config file")
	showVersion := fs.Bool("version", false, "print version metadata")
	fs.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	command, rest := splitCommand(args)

	if err := fs.Parse(rest); err != nil {
		return 2
	}

	if *showVersion || command == "version" {
		fmt.Println(version.String())
		return 0
	}

	switch command {
	case "", "help", "-h", "--help":
		fmt.Print(usage)
		return 0
	case "once":
		return runWithConfig(*configPath, func(a *app.App, ctx context.Context) error {
			return a.Once(ctx)
		})
	case "serve":
		return runWithConfig(*configPath, func(a *app.App, ctx context.Context) error {
			return a.Serve(ctx)
		})
	case "validate-config":
		return runValidate(*configPath)
	case "sources test":
		return runWithConfig(*configPath, func(a *app.App, ctx context.Context) error {
			return a.SourcesTest(ctx)
		})
	case "sources":
		fmt.Fprintln(os.Stderr, "usage: tux-letter sources test")
		return 2
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n%s", command, usage)
		return 2
	}
}

func splitCommand(args []string) (string, []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "sources" && i+1 < len(args) && args[i+1] == "test" {
			rest := append([]string{}, args[:i]...)
			rest = append(rest, args[i+2:]...)
			return "sources test", rest
		}
		if takesValue(a) {
			i++
			continue
		}
		if len(a) > 0 && a[0] != '-' {
			rest := append([]string{}, args[:i]...)
			rest = append(rest, args[i+1:]...)
			return a, rest
		}
	}
	return "", args
}

func takesValue(arg string) bool {
	return arg == "--config" || arg == "-config"
}

func loadConfig(path string) (config.Config, error) {
	cfg, used, err := config.Load(path)
	if err != nil {
		return cfg, err
	}
	if used != "" {
		fmt.Fprintf(os.Stderr, "using config: %s\n", used)
	}
	return cfg, nil
}

func runValidate(path string) int {
	cfg, err := loadConfig(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println("configuration is valid")
	return 0
}

func runWithConfig(path string, fn func(*app.App, context.Context) error) int {
	cfg, err := loadConfig(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.New(cfg)
	if err := fn(a, ctx); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	return 0
}
