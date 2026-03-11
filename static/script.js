function testAPI() {
    fetch('/api/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: 101, name: 'Gopher' })
    })
    .then(res => res.json())
    .then(data => {
        console.log('Success:', data);
        alert('Server received User: ' + data.name);
    })
    .catch(err => console.error('Error:', err));
}
