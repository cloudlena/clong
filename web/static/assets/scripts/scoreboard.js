"use strict";

const tbody = document.querySelector("#scoreboard tbody");

function init() {
  const ws = new WebSocket(wsURL("/ws/scoreboard"));

  // The server sends the best scores, best first, on connect and whenever they change
  ws.onmessage = (e) => drawScores(JSON.parse(e.data));

  ws.onclose = () => {
    showMessage("Reconnecting...");
    setTimeout(init, reconnectDelay);
  };
}

function drawScores(scores) {
  if (scores.length === 0) {
    showMessage("No scores yet...");
    return;
  }

  tbody.replaceChildren(...scores.map(scoreRow));
}

// Create a table row for a score
function scoreRow(score, i) {
  const rank = cell(i + 1 + ".");
  rank.className = "rank";

  const dot = document.createElement("span");
  dot.className = "dot";
  dot.style.backgroundColor = score.color;
  const name = cell(score.playerName);
  name.prepend(dot);

  const points = cell(score.finalScore);
  points.className = "score";

  const row = document.createElement("tr");
  row.append(rank, name, points);
  return row;
}

// Show a message instead of the scores
function showMessage(text) {
  const msg = cell(text);
  msg.colSpan = 3;
  msg.className = "msg";

  const row = document.createElement("tr");
  row.append(msg);
  tbody.replaceChildren(row);
}

// Create a table cell containing text
function cell(text) {
  const td = document.createElement("td");
  td.textContent = text;
  return td;
}

drawScores([]);
init();
