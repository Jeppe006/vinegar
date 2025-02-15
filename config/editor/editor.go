package editor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vinegarhq/vinegar/config"
)

// Editor retrieves the editor from the environment.
// If no environment variable is present it will fall back
// to nano, returning an error if it doesn't exist.
func Editor() (string, error) {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	log.Println("no EDITOR set, falling back to nano")

	return exec.LookPath("nano")
}

// Edit loops over editing the named configuration file name with
// an editor retrieved from [Editor] until it has no errors.
func Edit(name string) error {
	editor, err := Editor()
	if err != nil {
		return fmt.Errorf("editor: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}

	if err := fillTemplate(name); err != nil {
		return err
	}

	for {
		cmd := exec.Command(editor, name)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			return err
		}

		if _, err := config.Load(name); err != nil {
			log.Println(err)
			log.Println("Press enter to re-edit configuration file")
			fmt.Scanln()

			continue
		}

		break
	}

	return nil
}

func EditNonToml(name string) error {
	editor, err := Editor()
	if err != nil {
		return fmt.Errorf("editor: %w", err)
	}

	cmd := exec.Command(editor, name)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

func fillTemplate(name string) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.Size() > 1 {
		return nil
	}

	template := "# See how to configure Vinegar on the documentation website:\n" +
		"# https://vinegarhq.org/Configuration\n\n"

	log.Println("Writing Configuration template")

	if _, err := f.WriteString(template); err != nil {
		return err
	}

	return nil
}
