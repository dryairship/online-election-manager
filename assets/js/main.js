var userPassword;
var userRoll;

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
            if(response["success"]==0){
                document.getElementById("loginError").style="display:block";
            }else{
                userPassword = pass;
                userRoll = roll;
                $("body").load("home.html", loadPosts);
            }
        }
    });
    setTimeout(setVoteButtonsClickable, 3000);
}

function setVoteButtonsClickable(){
    $(".loading")[0].parentNode.removeChild($(".loading")[0]);
    $('input#voteButton').on('click', function() {
        var attr = this.attributes;    
        vote(attr["postid"].value,attr["pubkey"].value,1);
    });
}

function loadPosts(){
    console.log("Loading Posts");
    $.ajax({
        type: "GET",
        url: "/election/getVotablePosts/"+userRoll,
        cache: false,
        success: function(response){
            response.forEach(loadThisPost);
        }
    });
}

function loadThisPost(post, ind, allPosts){
    ind=ind+1;
    var postid = "post"+ind;
    $("#postsTable>tbody").append("<tr><td align='center' id='"+postid+"'></td></tr>");
    $("#"+postid).load("candidatePanel.html", function(){
        $("#"+postid+">#candidatePanel>.postname").html(post["PostName"])
        post["Candidates"].forEach(function(candidate, cid, allC){
            cid = cid+1;
            var candid = postid+"-cand"+cid;
            $('#'+postid+'>#candidatePanel').append("<div id='"+candid+"'></div><pre>  </pre>");
            $('#'+candid).load("election/getCandidateCard/"+candidate);
        });
    });
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

function vote(postid, pubkey, pref){
    console.log("Voted for Candidate with pubkey="+pubkey+"as pref number "+pref+" on postid="+postid);
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
    $("body").load("login.html")
});
