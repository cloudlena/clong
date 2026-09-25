"use strict";

// Initialize global variables
var ws = null;
var main = document.getElementById("main");
var spinner = document.getElementById("spinner");
var intro = document.getElementById("intro");
var reconnectingMsg = document.getElementById("reconnecting-msg");
var startButton = document.getElementById("start-button");
var secNum = document.getElementById("sec-num");
var scoreNum = document.getElementById("score-num");
var mc = new Hammer(main);
var w = 0;
var h = 0;
var myID = ensureCookie("userid", uuid);
var myColor = randomColor();
var gameRunning = false;
var locked = false;
var myPoints = 0;
var myTime = 0;

// Init player
ensureCookie("username", askUserName);
main.style.backgroundColor = myColor;
spinner.style.color = myColor;

// Listen for swipes to shoot balls
mc.get("swipe").set({ direction: Hammer.DIRECTION_ALL });
mc.on("swipe", function (e) {
  if (!gameRunning || locked) {
    return;
  }
  intro.hidden = true;

  var msg = {
    type: "BALL_INIT",
    color: myColor,
    posX: relW(e.center.x),
    velocityX: relW(e.velocityX),
    velocityY: relH(e.velocityY),
  };
  ws.send(JSON.stringify(msg));

  lock();
});

function init() {
  calcScreenSize();
  window.onresize = calcScreenSize;

  // Initialize WebSocket connection
  ws = new WebSocket(wsURL("/ws/controller"));

  ws.onopen = function () {
    unlock();
    reconnectingMsg.hidden = true;
  };

  // Listen for when ball is done to unlock screen
  ws.onmessage = function (e) {
    var msg = JSON.parse(e.data);
    if (gameRunning && msg.type === "BALL_DONE" && msg.player.id === myID) {
      myPoints += msg.points;
      scoreNum.textContent = myPoints;
      unlock();
    }
  };

  // Try to reconnect on close
  ws.onclose = function () {
    lock();
    reconnectingMsg.hidden = false;
    setTimeout(init, 3000);
  };
}

function startGame() {
  startButton.hidden = true;
  myTime = 60;
  secNum.textContent = myTime;
  myPoints = 0;
  scoreNum.textContent = myPoints;
  gameRunning = true;
  var gameLoop = setInterval(function () {
    myTime--;
    secNum.textContent = myTime;
    if (myTime <= 0) {
      gameRunning = false;
      unlock();
      startButton.hidden = false;
      clearInterval(gameLoop);
      var msg = {
        type: "GAME_FINISHED",
        finalScore: myPoints,
        color: myColor,
      };
      ws.send(JSON.stringify(msg));
    }
  }, 1000);
}

// Gets a cookie from the cookie jar
function getCookie(name) {
  var prefix = name + "=";
  var cookie = decodeURIComponent(document.cookie)
    .split("; ")
    .find(function (c) {
      return c.startsWith(prefix);
    });
  return cookie ? cookie.substring(prefix.length) : null;
}

// Returns the value of a cookie, creating it first if it doesn't exist yet
function ensureCookie(name, createValue) {
  var value = getCookie(name);
  if (value === null) {
    value = createValue();
    document.cookie = name + "=" + value;
  }
  return value;
}

// Asks the user for their name until they enter a valid one
function askUserName() {
  var userName = null;
  while (userName === null || userName.length < 2 || userName.length > 30) {
    userName = prompt("What's your full name?", "");
  }
  return userName;
}

function lock() {
  locked = true;
  main.style.display = "none";
}

function unlock() {
  locked = false;
  main.style.display = "block";
}

init();
