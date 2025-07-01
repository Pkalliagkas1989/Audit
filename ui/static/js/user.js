window.addEventListener('popstate', async () => {
  // User pressed back or forward button
  try {
    await fetch('/forum/api/session/logout', { method: 'POST', credentials: 'include' });
    // After logout, redirect to login page
    window.location.href = '/login';
  } catch (e) {
    console.error('Logout failed', e);
  }
});
