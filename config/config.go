package config

import (
	"log"

	"github.com/spf13/viper"
)

// Possible Election States
const (
	VotingNotYetStarted = iota
	AcceptingVotes
	VotingStopped
	ResultsAvailable
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
	ElectionName         string
	MailSenderAuthID     string
	MailSenderEmailID    string
	MailSenderPassword   string
	MailSignature        string
	MailSMTPHost         string
	MailSMTPPort         string
	MailSuffix           string
	MongoDialURL         string
	MongoDbName          string
	MongoUsername        string
	MongoPassword        string
	AssetsPath           string
	BallotIDsPath        string
	ElectionDataFilePath string
	ApplicationPort      string
	SessionsKey          string
	MaxTimeDelay         int
	RollNumberOfCEO      string
	PublicKeyOfCEO       string
	PrivateKeyOfCEO      string
	ResultProgress       float64
)

// Method to read the values of the global variables from environment variables.
func InitializeConfiguration() {
	viper.SetConfigName("config-online-election-manager")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/go")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[WARN] Unable to locate configuration file: ", err.Error())
	}

	viper.AutomaticEnv()

	switch viper.GetString("ElectionState") {
	case "VotingNotYetStarted":
		ElectionState = VotingNotYetStarted
	case "AcceptingVotes":
		ElectionState = AcceptingVotes
	case "VotingStopped":
		ElectionState = VotingStopped
	case "ResultsCalculated":
		ElectionState = ResultsAvailable
	default:
		log.Fatal("ElectionState should be one of {VotingNotYetStarted, AcceptingVotes, VotingStopped, ResultsCalculated}")
	}

	MailSenderEmailID = viper.GetString("MailSenderEmailID")
	MailSenderAuthID = viper.GetString("MailSenderAuthID")
	MailSenderPassword = viper.GetString("MailSenderPassword")
	MailSignature = viper.GetString("MailSignature")
	MailSMTPHost = viper.GetString("MailSMTPHost")
	MailSMTPPort = viper.GetString("MailSMTPPort")
	MailSuffix = viper.GetString("MailSuffix")

	MongoDialURL = viper.GetString("MongoDialURL")
	MongoDbName = viper.GetString("MongoDbName")
	MongoUsername = viper.GetString("MongoUsername")
	MongoPassword = viper.GetString("MongoPassword")

	AssetsPath = viper.GetString("AssetsPath")
	BallotIDsPath = viper.GetString("BallotIDsPath")
	ElectionDataFilePath = viper.GetString("ElectionDataFilePath")

	ApplicationPort = viper.GetString("ApplicationPort")
	SessionsKey = viper.GetString("SessionsKey")
	MaxTimeDelay = viper.GetInt("MaxTimeDelay")

	RollNumberOfCEO = viper.GetString("RollNumberOfCEO")
	ElectionName = viper.GetString("ElectionName")
}
