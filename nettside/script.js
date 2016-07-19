//The current user operating the website
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
    $("#usernameInput").focus();
    $("#loginBtn").click(logIn)
    $("#library-view").hide()
    $("#variables-view").hide()
    $("#deployment-view").hide()
    $("#description").hide()
    //document.getElementsByClassName('alert-box output')[0].style.visibility = 'hidden';

    //Trykker på "logg inn" knappen hvis enter trykkes på i inputbox
    $('#usernameInput').keypress(function (e) {
        if (e.keyCode == 13)
            $('#loginBtn').click();
    });
})


function logIn() {
    user = $("#usernameInput").val()
	  $.get({
		    url: users,
		    success: function(result) {
			      var e = false;
			      for (i=0; i< result.length ; i++) {
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
}

function makeUser(user) {
	console.log("make new user")
}

function importLibraryModules(path) {
    $.getJSON(path, function (resultModules) {
        var content = ""
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

    });
}

function importUserModules(path, user) {
    $.getJSON(path + "?user=" + user, function (resultModules) {
        console.log(resultModules)
        console.log(resultModules);
        var content = ""
        content += "<option selected disabled hidden>Brukermoduler...</option>"
        for (i = 0; i < resultModules.length; i++) {
            content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].name + "</option>"
        }

        $("#existing").html(content);
        $("#existing").change(function() {
            module = $("#existing option:selected").text()
            moduleSource = "USER"
            showModule(path + "?user=" + user +  "&module=" + module)
            $("#variables-view").fadeIn("slow")
        })
    });
}


function showModule(path) {
    $.getJSON(path, function(result) {
        var content = ""
        $("#moduleName").html(result.name)
        $("#moduleDescription").html(result.description)
        var myTable = ""
				myTable += "<thead><tr><th>Navn</th><th>Verdi</th></tr></thead>"
        for (i = 0; i < result.variables.length; i++) {
            console.log(result.variables[i])
            var value = result.variables[i].value
            if (value == "") {
                value = result.variables[i].defaultValue
            }

            var textInputBox = '<input type="text" class="form-control" value="' + result.variables[i].defaultValue + '" id="' [i] + '" />';
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

function showOutput() {
    var url = deployment + "?user=" + user + "&module=" + module
    $.getJSON(url, function(result) {
        $("#deploymentOutput").html(result.output)
        if (result.status == "Running") {
            setTimeout(showOutput, 1000)
        }

    });
}

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
