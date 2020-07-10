#!/bin/bash

# Email ID through which mails are sent
export OEMMailSenderEmailID="chief.election.officer.iitk@gmail.com"

# Only if using IITK mailer
# Username of the account (eg: 'darshi', if email ID is darshi@iitk.ac.in)
export OEMMailSenderUsername="darshi"

# Password of the account through which mails are sent
export OEMMailSenderPassword="1L0v3CatsMeowwwww"

# Subject of the verification mail
export OEMMailSubject="Gymkhana Elections Voter Registration Verification Code"

# SMTP details
export OEMMailSMTPHost="smtp.gmail.com"
export OEMMailSMTPPort="465"

# Suffix of the email IDs for the usernames contained in `students` collection in the database
# Leave empty if `email` value itself contains the full email address
export OEMMailSuffix="@iitk.ac.in"

# MongoDB URL. Leave Default if you did not change it while starting MongoDB
export OEMMongoDialURL="mongodb://0.0.0.0:27017"

# Name of the database in which the collections are stored
export OEMMongoDbName="ElectionDb"

# Username of the Database Admin
export OEMMongoUsername="ElectionAdmin"

# Password of the Database Admin
export OEMMongoPassword="1Al50L0veCatsMeowwwwww"

# Location of the build folder
# Check the Frontend section on README for more details on this.
export OEMAssetsPath="/path/to/the/application/build/folder/of/the/frontend"

# Location of the folder where ballot Ids are to be stored 
export OEMBallotIDsPath="/path/to/the/ballotids/folder"

# CSV file that contains the details of posts and candidates
export OEMElectionDataFilePath="/path/to/datafile/electionData.yml"

# Key for encrypting data that is stored in cookies
export OEMSessionsKey="An0th3rR@nD0Mpa55w0Rd"

# Maximum time delay (in seconds) between receiving a vote and inserting it into the database.
# An optimum value would be the time difference between 
# the time at which voting stops and the time at which the results are declared.
export OEMMaxTimeDelay=10

# Roll number of the CEO, as present in the `students` collection
export OEMRollNumberOfCEO="180561"

# Port on which the application should run
export OEMApplicationPort=":80"
