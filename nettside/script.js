//API endpoints
var apiRoot = "http://10.118.200.140:8080/api"
var LIBRARY_ENDPOINT = apiRoot + "/library"
var USERS_ENDPOINT = apiRoot + "/users"
var DEPLOY_ENDPOINT = apiRoot + "/deploy"

$(document).ready(function () {
    var user
    var module
    var moduleSource

    $("#copyright").html("&copy; Sopra Steria " + new Date().getFullYear())
    $.ajaxSetup({ cache: false })
    $("#login-view").hide()
    $("#login-view").fadeIn("slow")
    $("#usernameInput").focus()
    $("#loginBtn").click(function () {
        user = $("#usernameInput").val()
        logIn(user,
              function() {
                  getUserModules(user)
              },
              function() {
                  createNewUser(user)
              },
              function() {
                  transition("#login-view", "#library-view")
                  var content = "<span>" + "Du er logget inn som " + user + "</span>"
                  $("#bruker").html(content)
                  $("#bruker").show()
                  getLibraryModules()
              })
    })

    $("#library").change(function() {
        module = $("#library option:selected").attr("id")
        moduleSource = "LIB"
        getModule(LIBRARY_ENDPOINT +"?module=" + module)
        show("#variables-view")
    })

    $("#userLibrary").change(function () {
        module = $("#userLibrary option:selected").attr("id")
        moduleSource = "USER"
        getModule(USERS_ENDPOINT + "?user=" + user +  "&module=" + module)
        show("#variables-view")
    })

    $("#showDeployment").click(function() {
        if (moduleSource == "LIB") {
            copyModule(user, module)
        }
        show("#deployment-view")
    })

    $("#planBtn").click(function() {
        deploy(user, module, "plan")
        transition("#edit-view", "#output-view")
    })

    $("#applyBtn").click(function() {
        deploy(user, module, "apply")
        transition("#edit-view", "#output-view")
    })

    $("#destroyBtn").click(function() {
        deploy(user, module, "destroy")
        transition("#edit-view", "#output-view")
    })

    $("#backButton").click(function() {
        transition("#output-view", "#edit-view")
    })

    $("#library-view").hide()
    $("#variables-view").hide()
    $("#deployment-view").hide()
    $("#description").hide()
    $("#newUser").hide()
    $("#output-view").hide()

    //"logg inn" knapp aktiveres ved å trykke enter i inputbox
    $('#usernameInput').keypress(function (e) {
        if (e.keyCode == 13)
            $('#loginBtn').click()
    })
})

function transition(from, to) {
    $(from).fadeOut("slow", function() {
        $(to).fadeIn("slow")
    })
}

function show(view) {
    $(view).fadeIn("slow")
}

/**
 * Logs the user in with the provided username.
 * If the user does not exist on the server, a new user is created.
 * Fades out the login-view and fades in the library-view
 * Fetches all library modules as well as the users library modules
 * @param {string} user - The name of the user to log in
 * @param {function()} success - Function to call if user exists
 * @param {function()} failure - Function to call if user does not exist
 * @param {function()} regardless - Function to call regardless of user existing or not
 */
function logIn(user, success, failure, regardless) {
	  $.get({
		    url: USERS_ENDPOINT,
		    success: function(result) {
            if (result == null) {
                failure()
                regardless()
                return
            }

			      var exists = false
			      for (i=0; i< result.length ; i++) {
				        if (user == result[i]){
					          exists = true
                    break
				        }
			      }
			      if (exists == true) {
                success()
                regardless()
                return
			      }
            failure()
            regardless()
		    }
    })
}

/**
 * Creates a new user on the server
 * @param {string} user - The name of the user to create
 */
function createNewUser(user) {
    $.post({
        url: USERS_ENDPOINT + "?user=" + user,
        success: function (result) {
            var content= "<span>" + "Ny Bruker" + "</span><br>" + "Ny bruker opprettet for "+ user
            $("#newUser").html(content)
            $("#newUser").show()
        }
    })
}

