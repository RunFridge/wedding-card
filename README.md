**English** | [한국어](docs/README_ko.md)

# Wedding Invitation Website

[![Release](https://img.shields.io/github/v/release/RunFridge/wedding-card)](https://github.com/RunFridge/wedding-card/releases)
[![License](https://img.shields.io/github/license/RunFridge/wedding-card)](LICENSE)
[![GitHub](https://img.shields.io/badge/GitHub-RunFridge%2Fwedding--card-181717?logo=github)](https://github.com/RunFridge/wedding-card)

A self-hostable Korean wedding invitation website with a cozy Stardew Valley pixel-art aesthetic. Your guests get an interactive invitation — guestbook, mini-game, photo sharing — and you get an admin panel to manage everything without touching code.

Everything runs from a single Docker container. No cloud accounts, no monthly fees, no external services required.

> **🌐 Live demo:** [demo_wedding.runfridge.dev](https://demo_wedding.runfridge.dev) — try it before hosting your own.
> ⚠️ It's a demo server: it may be down and its data may be reset at any time.

## What Your Guests Get

- **Invitation** — Couple and family introductions, date and venue, and bank transfer info for congratulatory gifts
- **Directions** — Venue map with Kakao Map, Naver Map, T Map, and Google Maps links, plus subway, bus, parking, and charter bus guidance
- **Guestbook** — Congratulatory messages guests can edit or delete with their own password, with an option to leave secret messages only you two can read
- **Match Card Game** — A photo-matching mini-game using your wedding photos, with an arcade-style leaderboard, unlockable achievements, and a Hall of Fame
- **Photo Sharing** — Guests upload photos from the wedding day with a built-in editor (crop, rotate, filters); uploads open shortly before the ceremony with a live countdown
- **Little Delights** — Day/night themes, an animated pixel-art scene with wandering farm animals, confetti celebrations, and haptic feedback on mobile

## What You Get

The admin panel at `/-/admin/` puts everything in your hands:

- **Wedding details** — Names, family info, date, venue, greetings, bank accounts, and directions, all edited in the browser — even after sending out the link
- **Photo management** — Curate the photos for the game and gallery, pick the main photo, and approve or hide guest uploads
- **Guestbook moderation** — Hide or delete messages, and read the secret ones
- **Optional auto-moderation** — Let OpenAI's moderation service screen guest photos and messages automatically; you can always override
- **Analytics** — Daily visitor chart and live server logs
- **Leaderboard control** — Remove joke entries or reset the game rankings

## Getting Started

You need a server (any small VPS or home server works) with [Docker](https://docs.docker.com/get-docker/) installed.

**1. Start it**

```bash
docker run -d --name wedding-server \
  -p 8080:8080 \
  -v wedding-data:/data \
  -e TZ=Asia/Seoul \
  --restart unless-stopped \
  ghcr.io/runfridge/wedding-card:latest
```

Swap `Asia/Seoul` for your timezone if the wedding isn't in Korea. Prefer Docker Compose? Download [`docker-compose.yml`](docker-compose.yml) from this repository and run `docker compose up -d`.

**2. Get your admin password**

A password is generated on first launch and printed in the server logs:

```bash
docker logs wedding-server
```

It's also saved to `/data/admin_password.txt` (deleted once you set your own).

**3. Set up**

Open `http://localhost:8080/-/admin/` and log in. A setup wizard walks you through choosing your own password, photo storage (the default keeps photos on your server), and optional content moderation.

**4. Make it yours**

Fill in your wedding details under **Settings**, upload your photos, and share the link with your guests. 🎉

> **Tip:** To serve the site on a domain with HTTPS, put a reverse proxy (Caddy, Nginx, Traefik, …) in front of port 8080.

## Your Data & Backups

Everything lives in the `wedding-data` Docker volume: the database and, with the default local storage, all uploaded photos. Backing up that one volume backs up your whole wedding site.

## Content Moderation (Optional)

Don't want to approve every guest photo by hand? Enable moderation in **Admin → System** (or during setup) with an OpenAI API key. Clean photos are approved automatically, inappropriate guestbook messages are hidden automatically, and every decision can be overridden from the admin panel.

## Offline / Air-gapped Install

Each [GitHub Release](https://github.com/RunFridge/wedding-card/releases) attaches image tarballs for servers without registry access:

```bash
docker load < wedding-server-<version>-amd64.tar.gz   # or -arm64
```

## Languages

The invitation and admin panel are available in Korean and English. The language follows each visitor's browser and can be switched anytime; your own content — greetings, notices, directions — is shown exactly as you wrote it.

## For Developers

Want to build from source, explore the API, add a language, or contribute? See the [Development Guide](docs/DEVELOPMENT.md).

## License

[MIT](LICENSE)
