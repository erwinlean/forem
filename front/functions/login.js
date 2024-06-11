"use strict";

document.getElementById('login-form').addEventListener('submit', async function(event) {
    event.preventDefault();
    const email = document.getElementById('email-address').value;
    const password = document.getElementById('password').value;

    const response = await fetch('http://127.0.0.1:8000/users/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email: email, password:  password })
    });

    if (response.ok) {
        const data = await response.json();
        const token = data.token;
        localStorage.setItem('token', token);
        localStorage.setItem('email', email); 
        localStorage.setItem('token-timestamp', Date.now());
        window.location.href = 'home.html';
    } else {
        alert('Login failed');
    }
});