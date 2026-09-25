"use strict";

// Set game preferences
var ballRadius = 0.7;
var ballVelocityFactor = 2.5;
var maxTargets = 5;
var targetSpawnIntervalMin = 100;
var targetSpawnIntervalMax = 500;
var targetImages = [
  "docker",
  "go",
  "rust",
  "java",
  "nodejs",
  "python",
  "ruby",
  "mariadb",
  "mongodb",
  "rabbitmq",
  "redis",
  "kubernetes",
  "bespinian",
].map(function (name) {
  var img = new Image();
  img.src = "/assets/img/" + name + ".svg";
  return img;
});

// Initialize global variables
var ws = null;
var reqId = null;
var spawnTargetTimeout = null;
var canvas = document.getElementById("clong-canvas");
var gameMsg = document.getElementById("game-msg");
var ctx = canvas.getContext("2d");
var lastDrawnTime = Date.now();
var w = 0;
var h = 0;
var balls = [];
var targets = [];

function init() {
  calcScreenSize();
  window.onresize = calcScreenSize;

  // Initialize WebSocket connection
  ws = new WebSocket(wsURL("/ws/screen"));

  ws.onopen = function () {
    // Show game URL on screen
    gameMsg.textContent =
      "Visit " +
      window.location.protocol +
      "//" +
      window.location.host +
      " to play";

    // Trigger forever loops
    spawnTarget();
    reqId = requestAnimationFrame(draw);
  };

  // Listen for new balls coming in
  ws.onmessage = function (e) {
    var msg = JSON.parse(e.data);
    if (msg.type !== "BALL_INIT") {
      return;
    }
    if (msg.velocityY <= 0) {
      ballDone(msg.player, 0);
      return;
    }
    msg.radius = ballRadius;
    msg.velocityY = Math.max(msg.velocityY, 0.05) * ballVelocityFactor;
    balls.push(msg);
  };

  // Try to reconnect on close
  ws.onclose = function () {
    gameMsg.textContent = "Reconnecting...";
    clearTimeout(spawnTargetTimeout);
    cancelAnimationFrame(reqId);
    targets = [];
    balls = [];
    setTimeout(init, 3000);
  };
}

function ballDone(player, points) {
  var msg = {
    type: "BALL_DONE",
    player: player,
    points: points,
  };
  ws.send(JSON.stringify(msg));
}

// Regularly spawn new targets
function spawnTarget() {
  var delay = randInt(targetSpawnIntervalMin, targetSpawnIntervalMax);
  spawnTargetTimeout = setTimeout(function () {
    if (targets.length < maxTargets) {
      var width = randInt(2, 13);
      var height = width * (w / h);

      targets.push({
        img: targetImages[randInt(0, targetImages.length)],
        posX: randInt(20, 80),
        posY: randInt(height, 100),
        velocityX: Math.random() - 0.5,
        width: width,
        height: height,
      });
    }

    spawnTarget();
  }, delay);
}

// Draw current status onto canvas
function draw() {
  var now = Date.now();
  var dt = (now - lastDrawnTime) / 17;
  lastDrawnTime = now;

  canvas.width = w;
  canvas.height = h;
  ctx.clearRect(0, 0, w, h);

  // Draw and move balls, removing the ones that left the screen
  balls = balls.filter(function (b) {
    ctx.fillStyle = b.color;
    ctx.beginPath();
    ctx.arc(absW(b.posX), absH(100 - b.posY), absW(b.radius), 0, 2 * Math.PI);
    ctx.fill();

    b.posX += b.velocityX * dt;
    b.posY += b.velocityY * dt;

    var isOut = b.posY >= 100 || b.posX <= 0 || b.posX >= 100;
    if (isOut) {
      ballDone(b.player, 0);
    }
    return !isOut;
  });

  // Draw and move targets, removing the ones that got hit
  targets = targets.filter(function (t) {
    ctx.drawImage(
      t.img,
      absW(t.posX),
      absH(100 - t.posY),
      absW(t.width),
      absH(t.height),
    );

    // Invert velocity if target reaches end of screen
    if (
      (t.posX <= 0 && t.velocityX < 0) ||
      (t.posX >= 100 - t.width && t.velocityX > 0)
    ) {
      t.velocityX = -t.velocityX;
    }
    t.posX += t.velocityX * dt;

    var hit = balls.findIndex(function (b) {
      return doCollide(t, b);
    });
    if (hit === -1) {
      return true;
    }
    ballDone(balls[hit].player, calcPoints(t.posY, t.width, t.velocityX));
    balls.splice(hit, 1);
    return false;
  });

  reqId = requestAnimationFrame(draw);
}

init();
