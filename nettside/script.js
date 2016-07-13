document.getElementsByClassName('col-sm-5')[0].style.visibility = 'hidden';
document.getElementsByClassName('col-sm-4')[0].style.visibility = 'hidden';
function velgMal() {
	var e = document.getElementById("Select2");
	var templateName = e.options[e.selectedIndex].text;
	document.getElementById("if").innerHTML = templateName;
	var templateId = e.options[e.selectedIndex].value;
	if (templateId == "1") {
		Milestone1();
	} else if (templateId =="2") {
			Milestone2();
	} else {
		document.getElementById("if").innerHTML = "Mal finnes ikke";
	}
    document.getElementsByClassName('col-sm-5')[0].style.visibility = 'visible';
	
};

function lagFil() {
    document.getElementsByClassName('col-sm-4')[0].style.visibility = 'visible';
};

function lagArray() {
    var tabellVar11 = document.getElementById("innVarNavn1").innerHTML;
    var tabellVar12 = document.getElementById("Var1Verdi").value;
 
    var tabellVar21 = document.getElementById("innVarNavn2").innerHTML;
    var tabellVar22 = document.getElementById("Var2Verdi").value;

    var $verdier = $('#innVarNavn1, #Var1Verdi ')
    //var $verdier = $('#tabellVar11, #tabellVar12')
    var obj = {values: [
        {
            name: tabellVar11,
            value: tabellVar12
        },
        {
            name: tabellVar21,
            value: tabellVar22
        }
    ]
    };
    console.log(obj)

    var json = JSON.stringify(obj);

    //alert(tabellVar11 + " = '" + tabellVar12 + "'\n" + tabellVar21 + " = '" + tabellVar22 + "'");
    alert(json);

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
    $.ajaxSetup({ cache: false })
 //   $("#Button1").click(function () {
        $.getJSON('http://tfbrowser.routable.org/api/modules/milestone-1', function (result) {
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
   // });
});
}

function Milestone2(){
			$.ajaxSetup({cache: false})
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
}

function LastOpp(){
	var x = document.getElementById("myFile");
    var txt = "";
	if ('files' in x) {
	//filen heter nå "file"
	var file = x.files[0];
	if ('name' in file) {
		txt += "name: " + file.name + "<br>";
		}
	if ('size' in file) {
		txt += "size: " + file.size + " bytes <br>";
		}
	}
    else {
        if (x.value == "") {
            txt += "Select a file";
        } else {
            txt += "The files property is not supported by your browser!";
            txt  += "<br>The path of the selected file: " + x.value; // If the browser does not support the files property, it will return the path of the selected file instead.
        }
    } 
    document.getElementById("Fil").innerHTML = txt;
}
