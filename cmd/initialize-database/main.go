package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
)

type InitCandidate struct {
	Roll      string
	Manifesto string
}

type InitPost struct {
	Id         string
	Name       string
	Regex      string
	Candidates []InitCandidate
}

type InitData struct {
	Posts []InitPost
}

func main() {
	// Initialize the configuration from environment variables.
	config.InitializeConfiguration()

	// Connect to the database.
	electionDb, err := db.ConnectToDatabase()
	if err != nil {
		log.Fatal("[ERROR] Could not establish database connection.")
	}

	// Delete all entries from the database.
	err = electionDb.ResetDatabase()
	if err != nil {
		log.Fatal("[ERROR] Could not reset database.")
	}

	// Open the CSV file for reading the data about posts and candidates.
	fileData, err := ioutil.ReadFile(config.ElectionDataFilePath)
	if err != nil {
		log.Fatal("[ERROR] Election Data file not found.")
	}

	var data InitData
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("[ERROR] Election Data is not in the correct format.")
	}

	for _, post := range data.Posts {

		var candidatesUsernames []string

		for _, candidateData := range post.Candidates {
			candidate, err := CreateCandidate(&electionDb, post.Id, candidateData.Roll, candidateData.Manifesto)
			if err != nil {
				log.Fatalf("[ERROR] Cannot find candidate %s\n", candidateData.Roll)
			}

			candidatesUsernames = append(candidatesUsernames, candidate.Username)

			// Insert the newly created candidate into the database.
			err = electionDb.AddNewCandidate(&candidate)
			if err != nil {
				log.Fatalf("[ERROR] Cannot add candidate %s\n", candidateData.Roll)
			}
		}

		fullPost := models.Post{
			PostID:     post.Id,
			PostName:   post.Name,
			VoterRegex: post.Regex,
			Candidates: candidatesUsernames,
			Resolved:   false,
		}

		// Insert the newly created post into the database.
		err = electionDb.AddNewPost(&fullPost)

		if err != nil {
			log.Fatalf("[ERROR] Cannot insert post %s\n", post.Name)
		} else {
			log.Printf("Succesfully added the post of %s.\n", post.Name)
		}
	}
}

// Function to create a candidate from the data in the file.
func CreateCandidate(dB *db.ElectionDatabase, pID, roll, manifesto string) (models.Candidate, error) {
	skeleton, err := dB.FindStudentSkeleton(roll)
	if err != nil {
		return models.Candidate{}, err
	}

	candidate := models.Candidate{
		Roll:       roll,
		Name:       skeleton.Name,
		Email:      skeleton.Email + config.MailSuffix,
		Username:   fmt.Sprintf("P%sC%s", pID, roll),
		Password:   "",
		AuthCode:   "",
		PostID:     pID,
		Manifesto:  manifesto,
		PublicKey:  "",
		PrivateKey: "",
		KeyState:   config.KeysNotGenerated,
	}

	return candidate, nil
}
