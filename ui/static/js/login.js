document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const email = document.getElementById("email").value.trim();
  const password = document.getElementById("password").value;
  const message = document.getElementById("message");

  try {
    const response = await fetch("http://localhost:8080/forum/api/session/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include", // IMPORTANT to send and receive cookies
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (response.ok) {
      message.textContent = data.message;
      message.style.color = "green";
      // Redirect after successful login
      setTimeout(() => {
        window.location.href = "/user";  // Redirect to user page
      }, 1000);
    } else {
      message.textContent = data.message || "Login failed!";
      message.style.color = "red";
    }
  } catch (error) {
    message.textContent = "Error connecting to server.";
    message.style.color = "red";
  }
});
