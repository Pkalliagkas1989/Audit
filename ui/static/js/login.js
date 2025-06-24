document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('loginForm');
  const messageEl = document.getElementById('message');

  function showMessage(msg, isError = false) {
    if (!messageEl) return;
    messageEl.textContent = msg;
    messageEl.style.color = isError ? 'var(--color-warning)' : 'var(--color-accent)';
  }

  form?.addEventListener('submit', async (e) => {
    e.preventDefault();
    showMessage('');

    const email = document.getElementById('email').value.trim().toLowerCase();
    const password = document.getElementById('password').value;

    if (!email || !password) {
      showMessage('Email and password are required', true);
      return;
    }

    try {
      const resp = await fetch('http://localhost:8080/forum/api/session/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email, password })
      });

      const data = await resp.json().catch(() => ({}));

      if (resp.ok) {
        showMessage('Login successful! Redirecting...');
        setTimeout(() => {
          window.location.href = '/user';
        }, 1500);
      } else {
        const msg = data.message || data.error || 'Login failed';
        showMessage(msg, true);
      }
    } catch (err) {
      showMessage('Network error. Please try again.', true);
    }
  });
});
