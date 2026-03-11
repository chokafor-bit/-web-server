function testAPI() {
    const display = document.getElementById('api-data');
    const box = document.getElementById('response-box');

    fetch('/api/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: 1, name: 'Gopher' })
    })
    .then(res => res.json())
    .then(data => {
        box.classList.remove('hidden'); // Show the box
        display.innerText = JSON.stringify(data, null, 2); // Format JSON nicely
    })
    .catch(err => {
        display.innerText = "Error: " + err;
    });
}
