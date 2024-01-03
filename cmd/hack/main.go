package hack

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	// cfg := conf
	// var cfg *config = &config{
	// 	args:    allowedFlags,
	// 	flagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	// }
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	lineNum := 0
	for scanner.Scan() {
		input := scanner.Text()
		buf.WriteString(input)
		lineNum++
	}

	if len(buf.String()) > 0 {
		tempFile, err := os.CreateTemp(".", "stdFile")
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			os.Exit(1)
		}
		defer os.Remove(tempFile.Name())
		_, err = tempFile.Write(buf.Bytes())
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			os.Exit(1)
		}
		// cfg.cmd = stdIn
		// cfg.fileName = tempFile.Name()
	}
}
