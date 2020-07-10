
# Online Election Manager

A cryptographically secure portal to manage online elections, created as a part of the Project Track for the course ESC101A. The algorithm is explained in the [report](https://github.com/dryairship/online-election-manager/blob/master/report/report.pdf). The presentation slides are available [here](https://docs.google.com/presentation/d/1NUxbyJOmdoJwWQZVxAnTU4Inca59MdIHwtp_aervsp0).

## Getting Started
Follow these instructions to get a copy of the application on your computer :

### Prerequisites
- **Go**  
Install Go (go1.11+) by following the instructions given on the [Go installation page](https://golang.org/doc/install).
- **MongoDB**  
Install MongoDB by following the instructions given on the [MongoDB installation page](https://docs.mongodb.com/manual/installation/).

### Installing
``` 
go get github.com/dryairship/online-election-manager/cmd/online-election-manager
go get github.com/dryairship/online-election-manager/cmd/initialize-database
```
This creates two executables in `$GOPATH/bin` :
- `online-election-manager`
- `initialize-database`

## Configuration
You need to create a user in the database and populate the database with details of the students. You also need to create a configuration file which will export important values to the environment variables.

### Initializing the database
- Start MongoDB.  
`mongod`
- Connect to the database in a new terminal.  
`mongo`
- Create a new Election Database.  
`> use ElectionDb`
- Create a system admin for this database.  
`> db.createUser({user: "ElectionAdmin", pwd: "password", roles: ["readWrite", "dbAdmin"]})`
- Create a new collection called `students`.  
`> db.createCollection("students")`
- Fill this collection with the list of students (voters, candidates and the CEO). Each student is a separate document with three fields (all are strings) :
	- `roll` - Roll number of the student.
	- `email` - Just the username of the email address of the student (e.g. this field contains `darshi` if the email ID of the student is `darshi@iitk.ac.in`. *(The remaining part of the email ID will be specified in the configuration file)*.
	- `name` - Full name of the student.

### Creating the election data YML file
Follow the template given in `configurationTemplates/electionData.yml` to create a YML file containing the details of the posts for which the elections are being held, the voters eligible to vote for each post, and the candidates contesting the elections for various posts.  
- You can have any number of posts. Each post must have a unique ID.
- Each post can have any number of candidates. Each candidate for a post must have a unique roll number.
- All the fields (as specified in the sample YML file) are required.

### Creating the configuration file
Follow the template given in `configurationTemplates/configuration.sh` and replace the default values with your own values.

### Frontend
You need to build the frontend separately. The frontend is available at [dryairship/online-election-manager-frontend](https://github.com/dryairship/online-election-manager-frontend). Run `node run build` to build the frontend. Then specify the path of the build folder in the configuration file.

## Deployment
Follow these instructions to host an election using this application :
- Start MongoDB in auth mode.  
`mongod --auth`
- Open a new terminal and load values from the configuration file into the environment variables.  
`source /path/to/configuration/file/configuration.sh`
- Change to the directory where the application's executables were installed.  
`cd $GOPATH/bin`
- Initialize the database with the details of the posts and the candidates.  
`./initialize-database`
- Start the application.  
`./online-election-manager`

## Authors
- [Priydarshi Singh](https://dryairship.github.io)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/dryairship/online-election-manager/blob/master/LICENSE) file for details.
