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
    
    const name = nameInput.value.trim() || "New Gopher";

    fetch('/api/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: Math.floor(Math.random() * 9000) + 1000, name: name })
    })
    .then(res => res.json())
    .then(data => {
        // Remove empty message on first user
        if (userList.querySelector('.empty-msg')) userList.innerHTML = '';

        // Add success message
        responseBox.classList.remove('hidden');
        apiData.innerText = data.message;

        // Add User to List
        const userDiv = document.createElement('div');
        userDiv.className = 'user-item';
        userDiv.innerHTML = `
            <div>
                <div style="font-weight: 600;">${name}</div>
                <div style="font-size: 0.7rem; color: #94a3b8;">ID: ${data.id}</div>
            </div>
            <span style="font-size: 0.7rem; color: #22c55e;">● Active</span>
        `;
        
        userList.prepend(userDiv); // Add newest to top
        nameInput.value = ''; // Clear input
        
        // Hide success message after 3 seconds
        setTimeout(() => responseBox.classList.add('hidden'), 3000);
    });
}
