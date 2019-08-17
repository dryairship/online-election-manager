const mongoose = require('mongoose');

var CandidateSchema = new mongoose.Schema({
    roll : String,
    name : String,
    email : String,
    username : String,
    password : String,
    authcode : String,
    postid : String,
    manifesto : String,
    publickey : String,
    privatekey : String,
    keystate : Number,
});

var Candidate = mongoose.model('Candidate', CandidateSchema);

function getCandidates() {
    return new Promise((resolve, reject) => {
        Candidate.find(
            null,
            (err, candidates) => {
                if(err || !candidates)
                    reject(err);
                else
                    resolve(candidates);
            }
        );
    });
}

module.exports = {getCandidates}