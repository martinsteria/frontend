﻿//The current user operating the website
var user

//The current module the user is editing
var module

//Is the module from the library or from the users collection
//VALS: "USER" or "LIB"
var moduleSource

//API endpoints
var apiRoot = "http://52.169.232.92/api"
var modules = apiRoot + "/library"
var users = apiRoot + "/users"
var deployment = apiRoot + "/deploy"

$(document).ready(function () {
    $.ajaxSetup({ cache: false })
    $("#login-view").hide()
    $("#login-view").fadeIn("slow")
    $("#usernameInput").focus(); 
    $("#loginBtn").click(logIn)
    $("#library-view").hide()
    $("#variables-view").hide()
    $("#deployment-view").hide()
    $("#description").hide()
    $("#newUser").hide()

    //"logg inn" knapp aktiveres ved å trykke enter i inputbox
    $('#usernameInput').keypress(function (e) {
        if (e.keyCode == 13)
            $('#loginBtn').click();
    });
})

/*Checks wether the user exist by getting a list of all users from the server and comparing them with the entered username. 
If the user exist it calls importUserModules(), if not it calls makeUser()*/

function logIn() {
    user = $("#usernameInput").val() 
	  $.get({
		    url: users, //Get-request to server for url= ../api/users
		    success: function(result) {
			      var e = false;
            if (result == null) { 
                makeUser(user)
                return
            }
			      for (i=0; i< result.length ; i++) { // for all usernames check if it equals the enteres username
				        console.log(result[i]+"="+user);
				        if (user == result[i]){
					          e = true; 
				        }
			      }
			      if (e == true) { 
			          importUserModules(users, user);
			      }
			      else { 
				        makeUser(user);
			      }
		    }
	  })
    $("#login-view").fadeOut("slow", function() {
        importLibraryModules(modules)
        $("#library-view").fadeIn("slow")
    })
    var content = "<span>" + "Du er logget inn som " + user + "" + "</span>";
    $("#bruker").html(content)
    $("#bruker").show()
}

/* sends a post request to the server with the entered username*/
function makeUser(user) {
    console.log("make new user")
    $.post({
        url: users + "?user=" + user,
        success: function (result) {
            console.log("new User: " + user)
            var content= "<span>" + "Ny Bruker" + "</span><br>" + "Ny bruker opprettet for "+ user;
            $("#newUser").html(content)
            $("#newUser").show()
        }
    })
}

/*Sends a get-request for a JSON-file containing all available modules. Enteres the available modules to #library*/
function importLibraryModules(path) {
    $.getJSON(path, function (resultModules) {
        var content = ""
        if (resultModules != null){
        content += "<option selected disabled hidden>Biblioteksmoduler...</option>"
        for (i = 0; i < resultModules.length; i++) {
            content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].name + "</option>"
        }
        $("#library").html(content);
        $("#library").change(function() {
            module = $("#library option:selected").text()
            moduleSource = "LIB"
            showModule(path +"?module=" + module)
            $("#variables-view").fadeIn("slow")
        })
    }
    });
}

/*Sends a get-request for a JSON-file containing the modules for this user. Enteres the available modules to #userLibrary*/
function importUserModules(path, user) {
    $.getJSON(path + "?user=" + user, function (resultModules) {
        console.log(resultModules)
        console.log(resultModules);
        var content = ""
        content += "<option selected disabled hidden>Brukermoduler...</option>"
        if (resultModules == null) {
            return
        }
        for (i = 0; i < resultModules.length; i++) {
            content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].name + "</option>"
        }

        $("#userLibrary").html(content);
        $("#userLibrary").change(function () {
            module = $("#userLibrary option:selected").text()
            moduleSource = "USER"
            showModule(path + "?user=" + user +  "&module=" + module)
            $("#variables-view").fadeIn("slow")
        })
    });
}


/* Sends a get request for a JSON-file containing the variables for the selected module. Adds each variables to the table #variablesTable*/
function showModule(path) {
    $.getJSON(path, function(result) {
        var content = "2. " 
        var name = result.name
        content += name.charAt(0).toUpperCase() + name.slice(1);
        $("#moduleName").html(content)
        $("#moduleDescription").html(result.description)
        var myTable = ""
				myTable += "<thead><tr><th>Navn</th><th>Verdi</th></tr></thead>"
        for (i = 0; i < result.variables.length; i++) {
            console.log(result.variables[i])
            var value = result.variables[i].value
            if (value == "") {
                value = result.variables[i].defaultValue
            }

            var textInputBox = '<input type="text" class="form-control" value="' + value + '" />';
            myTable += '<tr>'
            myTable += '<td><a href="#" data-placement="left" data-toggle="tooltip" title="' + result.variables[i].description + '">' + result.variables[i].name + '</a></td>'
            myTable += '<td>' + textInputBox + '</td>'
            myTable += '</tr>'
        }
        $("#variablesTable").html(myTable)
        $("#showDeployment").click(showDeployment)
        $("[data-toggle=\"tooltip\"]").tooltip();
    });
}

/* Shows deployment-view. Defines method to call for the buttons "planBtn", "applyBtn" and "destroyBtn".*/

function showDeployment() {
    $("#deploymentOutput").hide()
    $("#deployment-view").fadeIn("slow")
    $("#planBtn").click(function() {
        deploy("plan")
    })
    $("#applyBtn").click(function() {
        deploy("apply")
    })
    $("#destroyBtn").click(function() {
        deploy("destroy")
    })
}


/* Post-request to server containing the command triggered by either "planBtn", "applyBtn" or "destroyBtn"*/

function deploy(command) {
    if (moduleSource == "LIB") {
        $.post({
            url: modules + "/copy?" + "user=" + user + "&module=" + module,
            success: function(result) {
                console.log(result)
            }
        })
        importUserModules(users, user)
    }

    var url = deployment + "?user=" + user + "&module=" + module + "&command=" + command
    $.post({
        url: url,
        data: getParameters(),
        success: function(result) {
            console.log(result)
            showOutput()
        },
        dataType: "json"
    })
    $("#deploymentOutput").fadeIn("slow")
}

/*Shows input from server*/
function showOutput() {
    var url = deployment + "?user=" + user + "&module=" + module
    $.getJSON(url, function(result) {
        $("#deploymentOutput").html(result.output)
        if (result.status == "Running") {
            setTimeout(showOutput, 1000)
        }

    });
}

/*wraps entered variables in a JSON-file*/
function getParameters() {
    var table = document.getElementById("variablesTable");
    var variables = []
    for (var i = 1, row; row = table.rows[i]; i++) {
        variables.push({
            name: row.cells[0].childNodes[0].innerHTML,
            value: row.cells[1].childNodes[0].value
        })
}
    console.log(variables)
	  return JSON.stringify(variables)
};
