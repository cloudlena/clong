// Convert vertical screen coordinates from relative to absolute
function absH(rel) {
  return Math.round(h * (rel / 100));
}

// Convert horizontal screen coordinates from relative to absolute
function absW(rel) {
  return Math.round(w * (rel / 100));
}

// Convert vertical screen coordinates from absolute to relative
function relH(abs) {
  return (abs / h) * -100;
}

// Convert horizontal screen coordinates from absolute to relative
function relW(abs) {
  return (abs / w) * 100;
}

// Calculate the points a certain target is worth
function calcPoints(posY, width, velocityX) {
  var bias = 5;
  var posYVal = 10 * (posY / 100);
  var widthVal = 10 * ((13 - width) / 11);
  var heightVal = 3 * ((13 - width) / 11);
  var velXVal = 25 * (2 * Math.abs(velocityX));
  return Math.round(bias + posYVal + widthVal + heightVal + velXVal);
}

// Calculate the current screen size
function calcScreenSize() {
  w = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
  h = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
}

// Check if a ball and a target collide
function doCollide(t, b) {
  var verticalRadius = b.radius * (w / h);
  var xMatch =
    b.posX + b.radius >= t.posX && b.posX - b.radius <= t.posX + t.width;
  var yMatch =
    b.posY + verticalRadius >= t.posY - t.height &&
    b.posY - verticalRadius <= t.posY;
  return xMatch && yMatch;
}

// Generate random integer
function randInt(min, max) {
  return Math.floor(Math.random() * (max - min)) + min;
}

// Generate random hex color value
function randomColor() {
  var n = Math.floor(Math.random() * 0x1000000);
  return "#" + n.toString(16).padStart(6, "0");
}

// Generate a UUID
function uuid() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1);
  }
  return [s4() + s4(), s4(), s4(), s4(), s4() + s4() + s4()].join("-");
}

// Return the WebSocket URL for a path on the current host
function wsURL(path) {
  var protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return protocol + "//" + window.location.host + path;
}
