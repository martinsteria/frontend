$( document ).ready( function() {
    console.log('test')
    $.getJSON( "http://127.0.0.1:8080/api/", function( result ) {
        console.log(result);
        content = "";
        for ( i = 0; i < result.people.length; i++ ) {
            content += "<tr>";
            content += "<td>" + result.people[i].name + "</td><td>" + result.people[i].catchPhrase + "</td>";
            content += "</tr>"
        }
        $( "#tab" ).html(content)
    });
})
