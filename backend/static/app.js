// Handling shortened links
$("#shorten").click(function () {
    var path = "http://" + window.location.hostname + ":" + window.location.port;
    var url = $("#urlField").val();
    $.ajax({
        type: 'post',
        dataType: 'json',
        data: JSON.stringify({
            "url": url
        }),
        url: path + '/shorten',
        success: function (msg) {
            var shortenedUrl = path + '/s/' + msg.tag
            $("#tag").text("Your shortened url: " + shortenedUrl);
            $("#tag").attr("href", shortenedUrl);

        }
    })
})

// Listen for user input in URLField and map it to the tag
$("#urlField").bind("input", function () {
    var str = $("#urlField").val()
    $("#tag").text(str);
});