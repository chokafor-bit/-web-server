function updateClock() {
    document.getElementById('live-clock').innerText = new Date().toLocaleTimeString('en-GB', { hour12: false });
}
setInterval(updateClock, 1000);
updateClock();

window.onload = function() {
    fetch('/api/status').then(res => res.json()).then(data => {
        document.getElementById('start-time').innerText = `Started: ${data.started_at} WAT`;
    });
};

function testAPI() {
    const nameInput = document.getElementById('userName');
    const userList = document.getElementById('user-list');
    const responseBox = document.getElementById('response-box');
    const apiData = document.getElementById('api-data');
    
    const name = nameInput.value.trim();
    if (!name) return;

    // Generate a random ID to send to the server
    const generatedId = Math.floor(Math.random() * 9000) + 1000;

    fetch('/api/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: generatedId, name: name })
    })
    .then(res => res.json().then(data => ({ status: res.status, data })))
    .then(({ status, data }) => {
        if (status !== 201) {
            // ERROR CASE (Duplicate)
            responseBox.classList.remove('hidden');
            apiData.innerText = data.message;
            apiData.style.color = "#ef4444"; // Red for error
            return;
        }

        // SUCCESS CASE
        responseBox.classList.remove('hidden');
        apiData.innerText = data.message;
        apiData.style.color = "var(--success-glow)"; // Green for success

        if (userList.querySelector('.empty-msg')) userList.innerHTML = '';

        // CREATE USER ITEM WITH ID
        const userDiv = document.createElement('div');
        userDiv.className = 'user-item';
        userDiv.innerHTML = `
            <div>
                <div style="font-weight: 600;">${data.name}</div>
                <div style="font-size: 0.7rem; color: #94a3b8;">ID: ${data.id}</div>
            </div>
            <span style="font-size: 0.7rem; color: #22c55e;">● Active</span>
        `;
        
        userList.prepend(userDiv);
        nameInput.value = '';
        
        // Hide success message after 4 seconds
        setTimeout(() => responseBox.classList.add('hidden'), 4000);
    })
    .catch(err => alert("Server Error: " + err))
}
