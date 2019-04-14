# Online Election Manager

A cryptographically secure portal to manage online elections, created as a part of the Project Track for the course ESC101A. 

This project aims to provide students the comfort of voting on their own devices, without the hassle of standing in a queue at the polling booth. Students register on the portal by providing their roll numbers. This sends an authentication code to their official IITK email addresses, which they use to choose a password. They then vote for their candidates and the votes are encrypted asymmetrically using the public keys of the candidates, and sent to the server. The votes are decrypted at the time of declaration of results using the private keys of the candidates.
