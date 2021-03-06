package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aquasecurity/bench-common/check"
	"github.com/aquasecurity/bench-common/outputter"
	"github.com/aquasecurity/bench-common/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func app(cmd *cobra.Command, args []string) {
	err := checkDefinitionFilePath(cfgFile)
	if err != nil {
		util.ExitWithError(err)
	}

	Main(cfgFile, define)
}

// Main entry point for benchmark functionality
func Main(filePath string, constraints []string) {
	controls, err := getControls(filePath, constraints)
	if err != nil {
		util.ExitWithError(err)
	}

	summary := runControls(controls, "")
	err = outputResults(controls, summary)
	if err != nil {
		util.ExitWithError(err)
	}
}

func outputResults(controls *check.Controls, summary check.Summary) error {
	config := &outputter.Config{
		Console: outputter.Console{
			NoRemediations:    noRemediations,
			IncludeTestOutput: includeTestOutput,
		},
		JSONFormat: jsonFmt,
	}
	o := outputter.BuildOutputter(summary, config)

	return o.Output(controls, summary)
}

func runControls(controls *check.Controls, checkList string) check.Summary {
	var summary check.Summary

	if checkList != "" {
		ids := util.CleanIDs(checkList)
		summary = controls.RunChecks(ids...)
	} else {
		summary = controls.RunGroup()
	}

	return summary
}

func getControls(path string, constraints []string) (*check.Controls, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	controls, err := check.NewControls([]byte(data), constraints)
	if err != nil {
		return nil, err
	}

	return controls, err
}

func checkDefinitionFilePath(filePath string) (err error) {
	glog.V(2).Info(fmt.Sprintf("Looking for config file: %s\n", filePath))
	_, err = os.Stat(filePath)

	return err
}
