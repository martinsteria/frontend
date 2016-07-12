document.getElementsByClassName('col-sm-3')[0].style.visibility = 'hidden';
document.getElementsByClassName('col-sm-2')[0].style.visibility = 'hidden';
function velgMal() { 
    document.getElementById("Mal").innerHTML = "Bolle!";
    var e = document.getElementById("Select2");
	var templateName = e.options[e.selectedIndex].text;
	document.getElementById("Mal").innerHTML = templateName;
	var templateId = e.options[e.selectedIndex].value;
	document.getElementById("if").innerHTML = templateId;
	if (templateId == "1") {
		Milestone1();
		document.getElementById("if").innerHTML = "Bolle!";
	} else if (templateId =="2") {
		    varName1 = "Variabel for milestone 2";
            $("#innVarNavn1").html(varName1);
            //$("#leseText").html(varName1);
            varBesk1 = "Variabel for milestone 2";
            $("#innBesk1").html(varBesk1);
            varName2 = "Variabel for milestone 2";
            $("#innVarNavn2").html(varName2);
            varBesk2 = "Variabel for milestone 2";
            $("#innBesk2").html(varBesk2);
			varName3 = "Variabel for milestone 2";
            $("#innVarNavn3").html(varName3);
            varBesk3 = "Variabel for milestone 2";
            $("#innBesk3").html(varBesk3);
            varName4 = "Variabel for milestone 2";
            $("#innVarNavn4").html(varName4);
            varBesk4 = "Variabel for milestone 2";
            $("#innBesk4").html(varBesk4);
		document.getElementById("if").innerHTML = "Ikke Bolle!";
	} else {
		document.getElementById("if").innerHTML = "Mal finnes ikke";
	}	
    document.getElementsByClassName('col-sm-3')[0].style.visibility = 'visible';
}

(function () {
    var textFile = null,
        makeTextFile = function (text) {
            var data = new Blob([text], { type: 'text/plain' });

            // If we are replacing a previously generated file we need to
            // manually revoke the object URL to avoid memory leaks.
            if (textFile !== null) {
                window.URL.revokeObjectURL(textFile);
            }

            textFile = window.URL.createObjectURL(data);

            return textFile;
        };


    var create = document.getElementById('create'),
        textbox = document.getElementById('leseText');

    create.addEventListener('click', function () {
        var link = document.getElementById('downloadlink');
        link.href = makeTextFile(textbox.value);
        link.style.display = 'block';
    }, false);
})();


function Milestone1() {
$(document).ready(function () {
    $.ajaxSetup({cache: false})
   /* $("#Button1").click(function () { */
        $.getJSON('http://tfbrowser.routable.org/api', function (result) {
            console.log(result);
            //document.write(result.name); skriver første til hovedside etter "last inn"
            // $("#leseText").html(result.name); skriver første navnet til textarea etter last inn
            
            varAll = result;
            $("#leseText").html(varAll);
            
            varName1 = result.variables[0].name;
            $("#innVarNavn1").html(varName1);
            //$("#leseText").html(varName1);
            varBesk1 = result.variables[0].description;
            $("#innBesk1").html(varBesk1);

            varName2 = result.variables[1].name;
            $("#innVarNavn2").html(varName2);
            varBesk2 = result.variables[1].description;
            $("#innBesk2").html(varBesk2);
			
			varName2 = result.variables[2].name;
            $("#innVarNavn3").html(varName2);
            varBesk2 = result.variables[2].description;
            $("#innBesk3").html(varBesk2);
			
			varName2 = result.variables[3].name;
            $("#innVarNavn4").html(varName2);
            varBesk2 = result.variables[3].description;
            $("#innBesk4").html(varBesk2);
            
            //content = "";
            // content2 = "";
            // content += result.name + "\n";
            // content += result.description + "\n";
            // content2 += result.variables[0].name;
            // content2 += result.variables[0].description;
            //content2 = result;
            //$("#leseText").html(content);
               // $.each(result, function (i, field) {
              //      $("#leseText").append(field + " ");
              
            //});
        });
    /*}); */
});

}

