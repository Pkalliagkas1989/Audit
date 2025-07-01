# Forum Application

This repository contains a simple forum backend (API) and frontend (UI). The API uses Google OAuth for authentication.

## Quick Start

1. Copy `.env.example` to `.env` and provide your Google OAuth credentials:
   ```bash
   cp .env.example .env
   # edit .env and fill GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET
   ```
2. Ensure `GOOGLE_CALLBACK_URL` in `.env` matches the authorized redirect URI in your Google console.
3. Build and start the application using Docker Compose:
   ```bash
   docker compose up --build
   ```
4. Open `http://localhost:8081` in your browser. Register with Google should now redirect correctly.

## Manual API Testing

See `instructions.md` for curl commands to interact with the API.

