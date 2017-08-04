// Handling shortened links

$("#form").submit(function(event) {
  event.preventDefault();
  var path = "http://" + window.location.hostname + ":" + window.location.port;
  var url = $("#urlField").val();
  $.ajax({
    type: "post",
    dataType: "json",
    data: JSON.stringify({
      url: url
    }),
    url: path + "/shorten",
    success: function(msg) {
      var shortenedUrl = path + "/s/" + msg.tag;
      var analytics = displayAnalytics(path, msg.tag, url);
      if (analytics.error) {
        $("#tag").text("Sorry an error occured");
      } else {
        // Click number
        $("#tag").text("Your shortened url: " + shortenedUrl);
        $("#tag").attr("href", shortenedUrl);
      }
    }
  });
});

// Display analytics
function displayAnalytics(path, tag, url) {
  cN = document.getElementById("clickNumber").getContext("2d");
  cT = document.getElementById("clickTimes").getContext("2d");
  $.ajax({
    url: path + "/analytics/" + tag,
    dataType: "json",
    success: function(data) {
      var cNChart = new Chart(cN, {
        type: "bar",
        data: {
          labels: [url],
          datasets: [
            {
              label: "Number of clicks",
              data: [data.clickNumber],
              backgroundColor: ["rgba(75, 192, 192, 0.2)"]
            }
          ]
        },
        options: {
          maintainAspectRatio: false
        }
      });

      var cTChart = new Chart(cT, {
        type: "line",
        data: {
          datasets: [
            {
              label: "Clicks over time",
              data: data.clickTimes,
              fill: false
            }
          ]
        },
        options: {
          maintainAspectRatio: false
        }
      });

      return analytics;
    }
  });
}
