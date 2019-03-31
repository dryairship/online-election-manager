var candidates;
var totalPosts;
var privateKeyOfCEO;
var names = [];
var scores = [];
var result = [];
var votesStr = [];
var usernames = [];
var privateKeys = [];

function calculate(){
    $("button").html("Display Final Tally");
    $("button").unbind('click');
    $("button").on('click', showResults);
    fetchPosts();
}

function findAllResults(){
    totalPosts.forEach(function(el, ind, all){
        $("#postsTable>tbody").append("<tr><td align='center' id='post"+el.postid+"'></td></tr>");
        $("#post"+el.postid).load("resultsTable.html", function(){
            $("#post"+el.postid+" legend").html(el.postid+") "+el.postname);
            findResult(el.postid);
        });
    });
}

function findResult(postid){
    votesStr[postid].forEach(function(el, ind, all){
        var newvote = [], vote = sjcl.decrypt(privateKeyOfCEO, el);
        var pref = 1;
        for(i=1; vote!=null; i++){
            newvote[i] = analyzeVote(postid, vote, i);
            if(newvote[i]==vote){
                showVote(postid, newvote, vote, i);
                break;
            }
            vote = newvote[i];
        }
    });
}

function analyzeVote(postid, vote, pref){
    var newVote = null;
    privateKeys[postid].every(function(el, ind, all){
        try {
            newVote = sjcl.decrypt(el, vote);
            if (result[usernames[postid][ind]]==undefined){
                result[usernames[postid][ind]]=[0,0,0,0];
            }
            if (scores[usernames[postid][ind]]==undefined){
                scores[usernames[postid][ind]]=0;
            }
            result[usernames[postid][ind]][pref] += 1;
            scores[usernames[postid][ind]]+=Math.pow(10,(3-pref)*3);
            return false;
        } catch {
            return true;
        }
    });
    if (newVote==null){
        return vote;
    }
    return newVote;
}

function showVote(postid, arr, vote, size){
    var voteData = vote.split("$");
    var voteID = voteData[size];
    $("#post"+postid+" table>tbody").append("<tr id='"+voteID+"'><td align='center'>"+voteID+"</td></tr>");
    for(var i=1; i<=3; i++){
        $("#"+voteID).append("<td align='center'>"+voteData[i]+"</td>");
    }
}

function parseVotes(votes, callback){
    votes.forEach(function(el, ind, all){
        if(votesStr[el.postid] == undefined){
            votesStr[el.postid] = [];
        }
        votesStr[el.postid].push(el.data);
    });
    callback();
}

function parseCandidatesData(){
    candidates.forEach(function(el, ind, all){
        if(privateKeys[el.PostID]==undefined){
            privateKeys[el.PostID] = [];
            names[el.PostID] = [];
            usernames[el.PostID] = [];
        }
        try {
            privateKeys[el.PostID].push(unserializePrivateKey(el.PrivateKey));
        } catch {}
        names[el.PostID].push(el.Name);
        usernames[el.PostID].push(el.Username);
    });
}

function fetchVotes(){
    $.ajax({
        type: "GET",
        url:  "ceo/fetchVotes/",
        cache:false,
        success: function(response){
            parseVotes(response,findAllResults);
        }
    });
}

function fetchPosts(){
    $.ajax({
        type: "GET",
        url:  "/ceo/fetchPosts",
        cache:false,
        success: function(response){
            totalPosts = response;
            fetchCandidates();
        },
        error: function(response){
            alert(response.responseText);
        }
    });
}

function fetchCandidates(){
    $.ajax({
        type: "GET",
        url:  "/ceo/fetchCandidates",
        cache:false,
        success: function(response){
            candidates = response;
            parseCandidatesData();
            fetchVotes();
        }
    });
}

function startVoting(){
    var pair = generateKeyPair();
    userData.publickey = serializePublicKey(pair.pub);
    privateKeyOfCEO = pair.sec;
    userData.privatekey = serializePrivateKey(pair.sec);
    $.ajax({
        type: "POST",
        url:  "/ceo/startVoting",
        data: JSON.stringify({
            'pubkey': userData.publickey,
            'privkey': userData.privatekey
        }),
        contentType: 'application/json; charset=utf-8',
        cache:false,
        success: function(response){
            $("button").html("Stop Voting");
            $("button").unbind("click");
            $("button").on('click', stopVoting);
        },
        error: function(response){
            alert(response.responseText);
        }
    });
}

function stopVoting(){
    $.ajax({
        type: "POST",
        url:  "/ceo/stopVoting",
        cache:false,
        success: function(response){
            $("button").html("Calculate Results");
            $("button").on('click', calculate);
        },
        error: function(response){
            alert(response.responseText);
        }
    });
}

function initializeCEO(){
    $.ajax({
        type: "GET",
        url:  "/election/getElectionState",
        cache:false,
        success: function(response){
            var button = $("button");
            if(response=="0"){
                button.html("Start Voting");
                button.on('click', startVoting);
            }else if(response=="1"){
                button.html("Stop Voting");
                button.on('click', stopVoting);
            }else if(response=="2"){
                button.html("Calculate Results");
                button.on('click', calculate);
            }
        },
        error: function(response){
            alert(response.responseText);
        }
    });
    privateKeyOfCEO = unserializePrivateKey(userData.privatekey);
}
