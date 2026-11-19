"use strict";

// Game preferences
const gameSeconds = 60; // must match gameDuration on the server
const ballTimeout = 15000; // ms after which a ball is considered lost
const swipeMinDistance = 10; // px
const swipeMinVelocity = 0.3; // px/ms
const swipeVelocityWindow = 50; // ms at the end of a swipe used to measure its velocity

const main = document.getElementById("main");
const spinner = document.getElementById("spinner");
const intro = document.getElementById("intro");
const reconnectingMsg = document.getElementById("reconnecting-msg");
const startButton = document.getElementById("start-button");
const secNum = document.getElementById("sec-num");
const scoreNum = document.getElementById("score-num");
const myColor = randomColor();
const outbox = [];
let ws = null;
let w = 0;
let h = 0;
let gameRunning = false;
let locked = false;
let lockTimeout = null;
let myPoints = 0;
let myTime = 0;
let pointerSamples = [];

ensureCookie("userid", () => crypto.randomUUID());
ensureCookie("username", askUserName);
main.style.backgroundColor = myColor;
spinner.style.color = myColor;
startButton.querySelector("button").addEventListener("click", startGame);
resize();
window.onresize = resize;

// Track pointer movements to detect swipes
main.addEventListener("pointerdown", (e) => {
  if (e.isPrimary) {
    main.setPointerCapture(e.pointerId);
    pointerSamples = [pointerSample(e)];
  }
});
main.addEventListener("pointermove", (e) => {
  if (e.isPrimary && pointerSamples.length > 0) {
    pointerSamples.push(pointerSample(e));
  }
});
main.addEventListener("pointercancel", () => {
  pointerSamples = [];
});
main.addEventListener("pointerup", (e) => {
  if (!e.isPrimary || pointerSamples.length === 0) {
    return;
  }
  const end = pointerSample(e);
  const start = pointerSamples[0];
  // Flicks accelerate, so velocity is measured at the end of the swipe
  const recent = pointerSamples.find((s) => end.t - s.t <= swipeVelocityWindow);
  pointerSamples = [];

  const distance = Math.hypot(end.x - start.x, end.y - start.y);
  const duration = end.t - start.t;
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
  ws = new WebSocket(wsURL("/ws/controller"));

  ws.onopen = () => {
    unlock();
    reconnectingMsg.hidden = true;
    flushOutbox();
  };

  // Unlock when our ball is done (the server only sends our own balls)
  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    if (gameRunning && msg.type === "BALL_DONE") {
      myPoints += msg.points;
      scoreNum.textContent = myPoints;
      unlock();
    }
  };

  ws.onclose = () => {
    lock();
    reconnectingMsg.hidden = false;
    setTimeout(init, reconnectDelay);
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
  const gameLoop = setInterval(() => {
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

function resize() {
  ({ w, h } = screenSize());
}

// Record the position of a pointer at the current time
function pointerSample(e) {
  return { x: e.clientX, y: e.clientY, t: e.timeStamp };
}

function getCookie(name) {
  const prefix = name + "=";
  const cookie = document.cookie.split("; ").find((c) => c.startsWith(prefix));
  return cookie ? decodeURIComponent(cookie.substring(prefix.length)) : null;
}

// Create a cookie if it doesn't exist yet
function ensureCookie(name, createValue) {
  if (getCookie(name) === null) {
    document.cookie = name + "=" + encodeURIComponent(createValue());
  }
}

// Ask the user for their name until they enter one the server accepts (1-30 characters)
function askUserName() {
  let userName = "";
  while (userName.length === 0 || [...userName].length > 30) {
    userName = prompt("What's your full name?", "") || "";
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

// Convert vertical screen coordinates from absolute to relative, pointing upwards
function relH(abs) {
  return (abs / h) * -100;
}

// Convert horizontal screen coordinates from absolute to relative
function relW(abs) {
  return (abs / w) * 100;
}

// Generate a random hex color
function randomColor() {
  const n = Math.floor(Math.random() * 0x1000000);
  return "#" + n.toString(16).padStart(6, "0");
}

init();
