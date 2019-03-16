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
    return false;
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
