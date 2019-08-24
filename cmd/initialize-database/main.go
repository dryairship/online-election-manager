package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
)

func main() {
	// Initialize the configuration from environment variables.
	config.InitializeConfiguration()

	// Connect to the database.
	electionDb, err := db.ConnectToDatabase()
	if err != nil {
		fmt.Println("[ERROR] Could not establish database connection.")
		os.Exit(1)
	}

	// Delete all entries from the database.
	electionDb.ResetDatabase()

	// Open the CSV file for reading the data about posts and candidates.
	file, err := ioutil.ReadFile(config.ElectionDataFilePath)
	if err != nil {
		fmt.Println("[ERROR] Election Data file not found.")
		os.Exit(1)
	}

	// Convert data into a simpler format.
	allData := strings.TrimSpace(string(file))
	posts := strings.Split(allData, "\n")

	for _, post := range posts {
		// Extract data about a post from the data read from the file.
		postData := strings.Split(post, ",")
		postID := postData[0]
		postName := postData[1] + ", " + postData[2]
		postRegex := postData[3]
		candidatesRollNumbers := make([]string, len(postData)/2-2)
		candidatesUsernames := make([]string, len(postData)/2-2)
		manifestoLinks := make([]string, len(postData)/2-2)

		for j, data := range postData {
			if j <= 3 || j%2 == 1 {
				continue
			}
			candidatesRollNumbers[j/2-1] = data
			manifestoLinks[j/2-1] = postData[j+1]
		}

		for j, cand := range candidatesRollNumbers {
			candidate, err := CreateCandidate(&electionDb, postID, cand, manifestoLinks[j])
			if err != nil {
				fmt.Printf("[ERROR] Cannot find candidate %s\n", cand)
				electionDb.ResetDatabase()
				os.Exit(1)
			}

			candidatesUsernames[j] = candidate.Username

			// Insert the newly created candidate into the database.
			err = electionDb.AddNewCandidate(&candidate)
			if err != nil {
				fmt.Printf("[ERROR] Cannot add candidate %s\n", cand)
				electionDb.ResetDatabase()
				os.Exit(1)
			}
		}

		fullPost := models.Post{
			PostID:     postID,
			PostName:   postName,
			VoterRegex: postRegex,
			Candidates: candidatesUsernames,
		}

		// Insert the newly created post into the database.
		err = electionDb.AddNewPost(&fullPost)

		if err != nil {
			fmt.Printf("[ERROR] Cannot insert post %s\n", postName)
			electionDb.ResetDatabase()
			os.Exit(1)
		} else {
			fmt.Printf("Succesfully added the post of %s.\n", postName)
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
