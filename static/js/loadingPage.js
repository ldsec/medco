
$("#s").click(function(e) {
    e.preventDefault();
    $.ajax({
        url: "http://localhost:47066/medco-loader/lala"
    }).done(function (data) {
        // Do whatever with returned data
    });
});