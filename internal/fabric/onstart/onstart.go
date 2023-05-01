package onstart

import (
	"fmt"
	"io"
	"os"

	"goAdvancedTpl/internal/fabric/logs"
)

func WriteMessage(buildVersion string, buildDate string, buildCommit string) {

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	message := fmt.Sprintf("Build version: %s\n", buildVersion)
	printMessage(message)

	if buildDate == "" {
		buildDate = "N/A"
	}
	message = fmt.Sprintf("Build date: %s\n", buildDate)
	printMessage(message)

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	message = fmt.Sprintf("Build commit: %s\n", buildCommit)
	printMessage(message)

}

func printMessage(message string) {
	if _, err := io.WriteString(os.Stdout, message); err != nil {
		logs.New().Println(err.Error())
	}
}
