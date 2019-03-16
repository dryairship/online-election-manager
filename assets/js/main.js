var userPassword;

function attemptLogin(){
    console.log("Called");
    var data = $('#loginform').serialize();
    console.log(data);
    $.ajax({
        type: "POST",
        url: "/login",
        data: data,
        cache: false,
        success: function(response){
            console.log("Success");
            console.log(response);
            if(response["success"]==0){
                document.getElementById("loginError").style="display:block";
            }else{
                alert("Successful login");
            }
        }
    });
}

function sendMail(){
    console.log("sendMail Called");
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
    console.log("reg Called");
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
                console.log(response);
                notif.html(response);
                notif.css("display","block");
                notif.removeClass("alert-info");
                notif.removeClass("alert-danger");
                notif.addClass("alert-success");
            },
            error: function(response){
                console.log(response);
                notif.html(response.responseText);
                notif.css("display","block");
                notif.removeClass("alert-info");
                notif.removeClass("alert-success");
                notif.addClass("alert-danger");
            }
        });
    }
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
