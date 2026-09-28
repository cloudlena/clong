"use strict";

// Set game preferences
var gameSeconds = 60;
var ballTimeout = 15000; // ms after which a ball is considered lost
var swipeMinDistance = 10; // px
var swipeMinVelocity = 0.3; // px/ms
var swipeVelocityWindow = 50; // ms at the end of a swipe used to measure its velocity

// Initialize global variables
var ws = null;
var outbox = [];
var main = document.getElementById("main");
var spinner = document.getElementById("spinner");
var intro = document.getElementById("intro");
var reconnectingMsg = document.getElementById("reconnecting-msg");
var startButton = document.getElementById("start-button");
var secNum = document.getElementById("sec-num");
var scoreNum = document.getElementById("score-num");
var w = 0;
var h = 0;
var myID = ensureCookie("userid", function () {
  return crypto.randomUUID();
});
var myColor = randomColor();
var gameRunning = false;
var locked = false;
var lockTimeout = null;
var myPoints = 0;
var myTime = 0;
var pointerSamples = [];

// Init player
ensureCookie("username", askUserName);
main.style.backgroundColor = myColor;
spinner.style.color = myColor;

// Track pointer movements to detect swipes
main.addEventListener("pointerdown", function (e) {
  if (e.isPrimary) {
    main.setPointerCapture(e.pointerId);
    pointerSamples = [pointerSample(e)];
  }
});
main.addEventListener("pointermove", function (e) {
  if (e.isPrimary && pointerSamples.length > 0) {
    pointerSamples.push(pointerSample(e));
  }
});
main.addEventListener("pointercancel", function () {
  pointerSamples = [];
});
main.addEventListener("pointerup", function (e) {
  if (!e.isPrimary || pointerSamples.length === 0) {
    return;
  }
  var end = pointerSample(e);
  var start = pointerSamples[0];
  // Flicks accelerate, so velocity is measured at the end of the swipe
  var recent = pointerSamples.find(function (s) {
    return end.t - s.t <= swipeVelocityWindow;
  });
  pointerSamples = [];

  var distance = Math.hypot(end.x - start.x, end.y - start.y);
  var duration = end.t - start.t;
  if (
    !recent ||
    end.t === recent.t ||
    distance < swipeMinDistance ||
    distance / duration < swipeMinVelocity
  ) {
    return;
  }
  shoot(
    end.x,
    (end.x - recent.x) / (end.t - recent.t),
    (end.y - recent.y) / (end.t - recent.t),
  );
});

// Shoot a ball at the screen
function shoot(x, velocityX, velocityY) {
  if (!gameRunning || locked) {
    return;
  }
  intro.hidden = true;

  send({
    type: "BALL_INIT",
    color: myColor,
    posX: relW(x),
    velocityX: relW(velocityX),
    velocityY: relH(velocityY),
  });

  // Don't wait forever if the ball gets lost, e.g. because the screen disconnected
  lock();
  lockTimeout = setTimeout(unlock, ballTimeout);
}

function init() {
  calcScreenSize();
  window.onresize = calcScreenSize;

  // Initialize WebSocket connection
  ws = new WebSocket(wsURL("/ws/controller"));

  ws.onopen = function () {
    unlock();
    reconnectingMsg.hidden = true;
    flushOutbox();
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
  myTime = gameSeconds;
  secNum.textContent = myTime;
  myPoints = 0;
  scoreNum.textContent = myPoints;
  gameRunning = true;
  send({ type: "GAME_STARTED" });
  var gameLoop = setInterval(function () {
    myTime--;
    secNum.textContent = myTime;
    if (myTime <= 0) {
      gameRunning = false;
      unlock();
      startButton.hidden = false;
      clearInterval(gameLoop);
      send({ type: "GAME_FINISHED", color: myColor });
    }
  }, 1000);
}

// Send a message to the server, keeping it until the connection is open
function send(msg) {
  outbox.push(JSON.stringify(msg));
  flushOutbox();
}

function flushOutbox() {
  while (ws.readyState === WebSocket.OPEN && outbox.length > 0) {
    ws.send(outbox.shift());
  }
}

// Record the position of a pointer at the current time
function pointerSample(e) {
  return { x: e.clientX, y: e.clientY, t: e.timeStamp };
}

// Gets a cookie from the cookie jar
function getCookie(name) {
  var prefix = name + "=";
  var cookie = document.cookie.split("; ").find(function (c) {
    return c.startsWith(prefix);
  });
  if (!cookie) {
    return null;
  }
  var value = cookie.substring(prefix.length);
  try {
    return decodeURIComponent(value);
  } catch {
    // Cookies from older versions weren't encoded
    return value;
  }
}

// Returns the value of a cookie, creating it first if it doesn't exist yet
function ensureCookie(name, createValue) {
  var value = getCookie(name);
  if (value === null) {
    value = createValue();
    document.cookie = name + "=" + encodeURIComponent(value);
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
  clearTimeout(lockTimeout);
  locked = true;
  main.style.display = "none";
}

function unlock() {
  clearTimeout(lockTimeout);
  locked = false;
  main.style.display = "block";
}

init();
