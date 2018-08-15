"use strict";

// Set game preferences
var ballRadius = 0.7;
var ballVelocityFactor = 2.5;
var maxTargets = 5;
var targetSpawnIntervalMin = 100;
var targetSpawnIntervalMax = 500;
var targetTypes = [
  {
    imgUrl: "/assets/img/swisscom-logo-lifeform.svg",
    sizeRatio: 1
  },

  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/docker-image.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/dotnetcore.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/go.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/java.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/node.js.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/php.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/python.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/buildpack/ruby.svg",
    sizeRatio: 1
  },

  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/atmos.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/elk.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/mariadb.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/mongodb.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/rabbitmq.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/redis.svg",
    sizeRatio: 1
  },
  {
    imgUrl:
      "https://api.lyra-836.appcloud.swisscom.com/static-media/images/service/smartmsg.svg",
    sizeRatio: 1
  }
];

// Initialize target images
for (var i = 0; i < targetTypes.length; i++) {
  targetTypes[i].img = new Image();
  targetTypes[i].img.src = targetTypes[i].imgUrl;
}

// Initialize global variables
var ws = null;
var reqId = null;
var spawnTargetTimeout = null;
var canvas = document.getElementById("clong-canvas");
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
  ws = new WebSocket(wsProtocol() + "//" + window.location.host + "/ws/screen");

  ws.onopen = function() {
    // Show game URL on screen
    $("#game-msg").text(
      "Visit " +
        window.location.protocol +
        "//" +
        window.location.host +
        " to play"
    );

    // Trigger forever loops
    spawnTarget(randInt(targetSpawnIntervalMin, targetSpawnIntervalMax));
    reqId = requestAnimationFrame(draw);
  };

  // Listen for new balls coming in
  ws.onmessage = function(e) {
    var msg = JSON.parse(e.data);
    if (msg.type === "ballInit") {
      if (msg.velocityY > 0) {
        msg.radius = ballRadius;
        msg.velocityY = Math.max(msg.velocityY, 0.05) * ballVelocityFactor;

        balls.push(msg);
      } else {
        ballDone(msg.player, 0);
      }
    }
  };

  // Try to reconnect on close
  ws.onclose = function() {
    $("#game-msg").text("Reconnecting...");
    clearTimeout(spawnTargetTimeout);
    targets = [];
    balls = [];
    cancelAnimationFrame(reqId);
    setTimeout(init, 3000);
  };
}

function ballDone(player, points) {
  var msg = {
    type: "ballDone",
    player: player,
    points: points
  };
  ws.send(JSON.stringify(msg));
}

// Regularly spawn new targets
function spawnTarget(ms) {
  spawnTargetTimeout = setTimeout(function() {
    if (targets.length < maxTargets) {
      var type = targetTypes[randInt(0, targetTypes.length)];
      var width = randInt(2, 13);
      var height = width * type.sizeRatio * (w / h);

      targets.push({
        img: type.img,
        posX: randInt(20, 80),
        posY: randInt(height, 100),
        velocityX: Math.random() - 0.5,
        velocityY: 0,
        width: width,
        height: height
      });
    }

    spawnTarget(randInt(targetSpawnIntervalMin, targetSpawnIntervalMax));
  }, ms);
}

// Draw current status onto canvas
function draw() {
  var now = Date.now();
  var dt = (now - lastDrawnTime) / 17;
  lastDrawnTime = now;

  ctx.canvas.width = w;
  ctx.canvas.height = h;

  // Clear canvas
  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

  // Draw balls
  for (var i = 0; i < balls.length; i++) {
    var b = balls[i];

    ctx.fillStyle = b.color;
    ctx.beginPath();
    ctx.arc(absW(b.posX), absH(100 - b.posY), absW(b.radius), 0, 2 * Math.PI);
    ctx.fill();

    balls[i].posX += b.velocityX * dt;
    balls[i].posY += b.velocityY * dt;

    if (b.posY >= 100 || b.posX <= 0 || b.posX >= 100) {
      ballDone(b.player, 0);
      balls.splice(i, 1);
    }
  }

  // Draw targets
  for (var i = 0; i < targets.length; i++) {
    var t = targets[i];

    ctx.drawImage(
      t.img,
      absW(t.posX),
      absH(100 - t.posY),
      absW(t.width),
      absH(t.height)
    );

    // Invert velocity if target reaches end of screen
    if (
      (t.posX <= 0 && t.velocityX < 0) ||
      (t.posX >= 100 - t.width && t.velocityX > 0)
    ) {
      targets[i].velocityX = -t.velocityX;
    }

    targets[i].posX += t.velocityX * dt;
    targets[i].posY += t.velocityY * dt;

    for (var j = 0; j < balls.length; j++) {
      var b = balls[j];
      if (doCollide(t, b)) {
        var points = calcPoints(t.posY, t.width, t.height, t.velocityX);
        ballDone(b.player, points);
        targets.splice(i, 1);
        balls.splice(j, 1);
      }
    }
  }

  requestAnimationFrame(draw);
}

init();
