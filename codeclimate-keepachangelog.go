package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/codeclimate/cc-engine-go/engine"
)

type validationError struct {
	Line   int
	Column int
}

var validationErrRgx = regexp.MustCompile(`at line (\d+), column (\d+)`)

func newValidationError(s string) *validationError {
	errSubmatch := validationErrRgx.FindStringSubmatch(s)

	if len(errSubmatch) == 0 {
		return nil
	}

	line, _ := strconv.Atoi(errSubmatch[1])
	column, _ := strconv.Atoi(errSubmatch[2])

	return &validationError{
		Line:   line,
		Column: column,
	}
}

func prefixInArr(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func getAnalysisFiles(rootPath string, config engine.Config) ([]string, error) {
	var analysisFiles []string

	err := filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, "CHANGELOG.md") && prefixInArr(path, engine.IncludePaths(rootPath, config)) {
			analysisFiles = append(analysisFiles, path)
			return nil
		}
		return err
	})

	return analysisFiles, err
}

func main() {
	rootPath := "/code/"

	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	analysisFiles, err := getAnalysisFiles(rootPath, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %v\n", err)
		os.Exit(1)
	}

	for _, path := range analysisFiles {
		output, err := execParse(path)

		vErr := newValidationError(output)
		fatalErr := err != nil && vErr == nil // not merely validation err

		relativePath := strings.SplitAfter(path, rootPath)[1]

		// fatal error, bail out
		if fatalErr {
			exitOnRunError(relativePath, output, err)
			return
		}

		// validation err, report it but keep going
		if vErr != nil  {
			registerIssue(relativePath, output, *vErr)
			continue
		}
	}
}

func execParse(path string) (string, error) {
	cmd := exec.Command("parse", path)

	outputBytes, err := cmd.CombinedOutput()
	output := strings.TrimPrefix(string(outputBytes), "ERROR: ")

	return output, err
}

func exitOnRunError(path string, msg string, err error) {
	fmt.Fprintf(os.Stderr, "Error analyzing path: %v\n", path)
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	fmt.Fprintf(os.Stderr, "parse_a_changelog output: %v\n", msg)

	os.Exit(1)
}

func registerIssue(path string, msg string, vErr validationError) {
	description := fmt.Sprintf(
		"Your changelog does not pass validation: %s",
		strings.TrimSuffix(msg, "\n"),
	)

	issue := &engine.Issue{
		Type:              "issue",
		Check:             "Changelog/Style/Changelog",
		Description:       description,
		RemediationPoints: int32(50000),
		Categories:        []string{"Style"},
		Location: &engine.Location{
			Path: path,
			Positions: &engine.LineColumnPosition{
				Begin: &engine.LineColumn{
					Line:   vErr.Line,
					Column: vErr.Column,
				},
				End: &engine.LineColumn{
					Line:   vErr.Line,
					Column: vErr.Column,
				},
			},
		},
	}
	engine.PrintIssue(issue)
}
