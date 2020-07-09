package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
)

func ExportBallotIdsToFile(ballotIds []models.UsedBallotID, fileName string) error {
	path := fmt.Sprintf("%s/ballotids/%s", config.AssetsPath, fileName)
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

	return nil
}
