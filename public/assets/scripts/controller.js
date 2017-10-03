"use strict";

// Initialize global variables
var ws = null;
var main = document.getElementById("main");
var spinner = document.getElementById("spinner");
var mc = new Hammer(main);
var w = 0;
var h = 0;
var myID = "";
var myColor = "white";
var gameRunning = false;
var locked = false;
var myPoints = 0;
var myTime = 0;

// Init player
myID = setUserID();
setUserName();
myColor = randomColor();
main.style.backgroundColor = myColor;
spinner.style.color = myColor;

function init() {
    calcScreenSize();
    window.onresize = calcScreenSize;

    // Initialize WebSocket connection
    ws = new WebSocket(wsProtocol() + "//" + window.location.host + "/ws/controller");

    ws.onopen = function() {
        unlock();
        $("#reconnecting-msg").hide();

        // Listen for new swipes
        mc.get("swipe").set({ direction: Hammer.DIRECTION_ALL });
        mc.on("swipe", function(e) {
            if (gameRunning && !locked) {
                $("#intro").hide();

                var msg = {
                    type: "ballInit",
                    color: myColor,
                    posX: relW(e.center.x),
                    velocityX: relW(e.velocityX),
                    velocityY: relH(e.velocityY),
                };
                ws.send(JSON.stringify(msg));

                lock();
            }
        });
    };

    // Listen for when ball is done to unlock screen
    ws.onmessage = function(e) {
        var msg = JSON.parse(e.data);
        if (gameRunning && msg.type === "ballDone" && msg.player.id === myID) {
            myPoints += msg.points;
            $("#score-num").text(myPoints);
            unlock();
        }
    };

    // Try to reconnect on close
    ws.onclose = function() {
        lock();
        $("#reconnecting-msg").show();
        setTimeout(init, 3000);
    };
}

function startGame() {
    $("#start-button").hide();
    myTime = 60;
    $("#sec-num").text(myTime);
    myPoints = 0;
    $("#score-num").text(myPoints);
    gameRunning = true;
    var gameLoop = setInterval(function() {
        myTime--;
        $("#sec-num").text(myTime);
        if (myTime <= 0) {
            gameRunning = false;
            unlock();
            $("#start-button").show();
            clearInterval(gameLoop);
            var msg = {
                type: "gameDone",
                finalScore: myPoints,
                color: myColor,
            };
            ws.send(JSON.stringify(msg));
        }
    }, 1000);
}

// Sets a cookie in the cookie jar
function setCookie(name, value) {
    document.cookie = name + "=" + value;
}

// Gets a cookie from the cookie jar
function getCookie(name) {
    var nameEQ = name + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(";");

    for(var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === " ") { c = c.substring(1); }
        if (c.indexOf(nameEQ) === 0) { return c.substring(nameEQ.length, c.length); }
    }

    return null;
}

function setUserID() {
    var userid = getCookie("userid");

    if (userid === null) {
        userid = uuid();
        setCookie("userid", userid);
    }

    return userid;
}

function getUserName() {
    var userName = prompt("Please enter your first and last name (2-30 characters)", "");

    if (userName === null || userName.length < 2 || userName.length > 30) {
        userName = getUserName()
    }

    return userName;
}

function setUserName() {
    var userName = getCookie("username");

    if (userName === null) {
        userName = getUserName();
        setCookie("username", userName);
    }

    return userName;
}

function lock() {
    locked = true;
    main.style.display = "none";
}

function unlock() {
    locked = false;
    main.style.display = "block";
}

if (window.WebSocket) {
    init();
} else {
    $("#start-button .btn").hide();
    $(".hud").hide();
    $("#start-button").append("<p style=\"text-align:center;margin-top:2em;\">Sorry, your browser isn't supported...</p>");
}
