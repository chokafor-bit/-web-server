// 1. Live Clock
function updateClock() {
    document.getElementById('live-clock').innerText = new Date().toLocaleTimeString();
}
setInterval(updateClock, 1000);
updateClock();

// 2. Load Server Start Time on Page Load
window.onload = function() {
    fetch('/api/status')
        .then(res => res.json())
        .then(data => {
            document.getElementById('start-time').innerText = `Started: ${data.started_at}`;
        })
        .catch(() => console.error("Could not fetch server status"));
};

// 3. Send User Data to Go
function testAPI() {
    const nameInput = document.getElementById('userName');
    const responseBox = document.getElementById('response-box');
    const apiData = document.getElementById('api-data');
    
    const name = nameInput.value.trim() || "Anonymous Gopher";

    fetch('/api/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: Math.floor(Math.random() * 100), name: name })
    })
    .then(res => res.json())
    .then(data => {
        responseBox.classList.remove('hidden');
        apiData.innerText = data.message; // Uses the message from your Go backend
        nameInput.value = ''; // Clear input
    })
    .catch(err => alert("Error connecting to server"));
}
