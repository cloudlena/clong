// Calculate the current screen size into the globals w and h
function calcScreenSize() {
  w = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
  h = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
}

// Return the WebSocket URL for a path on the current host
function wsURL(path) {
  var protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return protocol + "//" + window.location.host + path;
}
