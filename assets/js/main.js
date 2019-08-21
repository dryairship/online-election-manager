var userPassword;
var userRoll;
var userData;
var allPosts;
var votesCandidateNames = [["MEOW"]];
var votesCandidatePublicKeys = [["MEOW"]];
var ballotIDs = [];
var encryptedBallotIDs = [];
var finalVotes = [];
var unserializedPublicKeyOfCEO;

// Function to chack username and password and call the appropriate functions.
function attemptLogin(){
    var data = $('#loginform').serializeArray();
    var roll = data[0].value;
    var pass = data[1].value;
    var passHash = sjcl.codec.hex.fromBits(sjcl.hash.sha256.hash(pass));
    $.ajax({
        type: "POST",
        url: "/users/login",
        data: $.param({
                'roll':roll,
                'pass':passHash
            }),
        cache: false,
        success: function(response){
            userData = response;
            userPassword = pass;
            userRoll = roll;
            if(roll == "CEO"){
                $("body").load("ceo.html", initializeCEO);
            }else if(roll[0] == 'P'){
                $("body").load("candidateHome.html", initializeCandidateHome);
            }else{
                if(userData["Voted"]){
                    $("body").load("home.html", showUserHasVoted);
                }else if(userData["State"]==1){
                    $("body").load("home.html", loadPosts);
                    unserializedPublicKeyOfCEO = unserializePublicKey(userData.CEOKey);
                }else{
                    $("body").load("home.html", showIncorrectState);
                }
            }
        },
        error: function(response){
            // Show error message in the error box.
            document.getElementById("loginError").style="display:block";
            document.getElementById("loginError").innerHTML=response.responseText;
        }
    });
}

// Fetch posts from the server and store them.
function loadPosts(){
    $.ajax({
        type: "GET",
        url: "/election/getVotablePosts",
        cache: false,
        success: function(response){
            allPosts = response;
            response.forEach(loadThisPost);
        }
    });
}

// Load the candidates for a particular post on the client's device.
function loadThisPost(post, ind, all){
    var postid = post["PostID"];
    $("#postsTable>tbody").append("<tr><td align='center' id='post"+postid+"'></td></tr>");
    $("#post"+postid).load("candidatePanel.html", function(){
        $("#post"+postid+">#candidatePanel>.postname").html(post["PostName"])
        post["Candidates"].forEach(function(candidate, cid, allC){
            cid = cid+1;
            var candid = "post"+postid+"-cand"+cid;
            $('#post'+postid+'>#candidatePanel').append("<div id='"+candid+"'></div>");
            $('#'+candid).load("election/getCandidateCard/"+candidate, function(){
                $('#'+candid+' #voteButton').on('click', function() {
                    // If this button is clicked, it calls the vote() function with itself as the parameter.
                    vote(this);
                });
            });
        });
    });
}

// Reload the candidates after preferences have been reset,
function reloadPost(postid){
    var post;
    allPosts.forEach(function(el,ind,all){
        if(el["PostID"]==postid){
            post = el;
        }
    });
    votesCandidateNames[parseInt(postid)] = [];
    votesCandidatePublicKeys[parseInt(postid)] = [];
    $("#post"+postid).load("candidatePanel.html", function(){
        $("#post"+postid+">#candidatePanel>.postname").html(post["PostName"])
        post["Candidates"].forEach(function(candidate, cid, allC){
            cid = cid+1;
            var candid = "post"+postid+"-cand"+cid;
            $('#post'+postid+'>#candidatePanel').append("<div id='"+candid+"'></div>");
            $('#'+candid).load("election/getCandidateCard/"+candidate, function(){
                $('#'+candid+' #voteButton').on('click', function() {
                    vote(this);
                });
            });
        });
    });
}

