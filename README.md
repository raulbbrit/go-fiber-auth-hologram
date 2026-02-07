# 🔮 Hologram Auth

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?style=for-the-badge&logo=go&logoColor=white)](https://gofiber.io)
[![GORM](https://img.shields.io/badge/GORM-SQLite-336791?style=for-the-badge&logo=sqlite&logoColor=white)](https://gorm.io)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.x-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)](https://tailwindcss.com)
[![Alpine.js](https://img.shields.io/badge/Alpine.js-3.x-8BC0D0?style=for-the-badge&logo=alpine.js&logoColor=white)](https://alpinejs.dev)
![Render](https://img.shields.io/badge/Render-Deployed-46E3B7?style=for-the-badge&logo=render&logoColor=white)

**Sci-fi holographic projection authentication system with immersive 3D effects.** Features realistic hologram visuals with scan lines, glitch animations, projection beams, and floating particle effects.

## ✨ Features

- 🎭 **Holographic Projection Effects** - Realistic hologram with projection beam and pedestal
- 📺 **Scan Line Animation** - Authentic CRT-style scanning effect
- ⚡ **Glitch Effects** - Random digital glitches during loading states
- ✨ **Floating Particles** - Ambient holographic particles rising upward
- 🌐 **Grid Floor** - Tron-style perspective grid animation
- 💎 **3D Floating Forms** - Smooth levitation animation with depth
- 🔒 **Real-time Validation** - Async email checking and password strength meter
- 🎨 **Cyan Hologram Palette** - Cohesive sci-fi color scheme throughout
- 📱 **Responsive Design** - Works beautifully on all devices
- 🚀 **Zero Page Reloads** - All errors handled via Alpine.js

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.21+, Fiber v2 |
| **Database** | GORM + Pure Go SQLite (no CGO) |
| **Frontend** | Alpine.js 3.x, Tailwind CSS |
| **Auth** | bcrypt password hashing, Sessions |

> 💡 Uses `github.com/glebarez/sqlite` - pure Go SQLite driver, works on Windows without CGO!

## 🚀 Quick Start

Clone the repository:

```bash
git clone https://github.com/smart-developer1791/go-fiber-auth-hologram
cd go-fiber-auth-hologram
```

Initialize dependencies and run:

```bash
go mod tidy
go run .
```

Open browser:

```text
http://localhost:3000
```

## 📱 Demo Account

```text
Email: demo@hologram.io
Password: demo123
Phone: +1 (555) 123-4567
```

## 📁 Project Structure

```text
go-fiber-auth-hologram/
├── main.go              # Fiber server, routes, handlers
├── go.mod               # Go module definition
├── render.yaml          # Render deployment config
├── .gitignore           # Git ignore rules
├── README.md            # This file
└── templates/
    ├── login.html       # Hologram login page
    ├── register.html    # Hologram registration page
    └── dashboard.html   # User dashboard
```

## 🎨 Visual Effects

### Hologram Projection
- Projection beam emanating from pedestal
- Flickering hologram animation
- Scan line sweeping effect
- Static noise overlay

### 3D Elements
- Floating form with perspective
- Corner bracket decorations
- Pulsing glow effects
- Ambient particle system

### Interactive States
- Glitch effect on form submission
- Real-time field validation colors
- Password strength visualization
- Loading spinners

## 🔐 Validation Features

### Email Validation
- Format verification
- Real-time availability check
- Visual status indicators

### Password Strength
- 6-level strength meter
- Requirement hints
- Color-coded feedback

### Form Security
- bcrypt hashing
- Session-based auth
- CSRF protection ready

## 🌐 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/login` | Login page |
| POST | `/login` | Authenticate user |
| GET | `/register` | Registration page |
| POST | `/register` | Create account |
| GET | `/dashboard` | Protected dashboard |
| POST | `/logout` | End session |
| POST | `/api/validate/email` | Check email availability |
| POST | `/api/validate/password` | Get password strength |

## 🎯 Why This Design?

The holographic projection concept represents:

- **Futuristic Security** - Implies advanced protection
- **Immersive Experience** - Engages users visually
- **Memorable Identity** - Unique brand differentiation
- **Technical Sophistication** - Suggests quality engineering

---

## Deploy in 10 seconds

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)
