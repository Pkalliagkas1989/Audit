// register.js
document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const usernameInput = document.getElementById('username');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirmPassword');
    const messageDisplay = document.getElementById('message');
    const googleRegisterBtn = document.getElementById('googleRegisterBtn');
    const githubRegisterBtn = document.getElementById('githubRegisterBtn');

    // --- Utility Functions ---

    /**
     * Displays a message on the page.
     * @param {string} msg The message to display.
     * @param {boolean} isSuccess True for success message (greenish), false for error (reddish).
     */
    function displayMessage(msg, isSuccess) {
        messageDisplay.textContent = msg;
        messageDisplay.className = ''; // Clear existing classes
        messageDisplay.classList.add(isSuccess ? 'success' : 'error');
        // Ensure the error color is from --color-warning and success from --color-accent
        messageDisplay.style.color = isSuccess ? 'var(--color-accent)' : 'var(--color-warning)';
    }

    /**
     * Basic client-side validation for registration form.
     * @returns {boolean} True if all fields are valid, false otherwise.
     */
    function validateForm() {
        const username = usernameInput.value.trim();
        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();
        const confirmPassword = confirmPasswordInput.value.trim();

        if (!username || !email || !password || !confirmPassword) {
            displayMessage('All fields are required.', false);
            return false;
        }

        // Username validation (basic client-side, backend has stricter)
        if (username.length < 3 || username.length > 50 || !/^[a-zA-Z0-9_]+$/.test(username)) {
            displayMessage('Username must be 3-50 characters, letters/numbers/underscores only.', false);
            return false;
        }

        // Email validation (basic client-side, backend has stricter)
        if (!/\S+@\S+\.\S+/.test(email)) {
            displayMessage('Please enter a valid email address.', false);
            return false;
        }

        // Password validation (basic client-side, backend has stricter)
        if (password.length < 8) {
            displayMessage('Password must be at least 8 characters long.', false);
            return false;
        }
        if (!/[a-zA-Z]/.test(password) || !/\d/.test(password)) {
            displayMessage('Password must contain at least one letter and one digit.', false);
            return false;
        }

        if (password !== confirmPassword) {
            displayMessage('Passwords do not match.', false);
            return false;
        }

        return true;
    }

    // --- Manual Registration ---

    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent default form submission

        displayMessage('', true); // Clear previous messages

        if (!validateForm()) {
            return; // Stop if client-side validation fails
        }

        const username = usernameInput.value.trim();
        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();

        // Disable form elements during submission
        setFormEnabled(false);

        try {
            const response = await fetch('http://localhost:8080/forum/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    email: email,
                    password: password,
                }),
            });

            const data = await response.json();

            if (response.ok) {
                displayMessage('Registration successful! Redirecting to login...', true);
                // Redirect to login page after a short delay
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } else {
                // Handle backend errors
                const errorMessage = data.message || 'Registration failed. Please try again.';
                displayMessage(errorMessage, false);
            }
        } catch (error) {
            console.error('Network error during registration:', error);
            displayMessage('Network error. Please check your connection and try again.', false);
        } finally {
            setFormEnabled(true); // Re-enable form elements
        }
    });

    /**
     * Enables or disables form input fields and buttons.
     * @param {boolean} enabled True to enable, false to disable.
     */
    function setFormEnabled(enabled) {
        usernameInput.disabled = !enabled;
        emailInput.disabled = !enabled;
        passwordInput.disabled = !enabled;
        confirmPasswordInput.disabled = !enabled;
        registerForm.querySelector('.button1').disabled = !enabled;
        googleRegisterBtn.disabled = !enabled;
        githubRegisterBtn.disabled = !enabled;
    }

    // --- OAuth Registration ---

    // Note: These URLs will initiate the OAuth flow on your backend.
    // Ensure your Go backend has handlers for these endpoints that redirect to the OAuth provider.
    const GOOGLE_AUTH_URL = 'http://localhost:8080/forum/api/auth/google';
    const GITHUB_AUTH_URL = 'http://localhost:8080/forum/api/auth/github';

    googleRegisterBtn.addEventListener('click', () => {
        window.location.href = GOOGLE_AUTH_URL;
    });

    githubRegisterBtn.addEventListener('click', () => {
        window.location.href = GITHUB_AUTH_URL;
    });

    // Handle potential redirect messages after OAuth (e.g., if registration fails via OAuth)
    const urlParams = new URLSearchParams(window.location.search);
    const oauthMessage = urlParams.get('oauth_message');
    const oauthSuccess = urlParams.get('oauth_success'); // 'true' or 'false'

    if (oauthMessage) {
        displayMessage(decodeURIComponent(oauthMessage), oauthSuccess === 'true');
        // Clean the URL to remove the message params
        history.replaceState({}, document.title, window.location.pathname);
    }
});
