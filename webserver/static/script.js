$( document ).ready( function() {
    console.log('test')
    $.getJSON( "http://127.0.0.1:8080/api/", function( result ) {
        console.log(result)
    });
})
