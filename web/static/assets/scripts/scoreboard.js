"use strict";

// Set game preferences
var maxScores = 10;

// Initialize global variables
var ws = null;
var highScores = [];
var tbody = document.querySelector("#scoreboard tbody");

function init() {
  // Initialize WebSocket connection
  ws = new WebSocket(wsURL("/ws/screen"));

  ws.onopen = function () {
    // Get existing scores
    fetch("/api/scores")
      .then(function (res) {
        return res.json();
      })
      .then(function (data) {
        highScores = data;
        drawScores();
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
    showMessage("Reconnecting...");
    setTimeout(init, 3000);
  };
}

function drawScores() {
  if (highScores.length === 0) {
    showMessage("No scores yet...");
    return;
  }

  // Keep only the best scores
  highScores.sort(function (a, b) {
    return b.finalScore - a.finalScore;
  });
  highScores = highScores.slice(0, maxScores);

  tbody.replaceChildren(...highScores.map(scoreRow));
}

// Create a table row for a score
function scoreRow(score, i) {
  var rank = cell(i + 1 + ".");
  rank.className = "rank";

  var dot = document.createElement("span");
  dot.textContent = " ●";
  dot.style.color = score.color;
  var name = cell(score.player.name);
  name.append(dot);

  var row = document.createElement("tr");
  row.append(rank, name, cell(score.finalScore));
  return row;
}

// Show a message instead of the scores
function showMessage(text) {
  var msg = cell(text);
  msg.colSpan = 3;
  msg.style.textAlign = "center";

  var row = document.createElement("tr");
  row.append(msg);
  tbody.replaceChildren(row);
}

// Create a table cell containing text
function cell(text) {
  var td = document.createElement("td");
  td.textContent = text;
  return td;
}

init();
