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
    $("#login-view").hide()
    $("#library-view").show()
    importLibraryModules(modules, "library", "?")
}

function makeUser(user) {
    console.log("make new user")
    $.post({
        url: users + "?user=" + user,
        success: function (result) {
            console.log("new User: "+ user)
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
            for (i = 0; i < resultModules.length; i++) {
                content += "<option value=\"" + i + "\" id=\"" + resultModules[i].id + "\" >" + resultModules[i].name + "</option>"
            }

            $("#"+meny).html(content);
            $("#"+meny).click(function() {
            $("#variables-view").show()
            module = $("#"+meny+" option:selected").text()
            showModule(path + x +"module=" + module)
        })
    }
    });
}


function showModule(path) {
    $.getJSON(path, function(result) {
        var content = ""
		    content += "<span>" + result.name + "</span><br>" + result.description;
        $("#description").html(content)
        $("#description").show()
        console.log(result)
        var myTable = ""
				myTable += "<thead><tr><th>Navn</th><th></th><th>Verdi</th></tr></thead>"
        for (i = 0; i < result.variables.length; i++) {
            var textInputBox = '<input type="text" value="' + result.variables[i].defaultValue + '" id="' + [i] + '"name="' + [i] + '" />';
            myTable += '<tr>'
            myTable += '<td>' + result.variables[i].name + '</td>'
            myTable += '<td> <div class="help-tip"> <p>'+ result.variables[i].description +'</p> </div> </td>'
            myTable += '<td>' + textInputBox + '</td>'
            myTable += '</tr>'
        }
        $("#variablesTable").html(myTable)

        $("#showDeployment").click(showDeployment)
    });
}

function showDeployment() {
    $.post({
        url: modules + "/copy?" + "user=" + user + "&module=" + module,
        success: function(result) {
            console.log(result)
        }
    })
    $("#deployment-view").show()
    $("#planBtn").click(plan)
    $("#applyBtn").click(apply)
    $("#destroyBtn").click(destroy)
}

function plan() {
    var url = deployment + "?user=" + user + "&module=" + module + "&command=plan"
    $.post({
        url: url,
        data: getParameters,
        success: function(result) {
            console.log(result)
            setInterval(showOutput, 1000)
        },
        dataType: "json"
    })
}

function showOutput() {
    var url = deployment + "?user=" + user + "&module=" + module + "&command=plan"
    $.getJSON(url, function(result) {
        console.log(result)
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
