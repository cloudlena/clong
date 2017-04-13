function abs(percentage) {
    return Math.max(w, h) * (percentage/100);
}

function absH(percentage) {
    return h * (percentage/100);
}

function absW(percentage) {
    return w * (percentage/100);
}

function calcScreenSize() {
    w = Math.max(document.documentElement.clientWidth, window.innerWidth || 0);
    h = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
}

function doCollide(t, b) {
    var xMatch = false;
    var yMatch = false;

    if (b.posX+b.radius >= t.posX && b.posX-b.radius <= t.posX+t.width) {
        xMatch = true;
    }
    if (b.posY+b.radius >= t.posY && b.posY-b.radius <= t.posY+t.height) {
        yMatch = true;
    }

    if (xMatch && yMatch) {
        return true;
    } else {
        return false
    }
}

function randInt(min, max) {
    return Math.floor(Math.random() * (max - min)) + min;
}

function randomColor() {
    var letters = '0123456789ABCDEF';
    var color = '#';

    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }

    return color;
}

function rel(val) {
    return (val / Math.max(w,h)) * 100;
}

function relH(val) {
    return (val / h) * -100;
}

function relW(val) {
    return (val / w) * 100;
}

function wsProtocol() {
    var p = 'ws';

    if (window.location.protocol === "https:") {
        p = 'wss';
    }

    return p + ':';
}
