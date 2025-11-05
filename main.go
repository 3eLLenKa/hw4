package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	Count      bool
	Duplicates bool
	Unique     bool
	IgnoreCase bool
	SkipFields int
	SkipChars  int
}

func main() {
	opts := parseFlags()
	args := flag.Args()

	var input io.Reader
	var output io.Writer
	var err error

	switch len(args) {
	case 0:
		input = os.Stdin
		output = os.Stdout
	case 1:
		input, err = os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer input.(*os.File).Close()
		output = os.Stdout
	default:
		in, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer in.Close()
		out, err := os.Create(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer out.Close()
		input = in
		output = out
	}
	lines := readLines(input)
	result := Uniq(lines, opts)

	writeLines(output, result)
}

func parseFlags() Options {
	countFlag := flag.Bool("c", false, "подсчитать количество повторов")
	dupFlag := flag.Bool("d", false, "вывести только повторяющиеся строки")
	uniqFlag := flag.Bool("u", false, "вывести только уникальные строки")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	skipFields := flag.Int("f", 0, "пропустить первые num полей")
	skipChars := flag.Int("s", 0, "пропустить первые num символов")

	flag.Parse()

	modeCount := 0

	if *countFlag {
		modeCount++
	}
	if *dupFlag {
		modeCount++
	}
	if *uniqFlag {
		modeCount++
	}
	if modeCount > 1 {
		fmt.Fprintln(os.Stderr, "флаги -c, -d и -u нельзя использовать вместе")
		os.Exit(1)
	}

	return Options{
		Count:      *countFlag,
		Duplicates: *dupFlag,
		Unique:     *uniqFlag,
		IgnoreCase: *ignoreCase,
		SkipFields: *skipFields,
		SkipChars:  *skipChars,
	}
}

func readLines(r io.Reader) []string {
	scanner := bufio.NewScanner(r)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func writeLines(w io.Writer, lines []string) {
	writer := bufio.NewWriter(w)

	defer writer.Flush()

	for _, l := range lines {
		fmt.Fprintln(writer, l)
	}
}

func normalize(s string, opts Options) string {
	if opts.IgnoreCase {
		s = strings.ToLower(s)
	}

	fields := strings.Fields(s)
	if opts.SkipFields > 0 && opts.SkipFields < len(fields) {
		s = strings.Join(fields[opts.SkipFields:], " ")
	}
	if opts.SkipChars > 0 && opts.SkipChars < len(s) {
		s = s[opts.SkipChars:]
	}

	return s
}

func Uniq(lines []string, opts Options) []string {
	if len(lines) == 0 {
		return nil
	}

	type group struct {
		original string
		count    int
	}

	var groups []group
	prevNorm := normalize(lines[0], opts)
	current := group{original: lines[0], count: 1}

	for i := 1; i < len(lines); i++ {
		curNorm := normalize(lines[i], opts)
		if curNorm == prevNorm {
			current.count++
		} else {
			groups = append(groups, current)
			current = group{original: lines[i], count: 1}
			prevNorm = curNorm
		}
	}

	groups = append(groups, current)
	var result []string

	for _, g := range groups {
		switch {
		case opts.Duplicates && g.count > 1:
			result = append(result, g.original)
		case opts.Unique && g.count == 1:
			result = append(result, g.original)
		case opts.Count:
			result = append(result, fmt.Sprintf("%d %s", g.count, g.original))
		case !opts.Duplicates && !opts.Unique && !opts.Count:
			result = append(result, g.original)
		}
	}
	return result
}
