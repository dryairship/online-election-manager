package results

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
)

var ElectionDb *db.ElectionDatabase

func CalculateResult(database *db.ElectionDatabase) {
	ElectionDb = database

	err := ElectionDb.ClearResults()
	if err != nil {
		panic(err)
	}

	votedVoters, err := ElectionDb.FindVotedVoters()
	if err != nil {
		panic(err)
	}

	votedVotersFile, err := os.Create(config.AssetsPath + "/votedVoters.csv")
	if err != nil {
		panic(err)
	}
	defer votedVotersFile.Close()

	votedVotersFile.WriteString("Roll Number,Name,\n")
	for _, voter := range votedVoters {
		votedVotersFile.WriteString(voter.Roll + "," + voter.Name + ",\n")
	}

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

	votesData, err := os.Create(config.AssetsPath + "/votesData.csv")
	if err != nil {
		panic(err)
	}
	defer votesData.Close()

	votesData.WriteString("Post ID,Ballot ID,Preference 1,Preference 2,Preference 3\n")

	var inputLine string
	var splitLine []string
	var parsedVote models.ParsedVote
	var result models.Result

	for scanner.Scan() {
		inputLine = scanner.Text()

		splitLine = strings.Split(inputLine, ",")
		if splitLine == nil {
			continue
		}

		if splitLine[0] == "FT" {
			result = models.Result{
				PostID:      splitLine[1],
				Candidate:   splitLine[2],
				Preference1: splitLine[3],
				Preference2: splitLine[4],
				Preference3: splitLine[5],
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
			votesData.WriteString(strings.Join(splitLine[:5], ",") + "\n")
			ElectionDb.InsertParsedVote(&parsedVote)
		}
	}
}
