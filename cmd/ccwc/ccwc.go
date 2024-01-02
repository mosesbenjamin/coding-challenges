package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	cCmd     = "-c"
	lCmd     = "-l"
	mCmd     = "-m"
	wCmd     = "-w"
	hCmd     = "-h"
	hCmdLong = "--help"
	appName  = "ccwc"
)

type config struct {
	arg      string
	fileName string
	flagSet  *flag.FlagSet
}

var (
	ErrInvalidPosArgSpecified = errors.New("more than one positional argument specified")
	ErrInvalidFlag            = errors.New("invalid flag specified")
)

func (cfg *config) printUsage(w io.Writer) {
	cfg.flagSet.Usage()
}

func (cfg *config) runCmd(r io.Reader, w io.Writer) error {
	var err error
	switch cfg.arg {
	case hCmd:
		cfg.printUsage(w)
	case hCmdLong:
		cfg.printUsage(w)
	case cCmd:
		fmt.Println("handle '-c command'")
	}
	return err
}

func (cfg *config) validateArgs() (*config, error) {
	fs := cfg.flagSet

	if fs.NArg() == 2 {
		fs := cfg.flagSet
		cfg.arg = fs.Arg(1)
	}
	if fs.NArg() == 3 {
		cfg.arg = fs.Arg(1)
		cfg.fileName = fs.Arg(2)
	}
	if fs.NArg() > 3 {
		return nil, ErrInvalidPosArgSpecified
	}

	var validFlag bool
	allowedFlags := []string{cCmd, lCmd, mCmd, wCmd, hCmd, hCmdLong}
	for _, v := range allowedFlags {
		if cfg.arg == v {
			validFlag = true
		}
	}
	if !validFlag {
		return nil, ErrInvalidFlag
	}

	return cfg, nil
}

func (cfg *config) parseArgs(w io.Writer, args []string) (*config, error) {
	cfg.flagSet.SetOutput(w)
	cfg.flagSet.StringVar(&cfg.fileName, "c", "", "File to read from.")
	cfg.flagSet.Usage = func() {
		var usageString = `
	My own version of the Unix command line tool wc!

	Usage of %s: <options> [name]`
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

func main() {
	cfg := &config{
		flagSet: flag.NewFlagSet(appName, flag.ContinueOnError),
	}
	cfg, err := cfg.parseArgs(os.Stderr, os.Args[:])
	if err != nil {
		if errors.Is(err, ErrInvalidPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}

	cfg, err = cfg.validateArgs()
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
