// /static/js/error.js

function getQueryParam(name) {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(name);
}

window.addEventListener('DOMContentLoaded', () => {
  const message = getQueryParam('message') || 'An unknown error occurred.';
  document.getElementById('error-message').textContent = decodeURIComponent(message);
});

