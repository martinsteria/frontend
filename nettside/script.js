

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

$(document).ready(function () {
    $.ajaxSetup({cache: false})
    $("#leseFil").click(function () {
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
    });
});

