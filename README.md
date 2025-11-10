<h1 align="center">🔐 DevMess — Secure Realtime Chat Application</h1>

<p align="center">
  <img src="https://img.shields.io/badge/build-passing-brightgreen?style=for-the-badge" alt="build">
  <img src="https://img.shields.io/badge/frontend-React-blue?style=for-the-badge&logo=react" alt="react">
  <img src="https://img.shields.io/badge/backend-Golang-00ADD8?style=for-the-badge&logo=go" alt="golang">
  <img src="https://img.shields.io/badge/realtime-WebSocket-orange?style=for-the-badge&logo=websocket" alt="ws">
  <img src="https://img.shields.io/badge/storage-DigitalOcean%20Spaces-0069ff?style=for-the-badge&logo=digitalocean" alt="do">
</p>

<p align="center">
  🚀 <b>DevMess</b> is a secure real-time chat platform featuring end-to-end encryption,  
  lightning-fast performance, and a sleek developer-friendly interface ⚡
  <link>https://devmess.cloud</link>⚡
</p>

---

## ✨ Features

- 💬 **Realtime Chat** — instant message updates powered by WebSocket  
- 🔐 **End-to-End Encryption** — AES + JWT authentication  
- 🎨 **Modern UI/UX** — built with React + Framer Motion animations  
- ☁️ **Media Uploads** — securely store images & videos with DigitalOcean Spaces  
- 🔔 **Message Notifications**  
- 👀 **Online/Offline Presence** tracking  
- 🗑️ **Message Deletion & Recall**  
- 🌍 **Multi-Device & Cross-Platform** ready  

---

## 🧱 System Architecture

---

## 📂 Folder Structure

```bash
devmess/
├── backend/
│   ├── main.go
│   ├── handler/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── utils/
│
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   ├── components/
│   │   ├── services/
│   │   ├── hooks/
│   │   └── utils/
│   └── public/
│
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md
```
---

## 🔐 Example .env configuration

---
```bash
# DATABASE
DB_HOST=localhost
DB_USER=postgres
DB_PASS=123456
DB_NAME=devmess
DB_PORT=5432

# JWT
JWT_SECRET=supersecretkey

# DIGITALOCEAN SPACES
SPACES_KEY=DO_SPACES_ACCESS_KEY
SPACES_SECRET=DO_SPACES_SECRET_KEY
SPACES_REGION=sgp1
SPACES_BUCKET=devmess
SPACES_ENDPOINT=https://sgp1.digitaloceanspaces.com
```
---
🐳 Docker Deployment
---
```bash
docker-compose up --build
```
