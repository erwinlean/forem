"use strict";

function checkToken() {
  const token = localStorage.getItem('token');
  const tokenTimestamp = localStorage.getItem('token-timestamp');
  const oneDayInMilliseconds = 24 * 60 * 60 * 1000;

  if (!token || !tokenTimestamp || (Date.now() - tokenTimestamp) > oneDayInMilliseconds) {
      localStorage.removeItem('token');
      localStorage.removeItem('token-timestamp');
      window.location.href = 'login.html';
  }
}

document.addEventListener('DOMContentLoaded', function() {
  checkToken();

  const scrappingButton = document.getElementById('scrapping_mitutoyo');
  scrappingButton.addEventListener('click', async function() {
      const token = localStorage.getItem('token');
      const email = localStorage.getItem('email');

      const response = await fetch('http://127.0.0.1:8000/scrapper/mitutoyo', {
          method: 'GET',
          headers: {
              'Authorization': token,
              'X-User-Email': email
          }
      });

      if (response.ok) {
          alert('Scraping iniciado correctamente');
      } else {
          alert('Error al iniciar el scraping');
      }
  });
});

document.addEventListener('DOMContentLoaded', function() {
  const logoutButton = document.getElementById('logout-button');
  logoutButton.addEventListener('click', function() {
      localStorage.removeItem('token');
      localStorage.removeItem('token-timestamp');
      localStorage.removeItem('email');
      window.location.href = 'index.html';
  });
});
