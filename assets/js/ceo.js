var candidates;
var totalPosts;
var privateKeyOfCEO;
var names = {};
var postNames = {};
var postResults = {};
var scores = [];
var result = [];
var votesStr = [];
var usernames = [];
var privateKeys = [];
var finalUsernames = [];
var progress;
// Calculate the results.
function calculate(){
    $("button").html("Display Final Tally");
    $("button").unbind('click');
    $("button")[0].style.display='none';
    $("#ceomessage")[0].style.display='block';
    $("button").on('click', fetchResults);
    $.ajax({
        type: "GET",
        url:  "/ceo/calculateResult",
        cache:false,
        success: function(response){
            startProgress();
        }
    });
}

function startProgress() {
    setTimeout(() => {
        if(progress=="100.000%"){
            $("button")[0].style.display='block';
            $("#ceomessage")[0].style.display='none';
        } else {
            checkProgress();
        }
    }, 1000);
}

function checkProgress() {
    $.ajax({
        type: "GET",
        url:  "/ceo/resultProgress",
        cache:false,
        success: function(response){
            progress = response;
            $("#ceomessage")[0].innerHTML = "Progress : "+progress;
        }
    });
    startProgress();
}

function fetchResults() {
    $("button")[0].style.display='none';
    $.ajax({
        type: "GET",
        url:  "/ceo/getResult",
        cache:false,
        success: function(response){
            result = response;
            processResult();
        }
    });
}

function processResult() {
    result.forEach(el => {
        if(!postResults[el.PostID]) {
            postResults[el.PostID] = [];
        }
        if (el.Candidate == 0) {
            postResults[el.PostID].unshift({
                candidateRoll: el.Candidate,
                candidateName: names[el.Candidate],
                p1: el.Preference1,
                p2: el.Preference2,
                p3: el.Preference3
            });
        } else {
            postResults[el.PostID].push({
                candidateRoll: el.Candidate,
                candidateName: names[el.Candidate],
                p1: el.Preference1,
                p2: el.Preference2,
                p3: el.Preference3
            });
        }
    });
    displayResults();
}


/*
// Call the findResult function for all posts and append the data to the page.
function findAllResults(){
    totalPosts.forEach(function(el, ind, all){
        $("#postsTable>tbody").append("<tr><td align='center' id='post"+el.postid+"'></td></tr>");
        $("#post"+el.postid).load("resultsTable.html", function(){
            $("#post"+el.postid+" legend").html(el.postid+") "+el.postname);
            findResult(el.postid);
        });
    });
}

// Decrypts votes for a particular post.
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

// Decrypts a particular vote.
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

// Displays the details of the vote on the page.
function showVote(postid, arr, vote, size){
    var voteData = vote.split("$");
    var voteID = voteData[size];
    $("#post"+postid+" table>tbody").append("<tr id='"+voteID+"'><td align='center'>"+voteID+"</td></tr>");
    for(var i=1; i<=3; i++){
        if(i==size || voteData[i]==undefined)
            $("#"+voteID).append("<td align='center'>No Choice</td>");
        else
            $("#"+voteID).append("<td align='center'>"+voteData[i]+"</td>");
    }
}
*/
/*
// Find out winners and runners up from the decrypted votes.
function parseResults(callback){
    $("button").remove();
    totalPosts.forEach(function(el, ind){
        finalUsernames[el.postid] = [];
        usernames[el.postid].forEach(function(uname, uind){
            for(var i=3; i>=1; i--){
                var cs = parseInt(scores[finalUsernames[el.postid][i]]) || 0;
                if (scores[uname]>cs){
                    var tmp = finalUsernames[el.postid][i];
                    finalUsernames[el.postid][i] = uname;
                    finalUsernames[el.postid][i+1] = tmp;
                }
            }
        });
    });
    callback();
}
*/
// Display winners and runners up.
function displayResults(){
    $("#postsTable>tbody").append("<tr><td align='center'><div class=\"alert alert-success mx-auto my-auto d-inline-flex\">Final Tally</div></td></tr>");
    totalPosts.forEach(function(el, ind){
        $("#postsTable>tbody").append("<tr><td align='center' id='result"+el.postid+"'></td></tr>");
        $("#result"+el.postid).load("finalTally.html", function(){
            $("#result"+el.postid+" legend").html(el.postid+") "+el.postname);
            var candD = postResults[el.postid][0];
            if(parseInt(el.postid)<10) $("#result"+el.postid+" table>tbody").append("<tr id='res"+el.postid+"posNOTA'><td align='center'>NOTA</td><td align='center'>"+candD.candidateName+"</td><td align='center'>"+candD.p1+"</td><td align='center'>"+candD.p2+"</td><td align='center'>"+candD.p3+"</td></tr>");
            var i0 = (el.postid<10?1:0);
            for(var i=i0; i<postResults[el.postid].length; i++){
                candD = postResults[el.postid][i];
                $("#result"+el.postid+" table>tbody").append("<tr id='res"+el.postid+"posNOTA'><td align='center'>"+(i-i0+1)+"</td><td align='center'>"+candD.candidateName+"</td><td align='center'>"+candD.p1+"</td><td align='center'>"+candD.p2+"</td><td align='center'>"+candD.p3+"</td></tr>");
            }
        });
    });
}

// Get Name of a candidate.
function getName(candidate){
    var name = "";
    candidates.every(function(el){
        if(el.Username==candidate){
            name = el.Name;
            return false;
        }
        return true;
    });
    return name;
}
/*
// Properly store fetched votes from the server into a suitable format in the global variables.
function parseVotes(votes, callback){
    votes.forEach(function(el, ind, all){
        if(votesStr[el.postid] == undefined){
            votesStr[el.postid] = [];
        }
        votesStr[el.postid].push(el.data);
    });
    callback();
}
*/

// Fetch posts from the server.
function fetchPosts(){
    $.ajax({
        type: "GET",
        url:  "/ceo/fetchPosts",
        cache:false,
        success: function(response){
            totalPosts = response;
            totalPosts.forEach(el => {
                postNames[el.postid] = el.postname;
            })
        },
        error: function(response){
            alert(response.responseText);
        }
    });
}

// Fetch candidates from the server.
function fetchCandidates(){
    $.ajax({
        type: "GET",
        url:  "/ceo/fetchCandidates",
        cache:false,
        success: function(response){
            candidates = response;
            candidates.forEach(function(el, ind, all){
                names[el.Roll] = el.Name;
            });
            names[0] = "NOTA";
        }
    });
}

// Create a public-private key pair for the CEO and start the voting process.
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
            $("#ceowelcome").append(" Voting period has started now.");
        },
        error: function(response){
            alert(response.responseText);
        }
    });
}

// Tell the server to stop accepting votes.
function stopVoting(){
    $.ajax({
        type: "POST",
        url:  "/ceo/stopVoting",
        cache:false,
        success: function(response){
            $("#ceowelcome").append(" Voting period has stopped now.");
            $("button").html("Calculate Results");
            $("button").unbind("click");
            $("button").on('click', calculate);
        },
        error: function(response){
            alert(response.responseText);
        }
    });
}

// Set up CEO's homepage.
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
    fetchCandidates();
    fetchPosts();
}