// Calls the API to send verification mail.
function sendMail(){
    var notif = $('#mailNotification');
    notif.html("Sending mail...");
    notif.css("display","block");
    $.ajax({
        type: "GET",
        url: "/users/mail/"+document.getElementById("rollForAuthCode").value,
        cache: false,
        success: function(response){
            notif.html(response);
            notif.css("display","block");
            notif.removeClass("alert-info");
            notif.removeClass("alert-danger");
            notif.addClass("alert-success");
        },
        error: function(response){
            notif.html(response.responseText);
            notif.css("display","block");
            notif.removeClass("alert-info");
            notif.removeClass("alert-success");
            notif.addClass("alert-danger");
        }
    });
}

// Method to complete user registration.
function register(){
    var notif = $('#regNotification');
    var data = $('#registrationform').serializeArray();
    var roll = data[0].value;
    var pass = data[1].value;
    var pass2 = data[2].value;
    var auth = data[3].value;
    var passHash = sjcl.codec.hex.fromBits(sjcl.hash.sha256.hash(pass));

    if(pass!=pass2){
        notif.html("The passwords do not match.");
        notif.css("display","block");
        notif.removeClass("alert-info");
        notif.removeClass("alert-success");
        notif.addClass("alert-danger");
    }else{
        notif.html("Registering voter...")
        notif.css("display","block");
        notif.removeClass("alert-danger");
        notif.removeClass("alert-success");
        notif.addClass("alert-info");
        userPassword = pass;
        $.ajax({
            type: "POST",
            url:  "/users/register",
            data: $.param({
                'roll':roll,
                'pass':passHash,
                'auth':auth
            }),
            cache: false,
            success: function(response){
                showLoginForm();
                $("#loginError").removeClass("alert-danger");
                $("#loginError").addClass("alert-success");
                document.getElementById("loginError").style="display:block";
                document.getElementById("loginError").innerHTML="Registration successful.<br>You may now log in.";
            },
            error: function(response){
                notif.html(response.responseText);
                notif.css("display","block");
                notif.removeClass("alert-info");
                notif.removeClass("alert-success");
                notif.addClass("alert-danger");
            }
        });
    }
}

// Store the user's preferences in a global variable.
function vote(button){
    var postid = button.attributes["postid"].value;
    var serpubkey = button.attributes["pubkey"].value;
    var pubkey = unserializePublicKey(serpubkey);
    var pref = button.value[0];
    var candName = button.parentNode.firstChild.textContent;
    button.parentNode.parentNode.parentNode.parentNode.parentNode.remove();

    // To change the text in other buttons.
    document.querySelectorAll("#post"+postid+" #voteButton").forEach(function(el, ind, all){
        if(el.value[0]=="1")    el.value="2nd Preference";
        else if(el.value[0]=="2")    el.value="3rd Preference";
        else el.remove()
    });

    // Show the newly added preference in the list.
    showVoted(candName,pref, postid);
    var intID = parseInt(postid);
    var intPref = parseInt(pref);
    if(votesCandidateNames[intID]==undefined){
        votesCandidateNames[intID] = []
    }
    if(votesCandidatePublicKeys[intID]==undefined){
        votesCandidatePublicKeys[intID] = []
    }
    votesCandidateNames[intID][intPref] = candName;
    votesCandidatePublicKeys[intID][intPref] = pubkey;
}

// To display the current preferences.
function showVoted(candName, pref, postid){
    $("#post"+postid+" #heading").html("Preferences : <a href='#' onclick='reloadPost("+postid+"); return false;' class='badge badge-danger badge-pill'>Reset Preferences</a>");
    $("#post"+postid+" .list-group").append("<li id='pref"+pref+"' class='list-group-item'>"+pref+") "+candName+"</li>");
}

// Show final votes before submission.
function confirmVotes(){
    $("#confirmVotes .modal-body").html("");
    allPosts.forEach(function(post, ind, all){
        var pname = post["PostName"];
        var pid = parseInt(post["PostID"]);
        if(pid>=10 && (votesCandidateNames[pid]==undefined || votesCandidateNames[pid].length<4)){
            $("#badVotes").modal('show');
            return;
        }
        $("#confirmVotes .modal-body").append("<dl id='votes"+pid+"'><dt>"+pname+"</dt></dl>");
        if(votesCandidateNames[pid]==undefined || votesCandidateNames[pid].length==0){
            $("#votes"+pid).append("<dd>NOTA</dd>");
        }else{
            votesCandidateNames[pid].forEach(function(cand, indC, allC){
                $("#votes"+pid).append("<dd>"+indC+") "+cand+"</dd>");
            });
        }
        if(ind==all.length-1) $("#confirmVotes").modal('show');
    });
}

