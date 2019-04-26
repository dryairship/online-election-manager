
# Online Election Manager

A cryptographically secure portal to manage online elections, created as a part of the Project Track for the course ESC101A. The algorithm is explained in the the [report](https://github.com/dryairship/online-election-manager/blob/master/report/report.pdf). The presentation slides are available [here](https://docs.google.com/presentation/d/1NUxbyJOmdoJwWQZVxAnTU4Inca59MdIHwtp_aervsp0).

## Getting Started
Follow these instructions to get a copy of the application on your computer :

### Prerequisites
- **Go**
Install Go (go1.11+) by following the instructions given on the [Go installation page](https://golang.org/doc/install).
- **MongoDB**
Install MongoDB by following the instructions given on the [MongoDB installation page](https://docs.mongodb.com/manual/installation/).
- **sjcl**
	- You can either compile it yourself from the [source code](https://github.com/bitwiseshiftleft/sjcl) by following the instructions given [here](https://github.com/bitwiseshiftleft/sjcl/wiki/Getting-Started). [*Note that if you are compiling yourself, you need to do `./configure --with-ecc` before executing the `make` command.*]
	- Or you can download a pre-compiled version (with the ecc feature) from [here](https://ufile.io/8zvrzwdx).
	- You will need the `sjcl.js` file later.

### Installing
``` 
cd $GOPATH/src
go get github.com/dryairship/online-election-manager/cmd/online-election-manager
go get github.com/dryairship/online-election-manager/cmd/initialize-database
cd github.com/dryairship/online-election-manager
go install ./...
```
This creates two executables in `$GOPATH/bin` :
- `online-election-manager`
- `initialize-database`

Also, copy the `sjcl.js` file into `assets/js/`. 
## Configuration
You need to create a System Admin and populate the database with details of the students. You also need to create a configuration file which will export important values to the environment variables.
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
- Fill this collection with the list of students eligible to vote. Each student is a separate document with three fields (all are strings) :
	- `roll` - Roll number of the student.
	- `email` - Just the username of the email address of the student (e.g. this field contains `darshi` if the email ID of the student is `darshi@iitk.ac.in`. *(The remaining part of the email ID will be specified in the configuration file)*.
	- `name` - Full name of the student.
### Creating the election data CSV file
Follow the template given in `configurationTemplates/electionData.csv` to create a CSV file containing the details of the posts for which the elections are being held, the voters eligible to vote for each post, and the candidates contesting the elections for various posts.

### Creating the configuration file
Follow the template given in `configurationTemplates/configuration.sh` and replace the default values with your own values.
- You can have any number of posts. All the details of one post are in a single row.
- The first value is the name of the post.
- The second value is the regular expression that the students' roll number must match with, to check if they are eligible to vote for this post.
- The remaining values in each row are the details of the candidates. You need to fill two values for each candidate : 
	- The first value is the roll number of the candidate.
	- The second value is the link to the manifesto of the candidate.
All these values occur alternately on the same row. Thus, the format for each row is:
```
<Name Of The Post>,<Regular expression for voters>,<Roll no. of candidate 1>,<Manifesto of candidate 1>,<Roll no. of candidate 2>,<Manifesto of candidate 2>,<Roll no. of candidate 3>,<Manifesto of candidate 3>,...
```
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

