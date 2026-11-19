"use strict";

// ms to wait before reconnecting a closed WebSocket
const reconnectDelay = 3000;

// Return the current screen size
function screenSize() {
  return {
    w: Math.max(document.documentElement.clientWidth, window.innerWidth || 0),
    h: Math.max(document.documentElement.clientHeight, window.innerHeight || 0),
  };
}

// Return the WebSocket URL for a path on the current host
function wsURL(path) {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return protocol + "//" + window.location.host + path;
}
