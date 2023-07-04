package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"

	"golang.org/x/term"

	"github.com/hashicorp/vault/shamir"
)

var (
	partsFlag     int    = 3
	fileFlag      string = "-"
	outputDirFlag string = ""
	helpMode      bool
	thresholdFlag int = 2
)

func init() {
	flag.StringVar(&fileFlag, "f", "-", "read the secret from file instead of STDIN")
	flag.IntVar(&partsFlag, "parts", 3, "number of shares")
	flag.BoolVar(&helpMode, "help", false, "display this help and exit.")
	flag.StringVar(&outputDirFlag, "output", "CURDIR", "output directory (must not exist)")
	flag.IntVar(&thresholdFlag, "threshold", 2, "minimum number of shares required to reconstruct the secret")
	flag.Usage = usage
	flag.ErrHelp = nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("keyshtool: ")
	log.SetOutput(os.Stderr)
	flag.Parse()

	if helpMode {
		usage()
		return
	}

	if flag.NArg() < 1 {
		log.Fatal("invalid command")
	}

	switch cmd := flag.Arg(0); cmd {
	case "split":
		if err := split(); err != nil {
			log.Fatal(err)
		}

	case "combine":
		secret, err := combine()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("secret:", string(secret))

	default:
		log.Fatal("invalid command")
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: keyshtool [OPTION]... COMMAND")
	fmt.Fprintln(os.Stderr, "Commands: combine split")

	flag.PrintDefaults()
}

func split() error {
	var outputDir string

	parts, err := shamir.Split(bytes.TrimSpace(mustReadSecret()), partsFlag, thresholdFlag)
	if err != nil {
		return err
	}

	if outputDirFlag == "CURDIR" {
		outputDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		outputDir = outputDirFlag
	}

	outputDir = path.Join(outputDir, "PARTS")
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		log.Fatalf("couldn't open the directory %q", outputDir)
	}

	for i, part := range parts {
		if err := os.WriteFile(path.Join(outputDir, fmt.Sprintf("%06d-part.txt", i)), part, 0600); err != nil {
			return err
		}
	}

	return nil
}

func combine() ([]byte, error) {
	var parts [][]byte
	for i := 0; i < flag.NArg()-1; i++ {
		bs, err := os.ReadFile(flag.Args()[i+1])
		if err != nil {
			log.Fatal(err)
		}

		parts = append(parts, bs)
	}

	return shamir.Combine(parts)
}

func mustReadSecret() []byte {
	if fileFlag != "-" {
		bs, err := os.ReadFile(fileFlag)
		if err != nil {
			log.Fatal(err)
		}

		return bs
	}

	fmt.Fprint(os.Stderr, "Enter secret: ")
	bs, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	return bs
}
