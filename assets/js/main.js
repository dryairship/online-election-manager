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
