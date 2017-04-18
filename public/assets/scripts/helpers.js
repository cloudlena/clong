// Convert vertical screen coordinates from relative to absolute
function absH(rel) {
    return h * (rel/100);
}

// Convert horizontal screen coordinates from relative to absolute
function absW(rel) {
    return w * (rel/100);
}

// Convert vertical screen coordinates from absolute to relative
function relH(abs) {
    return (abs / h) * -100;
}

// Convert horizontal screen coordinates from absolute to relative
function relW(abs) {
    return (abs / w) * 100;
}

function calcPoints(posY, width, height, velocityX) {
    if (velocityX < 0) {
        velocityX = -velocityX;
    }
    var bias = 5;
    var posYVal = 10 * (posY/100);
    var widthVal = 10 * ((13-width)/11);
    var heightVal = 10 * ((13-width)/11);
    var velXVal = 25 * (2*velocityX);
    var finalVal = Math.round(bias + posYVal + widthVal + heightVal + velXVal);
    return finalVal;
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
    var verticalRadius = b.radius * (w/h);

    if (b.posX+b.radius >= t.posX && b.posX-b.radius <= t.posX+t.width) {
        xMatch = true;
    }
    if (b.posY+verticalRadius >= t.posY-t.height && b.posY-verticalRadius <= t.posY) {
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
    var letters = '0123456789ABCDEF';

    var color = '#';

    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }

    return color;
}

// Generate a UUID
function uuid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
        s4() + '-' + s4() + s4() + s4();
}

// Return correct WebSocket protocol
function wsProtocol() {
    var p = 'ws';

    if (window.location.protocol === "https:") {
        p = 'wss';
    }

    return p + ':';
}