// Decrypt Ballot IDs that a voter can use to verify his vote.
function decryptBallotIDs(){
    var alertBox = $(".alert");
    userData.BallotID.forEach(function(el, ind, all){
        encryptedBallotIDs[el.PostID] = el.BallotString;
        ballotIDs[el.PostID] = decryptFromPassword(el.BallotString);
        alertBox.html(alertBox.html()+"<br>Ballot ID for Post "+el.PostID+" = "+ballotIDs[el.PostID]);
    });
}

// Convert Public Key from JSON to Base64
function serializePublicKey(pub){
    return sjcl.codec.base64.fromBits(pub.get().x.concat(pub.get().y));
}

// Convert Private Key from JSON to Base64
function serializePrivateKey(priv){
    return sjcl.codec.base64.fromBits(priv.get());
}

// Convert Public Key from Base64 to JSON
function unserializePublicKey(serPub){
    return new sjcl.ecc.elGamal.publicKey(
        sjcl.ecc.curves.c256,
        sjcl.codec.base64.toBits(serPub)
    );
}

// Convert Private Key from Base64 to JSON
function unserializePrivateKey(serPriv){
    return new sjcl.ecc.elGamal.secretKey(
        sjcl.ecc.curves.c256,
        sjcl.ecc.curves.c256.field.fromBits(sjcl.codec.base64.toBits(serPriv))
    );
}

// Generate a public-private key pair.
function generateKeyPair(){
    return sjcl.ecc.elGamal.generateKeys(256);
}

// Decrypt something with the user's password as the key for symmetrically encrypted objects.
function decryptFromPassword(something){
    return sjcl.decrypt(userPassword, something);
}

// Encrypt something symmetrically with the user's password as the key.
function encryptWithPassword(something){
    return sjcl.encrypt(userPassword, something);
}

// Call the appropriate functions to submit votes.
function submitVotes(){
    encryptVotes();
    sendVotes();
}

// Creates BallotIDs and encrypts it using public keys of candidates and CEO.
function encryptVotes(){
    allPosts.forEach(function(post, ind, all){
        var ballotID = getRandomString();
        var pid = parseInt(post["PostID"]);
        ballotIDs[pid] = ballotID;
        encryptedBallotIDs[pid] = encryptWithPassword(ballotID);
        if(votesCandidateNames[pid] == undefined || votesCandidateNames[pid].length==0){
            currentVote = "$".concat(ballotID);
        }else{
            currentVote = votesCandidateNames[pid].join("$").concat("$").concat(ballotID);
            if(votesCandidatePublicKeys[pid][3]!=undefined){
                currentVote = sjcl.encrypt(votesCandidatePublicKeys[pid][3], currentVote);
            }
            if(votesCandidatePublicKeys[pid][2]!=undefined){
                currentVote = sjcl.encrypt(votesCandidatePublicKeys[pid][2], currentVote);
            }
            if(votesCandidatePublicKeys[pid][1]!=undefined){
                currentVote = sjcl.encrypt(votesCandidatePublicKeys[pid][1], currentVote);
            }
        }
        currentVote = sjcl.encrypt(unserializedPublicKeyOfCEO, currentVote);
        finalVotes[pid] = currentVote;
    });
}

// Submit votes to the server.
function sendVotes(){
    var voteData = [];
    allPosts.forEach(function(post, ind, all){
        var pid = parseInt(post["PostID"]);
        voteData = voteData.concat({
            "PostID"        : pid,
            "BallotString"  : encryptedBallotIDs[pid],
            "VoteData"      : finalVotes[pid]
        });
    });
    $.ajax({
        type: "POST",
        url: "/election/submitVote",
        dataType: 'json',
        data: JSON.stringify(voteData),
        cache: false,
        success: function(response){
            showUserHasVoted();
        }
    });
}

