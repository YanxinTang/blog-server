package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	errInvalidSequenceWidth = errors.New("digits must be positive")
	errInvalidName          = errors.New("name can not be empty")
)

var name string

func init() {
	migrateCreateCmd.Flags().StringVar(&name, "name", "", "name of the new migration")
	migrateMigrateCmd.MarkFlagRequired("name")
	migrateCmd.AddCommand(migrateCreateCmd)
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create migrations",
	Long:  "create migrations",
	Run: func(cmd *cobra.Command, args []string) {
		dir = filepath.Clean(dir)
		ext := ".sql"
		seqDigits := 6

		if err := create(dir, name, ext, seqDigits, true); err != nil {
			exitWithError(err)
		}
	},
}

func create(dir string, name string, ext string, seqDigits int, print bool) error {
	if len(name) == 0 {
		return errInvalidName
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		exitWithError(err)
	}
	version, err := nextSeqVersion(matches, seqDigits)
	if err != nil {
		exitWithError(err)
	}
	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err = filepath.Glob(versionGlob)
	if err != nil {
		return err
	}
	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			return err
		}

		if print {
			absPath, _ := filepath.Abs(filename)
			fmt.Println(absPath)
		}
	}

	return nil
}

func nextSeqVersion(matches []string, seqDigits int) (string, error) {
	if seqDigits <= 0 {
		return "", errInvalidSequenceWidth
	}

	nextSeq := uint64(1)

	if len(matches) > 0 {
		filename := matches[len(matches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("next sequence number %s too large. At most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	// os.Create() specifies 0666 as the FileMode, so we're doing the same
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}