/**
 * Fetches all library modules from server and displays them in the #library dropdown.
 * Also binds changes to the dropdown to show #variables-view
 */
function getLibraryModules() {
    $.getJSON(LIBRARY_ENDPOINT, function (resultModules) {
      var content = ""
        if (resultModules != null){
            content += "<option selected disabled hidden>Biblioteksmoduler...</option>"
            for (i = 0; i < resultModules.length; i++) {
                if (resultModules[i].provider != "") {
                    content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].provider + ": " + resultModules[i].name + "</option>"
                }
            }
            $("#library").html(content)
        }
    })
}

/**
 * Fetches all the user's modules from server and displays them in the #userLibrary dropdown.
 * Also binds changes to the dropdown to show #variables-view
 * @param {string} user - The user to fetch modules for
 */
function getUserModules(user) {
  $.getJSON(USERS_ENDPOINT + "?user=" + user, function (resultModules) {
      var content = ""
        content += "<option selected disabled hidden>Brukermoduler...</option>"
        if (resultModules == null) {
            return
        }
        for (i = 0; i < resultModules.length; i++) {
            content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].provider + ": " + resultModules[i].name + "</option>"
        }
        $("#userLibrary").html(content)
    })
}


/**
 * Fetches a module's documentation from the server and fills #variables-view with the information
 * @param {string} path - The API path for a module
 */
function getModule(path) {
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
                value = result.variables[i].default
            }

            var textInputBox = '<input type="text" class="form-control" value="' + value + '" />'
            content += '<tr>'
            content += '<td><a href="#" data-placement="left" data-toggle="tooltip" title="' + result.variables[i].description + '">' + result.variables[i].name + '</a></td>'
            content += '<td>' + textInputBox + '</td>'
            content += '</tr>'
        }
        $("#moduleVariables").html(content)
        $("[data-toggle=\"tooltip\"]").tooltip()
    })
}

/**
 * Shows #deployment-view. Binds the deploy function to buttons plan, apply and destroy
 */
function showDeployment() {
  $("#deployment-view").fadeIn("slow")
}

/**
 * Requests the server to execute a terraform command.
 * Also shows #output-view and incrementally fills it with the output from the server
 * @param {string} user - The user that should execute the command
 * @param {string} module - The module to deploy
 * @param {string} command - The command to execute
 */
function deploy(user, module, command) {
    var url =
        $.post({
            url: DEPLOY_ENDPOINT + "?user=" + user + "&module=" + module + "&command=" + command,
            success: function(result) {
                getOutput(user, module)
            },
            dataType: "json"
        })
}

function copyModule(user, module) {
    $.post({
        url: LIBRARY_ENDPOINT + "/copy?" + "user=" + user + "&module=" + module,
        success: function(result) {
            getUserModules(user)
        }
    })
}

/**
 * Reads the current output of the previously ran command from the server and fills in #output-view.
 * Recursively calls itself until the command is finished.
 * @param {string} user - The user that ran the command
 * @param {string} module - The module that was deployed
 */
function getOutput(user, module) {
    var url = DEPLOY_ENDPOINT + "?user=" + user + "&module=" + module
    $.getJSON(url, function(result) {
      $("#deploymentOutput").html(result.output)
        $("#deploymentError").html(result.error)
        if (result.status == "Running") {
            setTimeout(
                function() {
                    getOutput(user, module)
                },
                1000)
        }
    })
}

/**
 * Returns the variables filled in by the user in the #moduleVariables table
 * @return {string} A stringified JSON representation of the table with fields "name" and "value"
 */
function parseParameters() {
    var table = document.getElementById("moduleVariables")
    var variables = []
    for (var i = 1, row; row = table.rows[i]; i++) {
        variables.push({
            name: row.cells[0].childNodes[0].innerHTML,
            value: row.cells[1].childNodes[0].value
        })
    }
	  return JSON.stringify(variables)
}
