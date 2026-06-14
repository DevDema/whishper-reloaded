[![whishper banner](misc/banner.png)]

[![](https://img.shields.io/badge/website-066da5?style=for-the-badge&logo=icloud&logoColor=white)](https://whishper.net)
[![](https://img.shields.io/badge/self%20host%20guide-066da5?style=for-the-badge&logo=googledocs&logoColor=white)](https://whishper.net/guides/install)
[![](https://img.shields.io/badge/screenshots-5c1f87?style=for-the-badge&logo=slickpic&logoColor=white)](#screenshots)
[![](https://img.shields.io/docker/pulls/thespartan94/whishper-frontend?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/r/thespartan94/whishper-frontend)

**Whishper Reloaded** is an open-source, 100% local audio transcription and subtitling suite with a full-featured web UI. It is a fork of [Pluja's Whishper](https://github.com/pluja/whishper), further improved and now actively maintained.

> **Project status:** This project began as [Whishper](https://github.com/pluja/whishper) by [Pluja](https://github.com/pluja), who remains the original creator and main contributor of the core idea and codebase. Whishper Reloaded is now maintained solo by me, [devdema](https://github.com/devdema), going forward. If you find this useful, please consider supporting the original project too.

## Features

Everything from the original Whishper, plus a set of improvements added in this fork.

### Core features

- [x] 🗣️ **Transcribe any media** to text: audio, video, etc.
  - Transcribe from URLs (any source supported by yt-dlp).
  - Upload a file to transcribe.
- [x] 📥 **Download transcriptions in many formats**: TXT, JSON, VTT, SRT or copy the raw text to your clipboard.
- [x] 🌐 **Translate your transcriptions** to any language supported by [Libretranslate](https://libretranslate.com).
- [x] ✍️ **Powerful subtitle editor** so you don't need to leave the UI!
  - Transcription highlighting based on media position.
  - CPS (Characters per second) warnings.
  - Segment splitting.
  - Segment insertion.
  - Subtitle language selection.
- [x] 🏠 **100% Local**: transcription, translation and subtitle edition happen 100% on your machine (can even work offline!).
- [x] 🚀 **Fast**: uses FasterWhisper as the Whisper backend: get much faster transcription times on CPU!
- [x] 👍 **Quick and easy setup**: use the quick start script, or run through a few steps!
- [x] 🔥 **GPU support**: use your NVIDIA GPU to get even faster transcription times!
- [x] 🐎 **CPU support**: no GPU? No problem! Whishper can run on CPU too.

### What's new in Whishper Reloaded

- [x] **Decoupled containers**: run backend, frontend and whishper-api as separate containers ([DockerHub frontend](https://hub.docker.com/r/thespartan94/whishper-frontend), [DockerHub backend](https://hub.docker.com/r/thespartan94/whishper-backend)). Useful if your powerful transcription machine is not always on and you need the service available at all times!
- [x] **Search bar** to find specific transcriptions.
- [x] **Better page loading** through pagination, avoiding blocking UI threads that made browsers lag.
- [x] **Feedback banners** on connections to available services, with a more fail-tolerant approach to translations and new transcriptions.
- [x] **Rename transcriptions** after a successful transcription.
- [x] **Import/Export via JSON**: download the current transcription as JSON, edit it externally for pre-processing, and reupload it for real-time editing.
- [x] **Real-time editing improvements**:
  - [x] Audio-only mode for the Whishper editor.
  - [x] Playback shortcuts [F7, F8, F9].
  - [x] Navigate through segments using TAB.
  - [x] Go to current segment / navigate to a specific segment number.

## Project structure

Whishper is a collection of pieces that work together. The main pieces are:

- **Transcription-API**: the API that runs Faster-Whisper. Found in the `transcription-api` folder.
- **Whishper-Backend**: the backend that coordinates frontend calls, database, and tasks. Found in the `backend` folder.
- **Whishper-Frontend**: the frontend (web UI) of the application. Found in the `frontend` folder.
- **Translation (3rd party)**: the LibreTranslate container used for translating subtitles.
- **MongoDB (3rd party)**: the database that stores all information about your transcriptions.
- **Nginx (3rd party)**: the proxy that allows running everything from a single domain.

## Running with Docker

Thanks to the decoupled architecture, you can run all the components on a single machine, or split them across several machines (for example, keeping the GPU-heavy `whishper-api` on a separate, occasionally-on machine while the rest of the stack stays always available).

### Docker Compose

The following `docker-compose.yml` runs the backend, MongoDB and frontend on the same server, while `whishper-api` is hosted on a different machine:

```yaml
services:
  whishper-mongo:
    image: mongo:latest
    container_name: whishper-mongo
    env_file:
      - .backend.env
    restart: unless-stopped
    volumes:
      - ./db_data/db:/data/db
      - ./db_data/logs/:/var/log/mongodb/
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${DB_USER:-whishper}
      MONGO_INITDB_ROOT_PASSWORD: ${DB_PASS:-whishper}
    expose:
      - 27017
    ports:
      - 27017:27017
    command: mongod --logpath var/log/mongodb/mongod.log

  whishper-frontend:
    image: thespartan94/whishper-frontend:latest
    container_name: whishper-frontend
    ports:
      - "3000:3000"
    expose:
      - 3000
    env_file:
      - .frontend.env
    volumes:
      - ./whishper-frontend/logs:/var/log/whishper
    depends_on:
      - whishper-backend

  whishper-backend:
    image: thespartan94/whishper-backend:latest
    container_name: whishper-backend
    ports:
      - 8080:8080
    env_file:
      - .backend.env
    volumes:
      - ./uploads:/uploads
      - ./whishper-backend/logs:/var/log/whishper
    depends_on:
      - whishper-mongo
```

### Environment files

The `docker-compose.yml` references two env files, one for the backend and one for the frontend.

`.backend.env`:

```env
UPLOAD_DIR=/uploads
ASR_ENDPOINT=<external-ip-address:8000> # assuming default port
DB_USER=whishper # if default
DB_PASS=whishper # if default
DB_ENDPOINT=whishper-mongo:27017
TRANSLATION_ENDPOINT=<external-ip-address:5000> # assuming default port
```

`.frontend.env`:

```env
PUBLIC_API_HOST=<external-public-api-host>
PUBLIC_TRANSLATION_API_HOST=<external-public-translation-api-host>
PUBLIC_INTERNAL_API_HOST=http://whishper-backend:8080
PUBLIC_WHISHPER_PROFILE=gpu # or cpu
```

### Start the stack

Once your env files are in place, bring everything up with:

```bash
docker compose up -d
```

The web UI will be available at `http://localhost:3000`.

## Screenshots

These screenshots are available on [the official website](https://whishper.net/usage/transcriptions/), click any of the following links to see:

- [A transcription creation](https://whishper.net/usage/transcriptions/)
- [A transcription translation](https://whishper.net/usage/translate/)
- [A transcription download](https://whishper.net/usage/download/)
- [The subtitle editor](https://whishper.net/usage/editor/)

## Development

### Contributing

Contributions are welcome! Feel free to open a PR with your changes, or take a look at the issues to see if there is something you can help with.

### Development setup

Check out the development documentation [here](https://whishper.net/guides/develop/).

## Star History

<a href="https://www.star-history.com/#devdema/whishper&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=devdema/whishper&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=devdema/whishper&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=devdema/whishper&type=Date" />
 </picture>
</a>

## Credits

- **[Whishper](https://github.com/pluja/whishper) by [Pluja](https://github.com/pluja)** — the original project and main contributor. Whishper Reloaded would not exist without it; please consider supporting the original idea.
- [Faster Whisper](https://github.com/guillaumekln/faster-whisper)
- [LibreTranslate](https://github.com/LibreTranslate/LibreTranslate)
- Maintained by [devdema](https://github.com/devdema).
