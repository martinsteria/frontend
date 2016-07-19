var user
var module
var modules = "http://52.169.232.92/api/library"
var users = "http://52.169.232.92/api/users"
var deployment = "http://52.169.232.92/api/deploy"

$(document).ready(function () {
    $.ajaxSetup({ cache: false })
    $("#usernameInput").focus();
    $("#loginBtn").click(logIn)
    $("#library-view").hide()
    $("#variables-view").hide()
    $("#deployment-view").hide()
    $("#description").hide()
    $("#newUser").hide()

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
			      console.log(result)
			      var e = false;
			      for (i=0; i< result.length ; i++) {
				        console.log(result[i]+"="+user);
				        if (user == result[i]){
					          e = true;
				        }
			      }
			      if (e == true) {
			          importLibraryModules(users + "?user=" + user, "existing", "&");
			      }
			      else {
				        makeUser(user);
			      }
		    }
	  })
    $("#login-view").fadeOut("slow", function() {
        importLibraryModules(modules, "library", "?")
        $("#library-view").fadeIn("slow")
    })
    var content = "<span>" + "Du er logget in som " + user + "" + "</span>";
    $("#bruker").html(content)
    $("#bruker").show()
}

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

/*
x er enten lik '&' eller '?'
users + "?user=" + user + "&module=" + module
users + "?user=" + user + "?module=" + module
*/
function importLibraryModules(path, meny, x) {
    $.getJSON(path, function (resultModules) {
        console.log(resultModules);
        var content = ""
        if (resultModules != null){
        content += "<option selected disabled hidden>Biblioteksmoduler...</option>"
        for (i = 0; i < resultModules.length; i++) {
            content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].name + "</option>"
        }
        $("#"+meny).html(content);
        $("#"+meny).change(function() {
            module = $("#"+meny+" option:selected").text()
            showModule(path + x +"module=" + module)
            $("#variables-view").fadeIn("slow")
        })
    }
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
            var textInputBox = '<input type="text" class="form-control" value="' + result.variables[i].defaultValue + '" id="' + [i] + '" />';
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
    $.post({
        url: modules + "/copy?" + "user=" + user + "&module=" + module,
        success: function(result) {
            console.log(result)
        }
    })
    $("#deploymentOutput").hide()
    $("#deployment-view").show()
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
    var url = deployment + "?user=" + user + "&module=" + module + "&command=" + command
    $.post({
        url: url,
        data: getParameters,
        success: function(result) {
            console.log(result)
            showOutput()
        },
        dataType: "json"
    })
    $("#deploymentOutput").fadeIn("slow")
}

function showOutput() {
    var url = deployment + "?user=" + user + "&module=" + module + "&command=plan"
    $.getJSON(url, function(result) {
        $("#deploymentOutput").html(result.output)
        if (result.status == "Running") {
            setTimeout(showOutput, 1000)
        }

    });
}

function getParameters() {
	  var myRows = [];
	  var $headers = $("th");
	  var $rows = $("tbody tr").each(function(index) {
	      $cells = $(this).find("td");
	      myRows[index] = {};
	      $cells.each(function(cellIndex) {
		        if (cellIndex ==0){
		            myRows[index]['name'] = $(this).html();
		        }
		        if (cellIndex ==2){
		            var value = document.getElementById(index).value;
		            myRows[index]['value'] = value;
		        }
	      });
	  });
	  var myObj = myRows;
	  return JSON.stringify(myObj)
};
