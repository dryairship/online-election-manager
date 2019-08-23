const fs = require('fs');

const DBConnection = require('./db.js');
const sjcl = require('./sjcl.js');

const {getVotes} = require('./models/vote.js');
const {getCandidates} = require('./models/candidate.js');
const {insertVoteResult} = require('./models/vote-result.js');
const {getCEO} = require('./models/ceo.js');

DBConnection.dial();

var fetchedVotes = [];
var fetchedCandidates = [];
var fetchedCEO;
var privateKeyOfCEO;
var result = [];
var scores = [];

function unserializePrivateKey(serPriv){
    return new sjcl.ecc.elGamal.secretKey(
        sjcl.ecc.curves.c256,
        sjcl.ecc.curves.c256.field.fromBits(sjcl.codec.base64.toBits(serPriv))
    );
}

function verifyIntegrityOfVote(voteResult, postid) {
    var isValid = true;
    if(voteResult[1] && voteResult[0].postid != voteResult[1].postid) isValid = false;
    if(voteResult[2] && voteResult[0].postid != voteResult[2].postid) isValid = false;
    if(voteResult[3] && voteResult[0].postid != voteResult[3].postid) isValid = false;
    if(voteResult[1] && voteResult[2] && voteResult[1] == voteResult[2]) isValid = false;
    if(voteResult[1] && voteResult[3] && voteResult[1] == voteResult[3]) isValid = false;
    if(voteResult[2] && voteResult[3] && voteResult[2] == voteResult[3]) isValid = false;
    if(postid>=10 && voteResult[1] && !voteResult[3]) isValid = false;
    return isValid;
}

function analyzeVote(postid, vote, pref){
    var newVote = null;
    var correctCandidate = null;
    fetchedCandidates.every(thisCandidate => {
        try {
            newVote = sjcl.decrypt(thisCandidate.unserializedPrivateKey, vote);
            correctCandidate = thisCandidate;
            return false;
        } catch {
            return true;
        }
    });
    if (newVote==null){
        return {
            correctCandidate: correctCandidate,
            newVote: vote
        }
    }
    return {
        correctCandidate: correctCandidate,
        newVote: newVote
    };
}

function displayFinalTally() {
    result.forEach((postResult, postId) => {
        postResult.forEach(candidateResult => {
            console.log("FT,"+postId+","+candidateResult[0]+","+candidateResult[1]+","+candidateResult[2]+","+candidateResult[3]);
        });
    });
    process.exit();
}

function parseVotes(votes){
    var splitData;
    var percentageStatus;
    var postid;
    var voteResult;
    var analysisResult;
    var ballotId;
    var pref;

    votes.forEach((thisVote, index) => {
        percentageStatus = ((index+1)*100.0/votes.length);
        postid = parseInt(thisVote.postid);
        voteResult = [thisVote];
        pref = 1;
        vote = sjcl.decrypt(privateKeyOfCEO, thisVote.data);

        for(i=1; vote!=null; i++){
            analysisResult = analyzeVote(postid, vote, i);
            if(analysisResult.correctCandidate==null){
                splitData = analysisResult.newVote.split("$");
                ballotId = splitData[splitData.length-1];
                break;
            }else{
                voteResult[i] = analysisResult.correctCandidate;
            }
            vote = analysisResult.newVote;
        }

        if(verifyIntegrityOfVote(voteResult, postid)){
            if(voteResult[1]) result[postid][voteResult[1].roll][1] += 1;
            else result[postid][0][1]+=1;
            if(voteResult[2]) result[postid][voteResult[2].roll][2] += 1;
            if(voteResult[3]) result[postid][voteResult[3].roll][3] += 1;
            var p1 = voteResult[1] ? voteResult[1].roll : '0';
            var p2 = voteResult[2] ? voteResult[2].roll : '0';
            var p3 = voteResult[3] ? voteResult[3].roll : '0';
            console.log(postid+","+ballotId+","+p1+","+p2+","+p3+","+percentageStatus);
        } else {
            console.log(postid+","+ballotId+",Invalid vote,,,"+percentageStatus);
        }
        if(index==votes.length-1){
            displayFinalTally();
        }
    });
}

getVotes()
.then(votes => {
    fetchedVotes = votes;
    return getCandidates();
})
.then(candidates => {
    fetchedCandidates = candidates;
    fetchedCandidates.forEach(thisCandidate => {
        if(!result[thisCandidate.postid]) result[thisCandidate.postid] = [];
        if(!result[thisCandidate.postid][0]) result[thisCandidate.postid][0] = [0,0,0,0];
        result[thisCandidate.postid][thisCandidate.roll] = [thisCandidate.roll, 0, 0, 0];
        thisCandidate.unserializedPrivateKey = unserializePrivateKey(thisCandidate.privatekey);
    })
    return getCEO();
})
.then(ceo => {
    fetchedCEO = ceo;
    privateKeyOfCEO = unserializePrivateKey(ceo.privatekey);
})
.then(() => parseVotes(fetchedVotes))
.catch(err => {
    console.log(err);
    process.exit(1);
});