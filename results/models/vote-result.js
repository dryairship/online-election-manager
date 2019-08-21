const mongoose = require('mongoose');

var VoteResultSchema = new mongoose.Schema({
    voteid : String,
    pref1 : String,
    pref2 : String,
    pref3 : String
});

var VoteResult = mongoose.model('VoteResult', VoteResultSchema);

function insertVoteResult(id, p1, p2, p3) {
    return new Promise((resolve, reject) => {
        voteResult = new VoteResult({
            voteid : ""+id,
            pref1 : p1,
            pref2 : p2,
            pref3 : p3,
        });
        voteResult.save(err => {
            console.log(err);
            if(err){
                reject(err);
            }else{
                resolve(voteResult);
            }
        });
    });
}

module.exports = {insertVoteResult}