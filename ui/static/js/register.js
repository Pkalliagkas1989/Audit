document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('registerForm');
  const messageEl = document.getElementById('message');

  function showMessage(msg, isError = false) {
    if (!messageEl) return;
    messageEl.textContent = msg;
    messageEl.style.color = isError ? 'var(--color-warning)' : 'var(--color-accent)';
  }

  form?.addEventListener('submit', async (e) => {
    e.preventDefault();
    showMessage('');

    const username = document.getElementById('username').value.trim();
    const email = document.getElementById('email').value.trim().toLowerCase();
    const password = document.getElementById('password').value;
    const confirm = document.getElementById('confirmPassword').value;

    if (password !== confirm) {
      showMessage('Passwords do not match', true);
      return;
    }

    try {
      const resp = await fetch('http://localhost:8080/forum/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username, email, password })
      });

      const data = await resp.json().catch(() => ({}));

      if (resp.ok) {
        showMessage('Registration successful! Redirecting...');
        setTimeout(() => {
          window.location.href = '/login';
        }, 1500);
      } else {
        const msg = data.message || data.error || 'Registration failed';
        showMessage(msg, true);
      }
    } catch (err) {
      showMessage('Network error. Please try again.', true);
    }
  });
});
