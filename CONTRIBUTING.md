# Contributing to Whishper Reloaded

Thanks for your interest in improving **Whishper Reloaded**! Contributions of all
kinds are welcome — bug reports, feature requests, documentation, and code.

## Ways to contribute

- **Report a bug**: open a [bug report](https://github.com/DevDema/whishper-reloaded/issues/new?template=bug_report.md).
- **Request a feature**: open a [feature request](https://github.com/DevDema/whishper-reloaded/issues/new?template=feature_request.md).
- **Improve the docs**: even small fixes (typos, clarifications) are appreciated.
- **Submit code**: pick an open issue or propose a change via a pull request.

## Project structure

Whishper Reloaded is made of a few decoupled components:

- `transcription-api/` — the Faster-Whisper transcription API.
- `backend/` — coordinates frontend calls, the database, and background tasks.
- `frontend/` — the SvelteKit web UI.

Third-party services used at runtime: **LibreTranslate** (translation),
**MongoDB** (storage), and **Nginx** (reverse proxy).

## Local development setup

1. **Fork and clone** the repository:

   ```bash
   git clone https://github.com/<your-username>/whishper-reloaded.git
   cd whishper-reloaded
   ```

2. **Run the stack with Docker Compose.** The quickest way to get a working
   environment is to bring the full stack up as described in the
   [README](README.md#running-with-docker). Create your `.backend.env` and
   `.frontend.env` files first, then:

   ```bash
   docker compose up -d
   ```

3. **Work on a single component.** To iterate on the frontend or backend
   without rebuilding images every time, run that component locally and point
   its environment variables at the other containers. See each component's
   folder (`frontend/`, `backend/`, `transcription-api/`) for details.

## Pull request guidelines

- Create a topic branch from the default branch: `git checkout -b my-feature`.
- Keep changes focused; one logical change per pull request when possible.
- Write clear commit messages describing the *why*, not just the *what*.
- Make sure the app still builds and runs before opening the PR.
- Reference any related issue in the PR description (e.g. `Closes #123`).

## Code of conduct

Be respectful and constructive. We want this to be a welcoming project for
everyone.
