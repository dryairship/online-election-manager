var userPassword;
var userRoll;
var allPosts;
var votesCandidateNames = [["MEOW"]];
var votesCandidatePublicKeys = [["MEOW"]];

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
            userPassword = pass;
            userRoll = roll;
            $("body").load("home.html", loadPosts);
        },
        error: function(response){
            document.getElementById("loginError").style="display:block";
            document.getElementById("loginError").innerHTML=response.responseText;
        }
    });
    setTimeout(function(){setVoteButtonsClickable("all");}, 3000);
}

function setVoteButtonsClickable(postid){
    if(postid=="all"){
        $('input#voteButton').on('click', function() {
            vote(this);
        });
        $(".loading")[0].remove();
    }else{
        $('#post'+postid+' #voteButton').on('click', function() {
            vote(this);
        });
    }
}

function loadPosts(){
    console.log("Loading Posts");
    $.ajax({
        type: "GET",
        url: "/election/getVotablePosts/"+userRoll,
        cache: false,
        success: function(response){
            allPosts = response;
            response.forEach(loadThisPost);
        }
    });
}

function loadThisPost(post, ind, all){
    var postid = post["PostID"];
    $("#postsTable>tbody").append("<tr><td align='center' id='post"+postid+"'></td></tr>");
    $("#post"+postid).load("candidatePanel.html", function(){
        $("#post"+postid+">#candidatePanel>.postname").html(post["PostName"])
        post["Candidates"].forEach(function(candidate, cid, allC){
            cid = cid+1;
            var candid = "post"+postid+"-cand"+cid;
            $('#post'+postid+'>#candidatePanel').append("<div id='"+candid+"'></div>");
            $('#'+candid).load("election/getCandidateCard/"+candidate);
        });
    });
}

function reloadPost(postid){
    var post;
    allPosts.forEach(function(el,ind,all){
        if(el["PostID"]==postid){
            post = el;
        }
    });
    $("#post"+postid).load("candidatePanel.html", function(){
        $("#post"+postid+">#candidatePanel>.postname").html(post["PostName"])
        post["Candidates"].forEach(function(candidate, cid, allC){
            cid = cid+1;
            var candid = "post"+postid+"-cand"+cid;
            $('#post'+postid+'>#candidatePanel').append("<div id='"+candid+"'></div>");
            $('#'+candid).load("election/getCandidateCard/"+candidate);
        });
    });
    setTimeout(function(){setVoteButtonsClickable(postid)}, 1000);
}

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
}

function vote(button){
    var postid = button.attributes["postid"].value;
    var pubkey = button.attributes["pubkey"].value;
    var pref = button.value[0];
    var candName = button.parentNode.firstChild.textContent;
    button.parentNode.parentNode.parentNode.parentNode.parentNode.remove();
    console.log("Voted for Candidate "+candName+" with pubkey="+pubkey+" as pref number "+pref+" on postid="+postid);
    document.querySelectorAll("#post"+postid+" #voteButton").forEach(function(el, ind, all){
        if(el.value[0]=="1")    el.value="2nd Preference";
        else if(el.value[0]=="2")    el.value="3rd Preference";
        else el.remove()
    });
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

function showVoted(candName, pref, postid){
    $("#post"+postid+" #heading").html("Preferences : <a href='#' onclick='reloadPost("+postid+"); return false;' class='badge badge-danger badge-pill'>Reset Preferences</a>");
    $("#post"+postid+" .list-group").append("<li id='pref"+pref+"' class='list-group-item'>"+pref+") "+candName+"</li>");
}

function init(){
    console.log("Initializing.");
}
function showRegistrationForm(){
    document.getElementById("loginContainer").style="display:none";
    document.getElementById("registrationContainer").style="display:block";
}

function showLoginForm(){
    document.getElementById("loginContainer").style="display:block";
    document.getElementById("registrationContainer").style="display:none";
}

$(function(){
    $("body").load("login.html");
});
