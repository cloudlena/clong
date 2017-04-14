// Convert screen coordinates from relative to absolute
function abs(rel) {
    return Math.max(w, h) * (rel/100);
}

// Convert vertical screen coordinates from relative to absolute
function absH(rel) {
    return h * (rel/100);
}

// Convert horizontal screen coordinates from relative to absolute
function absW(rel) {
    return w * (rel/100);
}

// Convert screen coordinates from absolute to relative
function rel(abs) {
    return (abs / Math.max(w,h)) * 100;
}

// Convert vertical screen coordinates from absolute to relative
function relH(abs) {
    return (abs / h) * -100;
}

// Convert horizontal screen coordinates from absolute to relative
function relW(abs) {
    return (abs / w) * 100;
}

// Calculate the current screen size
function calcScreenSize() {
    w = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
    h = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
}

// Check if a ball and a target collide
function doCollide(t, b) {
    var xMatch = false;
    var yMatch = false;

    if (b.posX+b.radius >= t.posX && b.posX-b.radius <= t.posX+t.width) {
        xMatch = true;
    }
    if (b.posY+b.radius >= t.posY-t.height && b.posY-b.radius <= t.posY) {
        yMatch = true;
    }

    if (xMatch && yMatch) {
        return true;
    }

    return false
}

// Generate random integer
function randInt(min, max) {
    return Math.floor(Math.random() * (max - min)) + min;
}

// Generate random hex color value
function randomColor() {
    const letters = '0123456789ABCDEF';

    var color = '#';

    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }

    return color;
}

// Return correct WebSocket protocol
function wsProtocol() {
    var p = 'ws';

    if (window.location.protocol === "https:") {
        p = 'wss';
    }

    return p + ':';
}
