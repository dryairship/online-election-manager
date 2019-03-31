var privateKeyOfCEO;

function calculate(){

}

function startVoting(){

}

function stopVoting(){

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
