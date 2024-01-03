package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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
		fmt.Println("handle '-c command'")
	case mFlag:
		fmt.Println("handle '-m command'")
	case lFlag:
		fmt.Println("handle '-l command'")
	case wFlag:
		fmt.Println("handle '-w command'")
	default:
		fmt.Println("handle base case")
	}

	return err

}

func (cfg *config) validateArgs(w io.Writer) (*config, error) {
	fs := cfg.flagSet

	if fs.NArg() == 1 {
		cfg.printUsage(w)
		return nil, ErrEmptyArg
	}

	if fs.NArg() == 2 {
		for k, _ := range allowedFlags {
			if fs.Arg(1) == k && fs.Arg(1) != hFlag && fs.Arg(1) != hFlagLong {
				return cfg, ErrInvalidCommand
			}
		}
		cfg.fileName = fs.Arg(1)
		cfg.cmd = fs.Arg(1)
		return cfg, nil
	}
	if fs.NArg() == 3 {
		for k, _ := range allowedFlags {
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
