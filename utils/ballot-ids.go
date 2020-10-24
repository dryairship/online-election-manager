package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
)

func ExportSingleVoteBallotIdsToFile(ballotIds []models.UsedSingleVoteBallotID, fileName string) error {
	path := fmt.Sprintf("%s/%s", config.BallotIDsPath, fileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for _, ballotId := range ballotIds {
		_, err = w.WriteString(fmt.Sprintf("%s - %s - %s\n", ballotId.BallotString, ballotId.Roll, ballotId.Name))
		if err != nil {
			return err
		}
	}

	w.Flush()
	log.Println("[INFO] Ballot IDs saved in: ", fileName)

	return nil
}

func ExportBallotIdsToFile(ballotIds []string, fileName string) error {
	path := fmt.Sprintf("%s/%s", config.BallotIDsPath, fileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for _, ballotId := range ballotIds {
		_, err = w.WriteString(fmt.Sprintf("%s\n", ballotId))
		if err != nil {
			return err
		}
	}

	w.Flush()
	log.Println("[INFO] Ballot IDs saved in: ", fileName)

	return nil
}
