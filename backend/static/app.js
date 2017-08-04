// Handling shortened links
$("#form").submit(function (event) {
    event.preventDefault();
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