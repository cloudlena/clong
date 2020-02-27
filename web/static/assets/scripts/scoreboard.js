"use strict";

// Set game preferences
var maxScores = 10;

// Initialize global variables
var ws = null;
var highScores = [];

function init() {
  // Initialize WebSocket connection
  ws = new WebSocket(wsProtocol() + "//" + window.location.host + "/ws/screen");

  ws.onopen = function () {
    // Get existing scores
    $.get("/api/scores", function (data) {
      if (data !== null) {
        highScores = data;
        drawScores();
      }
    });
  };

  // Listen for new scores coming in
  ws.onmessage = function (e) {
    var msg = JSON.parse(e.data);
    if (msg.type === "GAME_FINISHED") {
      highScores.push(msg);
      drawScores();
    }
  };

  // Try to reconnect on close
  ws.onclose = function () {
    $("#scoreboard table").html(
      '<tr><td style="text-align:center;">Reconnecting...</td></tr>'
    );
    setTimeout(init, 3000);
  };
}

function drawScores() {
  // Sort by final score
  highScores.sort(function (a, b) {
    return b.finalScore - a.finalScore;
  });

  // Extract top 10
  highScores = highScores.slice(0, maxScores);

  var htmlString = "";
  for (var i = 0; i < highScores.length; i++) {
    htmlString +=
      "<tr>" +
      '<td class="rank">' +
      (i + 1) +
      ".</td>" +
      "<td>" +
      highScores[i].player.name +
      '<span style="color: ' +
      highScores[i].color +
      ';"> &#9679;</span></td>' +
      "<td>" +
      highScores[i].finalScore +
      "</td>" +
      "</tr>";
  }

  $("#scoreboard table").html(htmlString);
}

init();
