const mongoose = require('mongoose');

var VoteSchema = new mongoose.Schema({
    postid : String,
    data : String
});

var Vote = mongoose.model('Vote', VoteSchema);

function getVotes() {
    return new Promise((resolve, reject) => {
        Vote.find(
            null,
            (err, votes) => {
                if(err || !votes)
                    reject(err);
                else
                    resolve(votes);
            }
        );
    });
}

module.exports = {getVotes}