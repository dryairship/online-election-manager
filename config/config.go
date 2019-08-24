package config

import (
	"os"
	"strconv"
)

// Possible Election States
const (
	VotingNotYetStarted = iota
	AcceptingVotes
	VotingStopped
)

// Possible states for a candidate
const (
	KeysNotGenerated = iota
	KeysGenerated
	KeysDeclared
)

// Global variables used by the program
var (
	ElectionState        int
	MailSenderUsername   string
	MailSenderPassword   string
	MailSubject          string
	MailSMTPHost         string
	MailSMTPPort         string
	MailSuffix           string
	MongoDialURL         string
	MongoDbName          string
	MongoUsername        string
	MongoPassword        string
	AssetsPath           string
	ElectionDataFilePath string
	ApplicationPort      string
	SessionsKey          string
	MaxTimeDelay         int
	RollNumberOfCEO      string
	PublicKeyOfCEO       string
	PrivateKeyOfCEO      string
	ApplicationStage     string
	ResultProgress       float64
)

// Method to read the values of the global variables from environment variables.
func InitializeConfiguration() {
	switch os.Getenv("OEMElectionState") {
	case "VotingNotYetStarted":
		ElectionState = VotingNotYetStarted
	case "AcceptingVotes":
		ElectionState = AcceptingVotes
	case "VotingStopped":
		ElectionState = VotingStopped
	default:
		panic("OEMElectionState should be one of {VotingNotYetStarted, AcceptingVotes, VotingStopped}")
	}

	MailSenderUsername = os.Getenv("OEMMailSenderUsername")
	MailSenderPassword = os.Getenv("OEMMailSenderPassword")
	MailSubject = os.Getenv("OEMMailSubject")
	MailSMTPHost = os.Getenv("OEMMailSMTPHost")
	MailSMTPPort = os.Getenv("OEMMailSMTPPort")
	MailSuffix = os.Getenv("OEMMailSuffix")

	MongoDialURL = os.Getenv("OEMMongoDialURL")
	MongoDbName = os.Getenv("OEMMongoDbName")
	MongoUsername = os.Getenv("OEMMongoUsername")
	MongoPassword = os.Getenv("OEMMongoPassword")

	AssetsPath = os.Getenv("OEMAssetsPath")
	ElectionDataFilePath = os.Getenv("OEMElectionDataFilePath")
	ApplicationPort = os.Getenv("OEMApplicationPort")
	SessionsKey = os.Getenv("OEMSessionsKey")

	MaxTimeDelay, _ = strconv.Atoi(os.Getenv("OEMMaxTimeDelay"))

	RollNumberOfCEO = os.Getenv("OEMRollNumberOfCEO")

	ApplicationStage = os.Getenv("OEMApplicationStage")
}
