"use strict";

// Game preferences
const ballRadius = 0.7;
const ballVelocityFactor = 2.5;
const maxTargets = 5;
const targetSpawnIntervalMin = 100;
const targetSpawnIntervalMax = 500;
const targetWidthMin = 2;
const targetWidthMax = 13;
const targetImages = [
  "anthropic",
  "bespinian",
  "gemini",
  "go",
  "kubernetes",
  "mistral",
  "nine-dark",
  "openai",
  "python",
  "qwen",
  "redis",
  "rust",
  "stepping-stone",
].map((name) => {
  const img = new Image();
  img.src = "/assets/img/" + name + ".svg";
  return img;
});

const canvas = document.getElementById("clong-canvas");
const gameMsg = document.getElementById("game-msg");
const ctx = canvas.getContext("2d");
let ws = null;
let reqId = null;
let spawnTargetTimeout = null;
let lastDrawnTime = Date.now();
let w = 0;
let h = 0;
let balls = [];
let targets = [];

resize();
window.onresize = resize;

function init() {
  ws = new WebSocket(wsURL("/ws/screen"));

  ws.onopen = () => {
    gameMsg.textContent =
      "Visit " +
      window.location.protocol +
      "//" +
      window.location.host +
      " to play";
    spawnTarget();
    reqId = requestAnimationFrame(draw);
  };

  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    if (msg.type !== "BALL_INIT") {
      return;
    }
    if (msg.velocityY <= 0) {
      ballDone(msg.player, 0);
      return;
    }
    msg.posY = 0;
    msg.radius = ballRadius;
    msg.velocityY = Math.max(msg.velocityY, 0.05) * ballVelocityFactor;
    balls.push(msg);
  };

  // Show why the server closed the connection, if it did
  ws.onclose = (e) => {
    gameMsg.textContent = e.reason || "Reconnecting...";
    clearTimeout(spawnTargetTimeout);
    cancelAnimationFrame(reqId);
    targets = [];
    balls = [];
    setTimeout(init, reconnectDelay);
  };
}

function ballDone(player, points) {
  if (ws.readyState !== WebSocket.OPEN) {
    return;
  }
  ws.send(
    JSON.stringify({ type: "BALL_DONE", player: player, points: points }),
  );
}

// Match the canvas size to the screen size
function resize() {
  ({ w, h } = screenSize());
  canvas.width = w;
  canvas.height = h;
}

// Regularly spawn new targets
function spawnTarget() {
  const delay = randInt(targetSpawnIntervalMin, targetSpawnIntervalMax);
  spawnTargetTimeout = setTimeout(() => {
    // Skip images that are still loading or failed to load
    const imgs = targetImages.filter((i) => i.complete && i.naturalWidth > 0);
    if (targets.length < maxTargets && imgs.length > 0) {
      const img = imgs[randInt(0, imgs.length)];
      const width = randInt(targetWidthMin, targetWidthMax);
      // Keep the image's aspect ratio
      const height = width * (w / h) * (img.naturalHeight / img.naturalWidth);

      targets.push({
        img: img,
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

function draw() {
  const now = Date.now();
  const dt = (now - lastDrawnTime) / 17;
  lastDrawnTime = now;

  ctx.clearRect(0, 0, w, h);

  // Draw and move balls, removing the ones that left the screen
  balls = balls.filter((b) => {
    ctx.fillStyle = b.color;
    ctx.beginPath();
    ctx.arc(absW(b.posX), absH(100 - b.posY), absW(b.radius), 0, 2 * Math.PI);
    ctx.fill();

    b.posX += b.velocityX * dt;
    b.posY += b.velocityY * dt;

    const isOut = b.posY >= 100 || b.posX <= 0 || b.posX >= 100;
    if (isOut) {
      ballDone(b.player, 0);
    }
    return !isOut;
  });

  // Draw and move targets, removing the ones that got hit
  targets = targets.filter((t) => {
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

    const hit = balls.findIndex((b) => doCollide(t, b));
    if (hit === -1) {
      return true;
    }
    ballDone(balls[hit].player, calcPoints(t.posY, t.width, t.velocityX));
    balls.splice(hit, 1);
    return false;
  });

  reqId = requestAnimationFrame(draw);
}

// Convert vertical screen coordinates from relative to absolute
function absH(rel) {
  return Math.round(h * (rel / 100));
}

// Convert horizontal screen coordinates from relative to absolute
function absW(rel) {
  return Math.round(w * (rel / 100));
}

// Calculate the points a target is worth. Higher, smaller and faster targets are worth more.
function calcPoints(posY, width, velocityX) {
  const bias = 5;
  const posYVal = 10 * (posY / 100);
  const smallnessVal =
    13 * ((targetWidthMax - width) / (targetWidthMax - targetWidthMin));
  const velXVal = 50 * Math.abs(velocityX);
  return Math.round(bias + posYVal + smallnessVal + velXVal);
}

// Check if a ball and a target collide
function doCollide(t, b) {
  const verticalRadius = b.radius * (w / h);
  const xMatch =
    b.posX + b.radius >= t.posX && b.posX - b.radius <= t.posX + t.width;
  const yMatch =
    b.posY + verticalRadius >= t.posY - t.height &&
    b.posY - verticalRadius <= t.posY;
  return xMatch && yMatch;
}

// Generate a random integer in [min, max)
function randInt(min, max) {
  return Math.floor(Math.random() * (max - min)) + min;
}

init();
