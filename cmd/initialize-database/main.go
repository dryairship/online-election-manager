package main

import (
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/db"
    "io/ioutil"
    "strings"
    "fmt"
    "os"
)

func main() {
    electionDb, err := db.ConnectToDatabase()
    if err != nil {
        fmt.Println("[ERROR] Could not establish database connection.")
        os.Exit(1)
    }
    
    electionDb.ResetDatabase()
    
    file, err := ioutil.ReadFile(config.ElectionDataFilePath)
    if err != nil {
        fmt.Println("[ERROR] Election Data file not found.")
        os.Exit(1)
    }
    
    allData := strings.TrimSpace(string(file))
    posts := strings.Split(allData, "\n")
    
    for i, post := range posts {
        postID                  := fmt.Sprintf("%d", i+1)
        postData                := strings.Split(post, ",")
        postName                := postData[0]
        postRegex               := postData[1]
        candidatesRollNumbers   := make([]string, len(postData)/2-1)
        candidatesUsernames     := make([]string, len(postData)/2-1)
        manifestoLinks          := make([]string, len(postData)/2-1)
        
        for j, data := range postData {
            if j==0 || j%2==1 {
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
        
        err = electionDb.AddNewPost(&fullPost)
        if err != nil {
            fmt.Printf("[ERROR] Cannot insert post %s\n", postName)
            electionDb.ResetDatabase()
            os.Exit(1)
        }else{
            fmt.Printf("Succesfully added the post of %s.\n", postName)
        }
    }
}

func CreateCandidate(dB *db.ElectionDatabase, pID, roll, manifesto string) (models.Candidate, error) {
    skeleton, err := dB.FindStudentSkeleton(roll)
    if err != nil {
        return models.Candidate{}, err
    }
    
    candidate := models.Candidate{
        Roll:       roll,
        Name:       skeleton.Name,
        Email:      skeleton.Email+config.MailSuffix,
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
