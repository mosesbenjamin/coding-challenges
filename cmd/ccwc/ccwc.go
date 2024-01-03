package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	cFlag     = "-c"
	lFlag     = "-l"
	mFlag     = "-m"
	wFlag     = "-w"
	hFlag     = "-h"
	hFlagLong = "--help"
	name      = "ccwc"
)

var allowedFlags map[string]string

type config struct {
	args     map[string]string
	cmd      string
	fileName string
	flagSet  *flag.FlagSet
}

var (
	ErrInvalidPosArgSpecified = errors.New("more than one positional argument specified")
	ErrInvalidFlag            = errors.New("invalid flag specified")
	ErrInvalidCommand         = errors.New("invalid command specified")
	ErrEmptyArg               = errors.New("empty argument specified")
)

func (cfg *config) printUsage(w io.Writer) {
	cfg.flagSet.Usage()
}

func (cfg *config) runCmd(r io.Reader, w io.Writer) error {
	var err error

	switch cfg.cmd {
	case hFlag:
		cfg.printUsage(w)
	case hFlagLong:
		cfg.printUsage(w)
	case cFlag:
		numBytes, err := cfg.handleCFlag(w)
		if err != nil {
			return fmt.Errorf("flag %s in run cmd: %w", cFlag, err)
		}
		fmt.Fprintf(w, "%d %s\n", numBytes, cfg.fileName)
	case mFlag:
		numChars, err := cfg.handleMCmd(w)
		if err != nil {
			return fmt.Errorf("flag %s in run cmd: %w", mFlag, err)
		}
		fmt.Fprintf(w, "%d %s\n", numChars, cfg.fileName)
	case lFlag:
		numLines, err := cfg.handleLFlag(w)
		if err != nil {
			return fmt.Errorf("flag %s in run cmd: %w", lFlag, err)
		}
		fmt.Fprintf(w, "%d %s\n", numLines, cfg.fileName)
	case wFlag:
		wordLength, err := cfg.handleWCmd(w)
		if err != nil {
			return fmt.Errorf("flag %s in run cmd: %w", wFlag, err)
		}
		fmt.Fprintf(w, "%d %s\n", wordLength, cfg.fileName)
	default:
		numBytes, numLines, wordLength, err := cfg.handleBaseCase(w)
		if err != nil {
			return fmt.Errorf("base case in run cmd: %w", err)
		}
		fmt.Fprintf(w, "%d %d %d %s\n", numBytes, numLines, wordLength, cfg.fileName)

	}

	return err

}

func (cfg *config) handleBaseCase(w io.Writer) (int64, int64, int64, error) {
	var numBytes int64
	var numLines int64
	var wordLength int64

	numBytes, err := cfg.handleCFlag(w)
	if err != nil {
		return numBytes, numLines, wordLength, fmt.Errorf("base case: %w", err)
	}

	numLines, err = cfg.handleLFlag(w)
	if err != nil {
		return numBytes, numLines, wordLength, fmt.Errorf("base case: %w", err)
	}

	wordLength, err = cfg.handleWCmd(w)
	if err != nil {
		return numBytes, numLines, wordLength, fmt.Errorf("base case: %w", err)
	}

	return numBytes, numLines, wordLength, nil
}

func (cfg *config) handleCFlag(w io.Writer) (int64, error) {
	b, err := os.ReadFile(cfg.fileName)
	if err != nil {
		return 0, fmt.Errorf("handle c flag %w", err)
	}
	return int64(len(b)), nil
}

func (cfg *config) handleLFlag(w io.Writer) (int64, error) {
	var lineNum int64
	f, err := os.Open(cfg.fileName)
	if err != nil {
		return lineNum, fmt.Errorf("handle l flag: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
	}
	return lineNum, nil
}

func (cfg *config) handleWCmd(w io.Writer) (int64, error) {
	var wordLength int
	f, err := os.Open(cfg.fileName)
	if err != nil {
		return int64(wordLength), fmt.Errorf("handle w flag: %w", err)
	}
	defer f.Close()
	lineNum := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		words := strings.Fields(line)
		wordLength += len(words)
		lineNum++
	}
	return int64(wordLength), nil
}

func (cfg *config) handleMCmd(w io.Writer) (int64, error) {
	b, err := os.ReadFile(cfg.fileName)
	if err != nil {
		return 0, fmt.Errorf("handle c flag %w", err)
	}
	return int64(utf8.RuneCountInString(string(b))), nil
}

func (cfg *config) validateArgs(w io.Writer) (*config, error) {
	fs := cfg.flagSet

	if fs.NArg() == 1 {
		cfg.printUsage(w)
		return nil, ErrEmptyArg
	}

	if fs.NArg() == 2 {
		for k := range allowedFlags {
			if fs.Arg(1) == k && fs.Arg(1) != hFlag && fs.Arg(1) != hFlagLong {
				return nil, ErrInvalidCommand
			}
		}
		cfg.fileName = fs.Arg(1)
		cfg.cmd = fs.Arg(1)
		return cfg, nil
	}
	if fs.NArg() == 3 {
		for k := range allowedFlags {
			if fs.Arg(1) == k {
				cfg.cmd = fs.Arg(1)
				cfg.fileName = fs.Arg(2)
				return cfg, nil
			}
		}
		return nil, ErrInvalidCommand
	}
	if fs.NArg() > 3 {
		return nil, ErrInvalidPosArgSpecified
	}

	return cfg, nil
}

func (cfg *config) parseArgs(w io.Writer, args []string) (*config, error) {
	cfg.flagSet.SetOutput(w)
	for key, val := range cfg.args {
		if strings.HasPrefix(key, "--") {
			cfg.flagSet.StringVar(&key, strings.TrimPrefix(key, "--"), "", val)
		} else {
			cfg.flagSet.StringVar(&key, strings.TrimPrefix(key, "-"), "", val)
		}

	}
	cfg.flagSet.Usage = func() {
		var usageString = `
	My own version of the Unix command line tool wc!

	Usage of %s: ccwc [file] or ccwc <options> [file]`
		fmt.Fprintf(w, usageString, cfg.flagSet.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		cfg.flagSet.PrintDefaults()

	}
	err := cfg.flagSet.Parse(args)

	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func init() {
	allowedFlags = map[string]string{
		cFlag:     "outputs the number of bytes in a file",
		lFlag:     "outputs the number of lines in a file",
		wFlag:     "outputs the number of words in a file",
		mFlag:     "outputs the number of characters in a file",
		hFlag:     "display help on using this tool",
		hFlagLong: "display help on using this tool",
	}
}

func main() {
	cfg := &config{
		args:    allowedFlags,
		flagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}

	cfg, err := cfg.parseArgs(os.Stderr, os.Args[:])
	if err != nil {
		if errors.Is(err, ErrInvalidPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}

	cfg, err = cfg.validateArgs(os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	err = cfg.runCmd(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
