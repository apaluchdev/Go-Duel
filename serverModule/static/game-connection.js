var websocket 

export function ConnectToWebSocket() {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = urlParams.get('id');

    websocket = new WebSocket(`ws://localhost:8080/session/connect?id=${sessionId}`);

    websocket.onopen = function() {
    console.log("WebSocket connected.");
    };

    websocket.onclose = function() {
        console.log("WebSocket disconnected.");
    };

    websocket.onerror = function(event) {
        console.error("WebSocket error:", event);
    };

    websocket.onmessage = function(event) {
        const customEvent = new CustomEvent('gameDataUpdate', {
            detail: JSON.parse(event.data)
        });
        document.dispatchEvent(customEvent);
    };
}

export function SendData(score) {
    if (!websocket) return

    if (websocket.readyState === WebSocket.OPEN) {
        const data = {
            Score: score
        };
        websocket.send(JSON.stringify(data));
        console.log("Data sent:", data);
    } else {
        console.error("WebSocket connection is not open.");
    }
}

export function GetCookie() {
    fetch('http://localhost:8080/session/setuserid', {
    method: 'GET',
    credentials: 'same-origin' // include cookies in the request
    })
    .then(response => {
    // Handle response
    if (response.ok) {
        console.log('Cookie set successfully');
    } else {
        alert('Failed to set cookie');
    }
    })
    .catch(error => {
    console.error('Error:', error);
    });
}