import * as gameConnection from './game-connection.js';

gameConnection.GetCookie();
gameConnection.ConnectToWebSocket();

var spacebarPressed = false;
var score = 0;
var sendScoreFunc = null;

// Keypresses
document.addEventListener('keydown', function(event) {
    if (event.code === 'Space' && !spacebarPressed) {
        fillBar();
        spacebarPressed = true;
    }
});

document.addEventListener('keyup', function(event) {
    if (event.code === 'Space') {
        spacebarPressed = false;
    }
});

// Connect Event
document.addEventListener('gameDataUpdate', function(event) {
    //console.log(JSON.stringify(event.detail));
    //console.log("\nSessionId: " + event.detail.SessionId)
    //console.log("SessionStart: " + event.detail.SessionStartTime)
    console.log(event.detail.PlayerScores)
    if (Object.keys(event.detail.PlayerScores).length > 1) {
        console.log("2 Players")
        document.getElementById('game').style.display = 'block'
        document.getElementById('cover').style.display = 'none'
    }
    else {
        console.log("1 Player")
        document.getElementById('game').style.display = 'none'
        document.getElementById('cover').style.display = 'block'
    }
    document.getElementById('join-link').href = `http://localhost:8080/static?id=${event.detail.SessionId}`;
    document.getElementById('join-link').innerText = `http://localhost:8080/static?id=${event.detail.SessionId}`;
    console.log("\n" + JSON.stringify(event.detail.PlayerScores))
  });

// Graphics
function fillBar() {
    if (score >= 100) return
    document.getElementById('score').innerText = ++score;
    var fill = document.getElementById('fill');
    fill.style.width = document.getElementById('container').offsetWidth * (score/100) + 'px'
    gameConnection.SendData(score)
}
