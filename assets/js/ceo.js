var candidates;
var totalPosts;
var privateKeyOfCEO;
var names = [];
var usernames = [];
var privateKeys = [];

function calculate(){
    $("button").html("Display Final Tally");
    $("button").unbind('click');
    $("button").on('click', showResults);
    fetchPosts();
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
