//The current user operating the website
var user

//The current module the user is editing
var module

//Is the module from the library or from the users collection
//VALS: "USER" or "LIB"
var moduleSource

//API endpoints
var apiRoot = "http://localhost/api"
var libModulesAPIEndpoint = apiRoot + "/library"
var usersAPIEndpoint = apiRoot + "/users"
var deploymentAPIEndpoint = apiRoot + "/deploy"

$(document).ready(function () {
    $("#copyright").html("&copy; Sopra Steria " + new Date().getFullYear())
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
If the user exist it calls importUserModules(), if not it calls createNewUser()*/

function logIn() {
    user = $("#usernameInput").val()
	  $.get({
		    url: usersAPIEndpoint, //Get-request to server for url= ../api/users
		    success: function(result) {
			      var e = false;
            if (result == null) {
                createNewUser(user)
                return
            }
			      for (i=0; i< result.length ; i++) { // for all usernames check if it equals the enteres username
				        if (user == result[i]){
					          e = true;
				        }
			      }
			      if (e == true) {
			          importUserModules(usersAPIEndpoint, user);
			      }
			      else {
				        createNewUser(user);
			      }
		    }
	  })
    $("#login-view").fadeOut("slow", function() {
        importLibraryModules(libModulesAPIEndpoint)
        $("#library-view").fadeIn("slow")
    })
    var content = "<span>" + "Du er logget inn som " + user + "" + "</span>";
    $("#bruker").html(content)
    $("#bruker").show()
}

/* sends a post request to the server with the entered username*/
function createNewUser(user) {
    $.post({
        url: usersAPIEndpoint + "?user=" + user,
        success: function (result) {
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
            $("#library").unbind()
            $("#library").change(function() {
                module = $("#library option:selected").attr("id")
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
            module = $("#userLibrary option:selected").attr("id")
            moduleSource = "USER"
            showModule(path + "?user=" + user +  "&module=" + module)
            $("#variables-view").fadeIn("slow")
        })
    });
}


/* Sends a get request for a JSON-file containing the variables for the selected module. Adds each variables to the table #moduleVariables*/
function showModule(path) {
    console.log("Show module called")
    $.getJSON(path, function(result) {
        var content = "2. "
        var name = result.name
        $("#moduleName").html(result.name)
        $("#moduleProvider").html("<b>Provider: </b>" + result.provider)
        $("#moduleDescription").html(result.description)

        //Fill outputs table
        var content = ""
        content += "<thead><tr><th>Output</th><th>Beskrivelse</th></tr></thead>"
        for (i = 0; i < result.outputs.length; i++) {
            content += "<tr>"
            content += "<td>" + result.outputs[i].name + "</td>"
            content += "<td>" + result.outputs[i].description + "</td>"
            content += "</tr>"
        }
        $("#moduleOutputs").html(content)

        //Fill variables table
        var content = ""
				content += "<thead><tr><th>Variabel</th><th>Verdi</th></tr></thead>"
        for (i = 0; i < result.variables.length; i++) {
            var value = result.variables[i].value
            if (value == "") {
                value = result.variables[i].defaultValue
            }

            var textInputBox = '<input type="text" class="form-control" value="' + value + '" />';
            content += '<tr>'
            content += '<td><a href="#" data-placement="left" data-toggle="tooltip" title="' + result.variables[i].description + '">' + result.variables[i].name + '</a></td>'
            content += '<td>' + textInputBox + '</td>'
            content += '</tr>'
        }
        $("#moduleVariables").html(content)
        $("#showDeployment").unbind()
        $("#showDeployment").click(showDeployment)
        $("[data-toggle=\"tooltip\"]").tooltip();
    });
}

/* Shows deployment-view. Defines method to call for the buttons "planBtn", "applyBtn" and "destroyBtn".*/

function showDeployment() {
    console.log("show deployment called")
    $("#deployment-view").fadeIn("slow")

    $("#planBtn").unbind()
    $("#planBtn").click(function() {
        deploy("plan")
    })

    $("#applyBtn").unbind()
    $("#applyBtn").click(function() {
        deploy("apply")
    })

    $("#destroyBtn").unbind()
    $("#destroyBtn").click(function() {
        deploy("destroy")
    })
}


/* Post-request to server containing the command triggered by either "planBtn", "applyBtn" or "destroyBtn"*/

function deploy(command) {
    console.log("deploy called")
    if (moduleSource == "LIB") {
        $.post({
            url: libModulesAPIEndpoint + "/copy?" + "user=" + user + "&module=" + module,
            success: function(result) {
                importUserModules(usersAPIEndpoint, user)
            }
        })
    }

    var url = deploymentAPIEndpoint + "?user=" + user + "&module=" + module + "&command=" + command
    $.post({
        url: url,
        data: getParameters(),
        success: function(result) {
            showOutput()
            $("#edit-view").fadeOut("slow", function() {
                $("#output-view").fadeIn("slow")
            })
        },
        dataType: "json"
    })
}

/*Shows input from server*/
function showOutput() {
    var url = deploymentAPIEndpoint + "?user=" + user + "&module=" + module
    $.getJSON(url, function(result) {
        $("#deploymentOutput").html(result.output)
        if (result.status == "Running") {
            setTimeout(showOutput, 1000)
        }

    });
}

/*wraps entered variables in a JSON-file*/
function getParameters() {
    var table = document.getElementById("moduleVariables");
    var variables = []
    for (var i = 1, row; row = table.rows[i]; i++) {
        variables.push({
            name: row.cells[0].childNodes[0].innerHTML,
            value: row.cells[1].childNodes[0].value
        })
    }
	  return JSON.stringify(variables)
}