// Creates a random string.
function getRandomString(){
    var randBytes = sjcl.random.randomWords(8);
    var randHex = sjcl.codec.hex.fromBits(randBytes);
    return randHex;
}

// Toggle from login page to registration page.
function showRegistrationForm(){
    document.getElementById("loginContainer").style="display:none";
    document.getElementById("registrationContainer").style="display:block";
}

// Toggle from registration page to login page.
function showLoginForm(){
    document.getElementById("loginContainer").style="display:block";
    document.getElementById("registrationContainer").style="display:none";
}

// Show that the user has submitted his vote.
function showUserHasVoted(){
    $("body").addClass("d-flex");
    $("body").html("<div class=\"alert alert-success mx-auto my-auto d-inline-flex\">Your vote has been submitted.</div>");
    decryptBallotIDs();
}

// Show that the time for voting has ended/not started.
function showIncorrectState(){
    var msg;
    if(userData["State"]==0) msg = "Voting has not yet started.";
    else msg = "Voting period is over now."
    $("body").addClass("d-flex");
    $("body").html("<div class=\"alert alert-danger mx-auto my-auto d-inline-flex\">"+msg+"</div>");
}

// Set up home page when a candidate logs in.
function initializeCandidateHome(){
    var button = $("button");
    var text = $(".welcomed");
    if(userData.KeyState==0 && userData.State=="0" ){
        text.html("Click on the button to confirm your candidature.");
    }else if(userData.State=="1" || (userData.State=="0" && userData.KeyState==1)){
        button.remove();
        text.html("Wait for the elections to end.");
    }else if(userData.State=="2" && userData.KeyState==1){
        text.html("Click on the button to send your private key for decryption.");
        button.html("Send private key");
        button.unbind('click');
        button.on('click', sendPrivateKey);
    }else if(userData.KeyState==2){
        text.html("All the best for the results.");
        button.remove();
    }else{
        button.remove();
        text.removeClass("alert-success");
        text.addClass("alert-danger");
        text.html("An error occured.");
    }
}

// Generates a public-private key pair for a candidate and submits it.
function confirmCandidature(){
    var pair = generateKeyPair();
    var pubKey = serializePublicKey(pair.pub);
    var privKey = encryptWithPassword(serializePrivateKey(pair.sec));
    $.ajax({
        type: "POST",
        url: "/candidate/confirmCandidature",
        contentType: 'application/json; charset=utf-8',
        data: JSON.stringify({
            'pubkey': pubKey,
            'privkey': privKey
        }),
        cache: false,
        success: function(response){
            showCandidatureConfirmed();
        }
    });
}

// Sends the unencrypted private key to the server.
function sendPrivateKey(){
    var privKey = decryptFromPassword(userData.PrivateKey);
    $.ajax({
        type: "POST",
        url: "/candidate/declarePrivateKey",
        contentType: 'application/json; charset=utf-8',
        data: JSON.stringify({
            'pubkey': userData.PublicKey,
            'privkey': privKey
        }),
        cache: false,
        success: function(response){
            showPrivateKeySubmitted();
        }
    });
}

// Show the message that the private key has been submitted.
function showPrivateKeySubmitted(){
    $(".welcomed").html("Your key has been submitted. All the best for the results.");
    $("button").remove();
}

// Show the message that the candidature has been confirmed.
function showCandidatureConfirmed(){
    $(".welcomed").html("Your candidature has been confirmed. Wait for the elections to end.");
    $("button").remove();
}

// Load login page and seed the random bytes generator in SJCL.
$(function(){
    $("body").load("login.html");
    var arr = new Uint32Array(128);
    crypto.getRandomValues(arr);
    sjcl.random.addEntropy(arr, 1024, "crypto.randomBytes");
});

