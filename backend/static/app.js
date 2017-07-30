
// Handling shortened links
$("#shorten").click(function () {
    var path = "http://" + window.location.hostname + ":" + window.location.port + "/shorten";
    var url = String($("#urlField").val());
    $.ajax({
        type: 'post',
        dataType: 'json',
        data: JSON.stringify({
            "url": url
        }),
        url: path,
        success: function(msg, data) {
            console.log(data, msg)
        }
    })
})