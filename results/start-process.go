package results

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
)

var ElectionDb *db.ElectionDatabase
var finalTally = false

func CalculateResult(database *db.ElectionDatabase) {
	ElectionDb = database
	var cmd *exec.Cmd
	if config.ApplicationStage == "development" {
		cmd = exec.Command("node", "results/calculate-result.js")
	} else {
		cmd = exec.Command("./calculate-result")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go updateStatus(stdout)
	cmd.Wait()
}

func updateStatus(r io.Reader) {
	scanner := bufio.NewScanner(r)
	var inputLine string
	var splitLine []string
	var parsedVote models.ParsedVote
	var result models.Result

	for scanner.Scan() {
		inputLine = scanner.Text()

		if inputLine == "FinalTally" {
			finalTally = true
			continue
		}

		splitLine = strings.Split(inputLine, ",")
		if splitLine == nil {
			continue
		}

		if finalTally {
			result = models.Result{
				PostID:      splitLine[0],
				Candidate:   splitLine[1],
				Preference1: splitLine[2],
				Preference2: splitLine[3],
				Preference3: splitLine[4],
			}
			ElectionDb.InsertResult(&result)
		} else {
			config.ResultProgress, _ = strconv.ParseFloat(splitLine[5], 64)
			parsedVote = models.ParsedVote{
				PostID:      splitLine[0],
				BallotID:    splitLine[1],
				Preference1: splitLine[2],
				Preference2: splitLine[3],
				Preference3: splitLine[4],
			}
			ElectionDb.InsertParsedVote(&parsedVote)
		}
	}
}
