**Whishper Reloaded** is an open-source, 100% local audio transcription and subtitling suite with a full-featured web UI. It is based on [Pluja's Whishper](https://github.com/pluja/whishper) and further improved.

## What's different?

- [x] Load whishper backend, frontend and whishper-api as separate containers ([DockerHub frontend](https://hub.docker.com/r/thespartan94/whishper-frontend), [DockerHub backend](https://hub.docker.com/r/thespartan94/whishper-backend)). Useful if your powerful transcriptive is not always on and you need the service at all times!
- [x] Search for specific transcriptions via the search bar!
- [x] Better page loading experience using pagination, solving blocking UI threads making browsers lag.
- [x] Feedback banners on connections to available services, having a more fail-prone approach to translations and new transcriptions services.
- [x] Rename your transcriptions after a successful transcription.
- [x] Upload your whole transcription through JSON file for external pre-processing before real-time editing! Download the current transcription in JSON format, edit it externally, and reupload.  
- [x] Real-time editing improvements:
    - [x] Audio-only mode for the Whishper editor
    - [x] Use shortcut to control the playback [F7, F8, F9]
    - [x] Navigate through segments using TAB on your keyboard
    - [x] Go to current segment / Navigate to segment number  

## Docker-Compose file

The following is an example of a docker-compose.yml that puts backend, mongoDB and frontend on the same server, while whishper-api is hosted on a different machine:

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

The docker-compose.yml references two different env files, one for backend and another for frontend:

.backend.env:

    UPLOAD_DIR=/uploads
    ASR_ENDPOINT=<external-ip-address:8000> # assuming default port
    DB_USER=whishper # if default
    DB_PASS=whishper # if default
    DB_ENDPOINT=whishper-mongo:27017 
    TRANSLATION_ENDPOINT=<external-ip-address:5000> # assuming default port

.frontend.env:

    PUBLIC_API_HOST=<external-public-api-host>
    PUBLIC_TRANSLATION_API_HOST=<external-public-translation-api-host>
    PUBLIC_INTERNAL_API_HOST=http://whishper-backend:8080
    PUBLIC_WHISHPER_PROFILE=gpu # or cpu
    
# Original Whishper description

[![whishper banner](misc/banner.png)](https://whishper.net)

[![](https://img.shields.io/badge/website-066da5?style=for-the-badge&logo=icloud&logoColor=white)](https://whishper.net)
[![](https://img.shields.io/badge/self%20host%20guide-066da5?style=for-the-badge&logo=googledocs&logoColor=white)](https://whishper.net/guides/install)
[![](https://img.shields.io/badge/screenshots-5c1f87?style=for-the-badge&logo=slickpic&logoColor=white)](#screenshots)
[![](https://img.shields.io/docker/pulls/pluja/whishper?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/r/pluja/whishper)


## Features

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

## Project structure

Whishper is a collection of pieces that work together. The three main pieces are:

- Transcription-API: This is the API that enables running Faster-Whisper. You can find it in the `transcription-api` folder.
- Whishper-Backend: This is the backend that coordinates frontend calls, database, and tasks. You can find it in `backend` folder.
- Whishper-Frontend: This is the frontend (web UI) of the application. You can find it in `frontend` folder.
- Translation (3rd party): This is the libretranslate container that is used for translating subtitles.
- MongoDB (3rd party): This is the database that stores all the information about your transcriptions.
- Nginx (3rd party): This is the proxy that allows running everything from a single domain.

### Contributing

Contributions are welcome! Feel free to open a PR with your changes, or take a look at the issues to see if there is something you can help with.

### Development setup

Check out the development documentation [here](https://whishper.net/guides/develop/).

## Screenshots

These screenshots are available on [the official website](https://whishper.net/usage/transcriptions/), click any of the following links to see:

- [A transcription creation](https://whishper.net/usage/transcriptions/)
- [A transcription translation](https://whishper.net/usage/translate/)
- [A transcription download](https://whishper.net/usage/download/)
- [The subtitle editor](https://whishper.net/usage/editor/)

## Support

- [Monero](https://www.getmonero.org/): `82x6cn628oTUXV63DxBd6MJB8d997FhaSaGFvoWMgwihVmgiXYQPAwm2BCH31AovA9Qnnv1qQRrJk83TaJ8DaSZU2zkbWfM`
- [Bitcoin](https://bitcoin.org/en/): `bc1qfph44jl4cy03stwfkk7g0qlwx2grldr9xpk086`
- [Lightning Network (kycnotme)](https://getalby.com/p/kycnotme)

## Star History

<a href="https://www.star-history.com/#devdema/whishper&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=devdema/whishper&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=devdema/whishper&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=devdema/whishper&type=Date" />
 </picture>
</a>

## Credits

- [Faster Whisper](https://github.com/guillaumekln/faster-whisper)
- [LibreTranslate](https://github.com/LibreTranslate/LibreTranslate)


