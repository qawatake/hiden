package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/qawatake/hiden/internal/config"
	"github.com/qawatake/hiden/internal/finder"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "ls":
		if err := runLs(); err != nil {
			if errors.Is(err, finder.ErrCancelled) {
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("hiden version %s\n", version)
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func runLs() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	path, err := finder.Run(cfg.Dirname)
	if err != nil {
		return err
	}

	if path != "" {
		fmt.Println(path)
	}
	return nil
}

func printHelp() {
	fmt.Println(`hiden - Search personal memo/script directories across ghq repositories

Usage:
  hiden <command>

Commands:
  ls       Search and select files from hiden directories
  version  Print version information
  help     Print this help message`)
}
