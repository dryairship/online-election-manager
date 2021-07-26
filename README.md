
# Online Election Manager

A cryptographically secure portal to manage online elections, originally created as a part of the Project Track for the course ESC101A. The version of the project presented for the Project Track is available in the [project-track](https://github.com/dryairship/online-election-manager/tree/project-track) branch.

## Usage

### Required data

To set up the election portal, the following data is needed:
 - A mongodb database with a collection called `students`, which contains the details of the students (the voters, candidates and CEO all should be present in this table). The following fields (of type string) must be present for each record in the collection:  
   - `name`: The name of the student
   - `roll`: The roll number of the student
   - `email`: The email of the student (without the @org.tld suffix). Eg: If the email of the student is `darshi@iitk.ac.in` then this field contains `darshi`.
 - A YML file ([`resources/electionData.yml`](https://github.com/dryairship/online-election-manager/blob/master/resources/electionData.yml)) containing the details of the posts for which the elections are being held. The following data is needed for each post:  
   - `id`: An integer (should be unique for each post)
   - `name`: Name of the post
   - `hasNota`: Whether or not the post allows the voter to choose NOTA
   - `candidates`: The roll numbers of the candidates contesting for this post
 - The list of voters eligible to vote, for each post. This list should be in the form of line separated values, in a file named `<postID>.txt` and put in the [`resources/voters`](https://github.com/dryairship/online-election-manager/tree/master/resources/voters) directory.
 - Configuration values in `backend-config.yml` (note that this file needs to be updated both in the [main repo](https://github.com/dryairship/online-election-manager/blob/master/backend-config.yml) and in the [backend submodule](https://github.com/dryairship/online-election-manager-backend/blob/master/backend-config.yml)).
 - Frontend arguments values in [`docker-compose.yml`](https://github.com/dryairship/online-election-manager/blob/master/docker-compose.yml).

### Before deployment

 - Run `go run init-voters`. This will populate the `voters` and `ceo` collections.
 - Send authcodes to all the voters. (TODO: Add details)
 - Run `go run initialize-database`. This will populate the `posts` and `candidates` collections.

### Deployment

```
docker-compose up -d --build frontend
docker-compose up -d --build backend
```

## Authors
- [Priydarshi Singh](https://dryairship.github.io)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/dryairship/online-election-manager/blob/master/LICENSE) file for details.
