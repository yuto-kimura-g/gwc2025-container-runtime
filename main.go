package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	switch os.Args[1] {
	case "run":
		if err := run(os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}

func run(command []string) error {
	if err := SetupNamespace(); err != nil {
		return err
	}

	if err := SetupChroot(); err != nil {
		return err
	}

	if err := SetupCgroup(); err != nil {
		return err
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
