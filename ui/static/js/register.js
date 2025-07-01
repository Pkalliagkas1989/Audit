document.getElementById("registerForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const email = document.getElementById("email").value.trim();
  const password = document.getElementById("password").value;
  const confirmPassword = document.getElementById("confirmPassword").value;
  const message = document.getElementById("message");

  if (password !== confirmPassword) {
    message.textContent = "Passwords do not match!";
    message.classList.remove("success");
    return;
  }

  try {
    const response = await fetch("http://localhost:8080/forum/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include", // ✅ Include cookies!
      body: JSON.stringify({ username, email, password }),
    });

    const data = await response.json();

    if (response.ok) {
      // ✅ Save CSRF token if needed
      if (data.csrf_token) {
        localStorage.setItem("csrfToken", data.csrf_token);
      }

      message.textContent = "Registration successful!";
      message.classList.add("success");

      setTimeout(() => {
        window.location.replace("/user");
      }, 1000);
    } else {
      message.textContent = data.message || "Registration failed!";
      message.classList.remove("success");
    }
  } catch (error) {
    message.textContent = "Error connecting to server.";
    message.classList.remove("success");
  }
});

document.getElementById("googleRegisterBtn").addEventListener("click", () => {
  window.location.href = "http://localhost:8080/auth/google/login";
});

document.getElementById("githubRegisterBtn").addEventListener("click", () => {
  window.location.href = "http://localhost:8080/auth/github/login";
});
