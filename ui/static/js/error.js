window.addEventListener('DOMContentLoaded', () => {
  const params = new URLSearchParams(window.location.search);
  const msg = params.get('msg');

  const errorMessage = document.getElementById('error-message');

  if (msg && errorMessage) {
    errorMessage.textContent = decodeURIComponent(msg);
  } else if (errorMessage) {
    errorMessage.textContent = "An unexpected error occurred.";
  }

  document.getElementById("back-button")?.addEventListener("click", () => {
    window.history.back();
  });
});
// 