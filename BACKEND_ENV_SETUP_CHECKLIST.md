# Backend Environment Variables - Setup Checklist ✅

## Current Status

### ✅ Completed
- [x] `.env` file created in `/backend/` directory
- [x] All template variables from `.env.example` copied to `.env`
- [x] Deployment guide created
- [x] Environment variables documented

---

## Required Variables for Vercel Deployment

### 1. Database Configuration

**Variable:** `DATABASE_URL`

**Current Status:** ❌ NEEDS TO BE SET
**Priority:** 🔴 CRITICAL

**How to set:**

#### Option A: Neon PostgreSQL (Recommended)
```
Steps:
1. Go to https://console.neon.tech
2. Create project → create database "gobank"
3. Copy connection string
4. Add to Vercel: DATABASE_URL = postgresql://...
```

**Format Example:**
```
DATABASE_URL=postgresql://user:neon_password@ep-xxxxx.us-east-1.neon.tech:5432/gobank?sslmode=require
```

#### Option B: AWS RDS / Self-hosted PostgreSQL
```
Format:
DATABASE_URL=postgresql://user:password@host:5432/gobank?sslmode=disable
```

---

### 2. JWT Secret

**Variable:** `JWT_SECRET`

**Current Status:** ✅ SET (for local dev)
```
NhDmIdl2BeGZbwklk8pBv2yIBZWsfrC8tI4x3o993aWwo=
```

**For Production:** 🔴 GENERATE NEW SECRET

**Generate new secret:**
```bash
# macOS/Linux
openssl rand -base64 32

# Or online
# https://generate-random.org/
```

**Action Steps:**
1. Generate new 32-character secret
2. Update `.env` for local testing
3. Add to Vercel environment variables

---

## Optional Variables (Already Set with Defaults)

### 3. SMTP Configuration

**Status:** ✅ CONFIGURED

```
SMTP_EMAIL=ahnaf.shahriar2003@gmail.com
SMTP_PASSWORD=oiyo jwqw bflo cgwa
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

**Note:** These are for sending verification/notification emails. If not set, emails are logged to console instead.

⚠️ **Security Warning:** For production, use Gmail App Password, not your actual password.

**To set up Gmail App Password:**
1. Enable 2-Factor Authentication in Google Account
2. Go to https://myaccount.google.com/apppasswords
3. Select "Mail" and "Other (custom name)"
4. Copy generated password
5. Update Vercel settings

---

### 4. Coupon Code

**Status:** ✅ SET
```
COUPON_CODE=OFFER1000
```

---

### 5. WebAuthn Configuration

**Status:** ✅ SET (for local dev)

```
WEBAUTHN_RP_ORIGIN=http://localhost:8080
WEBAUTHN_RP_ID=localhost
WEBAUTHN_DISPLAY_NAME=GoBank
```

**For Production:** Update these values

**When deploying to production:**
- `WEBAUTHN_RP_ID`: Your domain WITHOUT protocol (e.g., `example.com`)
- `WEBAUTHN_RP_ORIGIN`: Full URL WITH protocol (e.g., `https://example.com`)
- `WEBAUTHN_DISPLAY_NAME`: Your app name

---

## Vercel Environment Variables Form

Use this template when adding environment variables to your Vercel project:

```
Key                     Value
─────────────────────────────────────────────────────────────────
DATABASE_URL            [YOUR_NEON_CONNECTION_STRING]
JWT_SECRET              [YOUR_NEW_32_CHAR_SECRET]
COUPON_CODE             OFFER1000
SMTP_EMAIL              ahnaf.shahriar2003@gmail.com
SMTP_PASSWORD           oiyo jwqw bflo cgwa
SMTP_HOST               smtp.gmail.com
SMTP_PORT               587
WEBAUTHN_RP_ORIGIN      https://your-frontend-domain.vercel.app
WEBAUTHN_RP_ID          your-frontend-domain.vercel.app
WEBAUTHN_DISPLAY_NAME   GoBank
```

---

## Step-by-Step Vercel Setup

### 1. Go to Vercel Project Settings

```
https://vercel.com/dashboard → Select project → Settings → Environment Variables
```

### 2. Add Variables

Click "Add new" and fill in:

| Key | Value | Environments |
|-----|-------|--------------|
| DATABASE_URL | postgresql://... | Production, Preview, Development |
| JWT_SECRET | your-32-char-secret | Production, Preview, Development |
| Others | (see form above) | (Select as needed) |

### 3. Redeploy

After adding variables, redeploy:
```bash
git push origin main
# or
vercel deploy --prod
```

---

## Files Modified/Created

```
backend/
├── .env                                    ✅ Created (local configuration)
└── .env.example                            ✓ Already exists (template)

root/
├── BACKEND_DEPLOYMENT_COMPLETE.md          ✅ Created (full guide)
└── BACKEND_ENV_SETUP_CHECKLIST.md         ✅ You are here
```

---

## Testing Environment Variables

### Local Testing

```bash
cd backend

# Verify .env exists
ls -la .env

# Check Go can read variables
go run . # Should load .env automatically
```

### After Deployment to Vercel

```bash
# Test health endpoint
curl https://your-vercel-project.vercel.app/health

# Expected response:
# {"status":"ok"}
```

---

## Common Issues & Fixes

### ❌ "DATABASE_URL not set"

**Fix:**
1. Verify DATABASE_URL is in Vercel Environment Variables
2. Redeploy with `git push origin main`
3. Check variable is set for correct environment (Production/Preview/Development)

### ❌ "JWT token invalid"

**Fix:**
1. Ensure same JWT_SECRET in frontend and backend
2. Token must not be expired
3. Check Authorization header format: `Bearer <token>`

### ❌ "CORS error from frontend"

**Fix:**
1. Verify backend is running: `curl /health`
2. Check frontend API_URL is correct
3. Ensure CORS headers are present in backend (auto-enabled)

### ❌ "WebAuthn not working"

**Fix:**
1. Verify WEBAUTHN_RP_ID matches domain exactly (no protocol)
2. Verify WEBAUTHN_RP_ORIGIN includes protocol
3. Check browser dev tools for errors
4. Must use HTTPS in production

---

## Priority Action Items

### 🔴 CRITICAL (Do Now)
- [ ] Set `DATABASE_URL` in Vercel
  - Create Neon account → database → get connection string
  
### 🟠 IMPORTANT (Before First Deploy)
- [ ] Generate new `JWT_SECRET` for production
- [ ] Update `WEBAUTHN_RP_ORIGIN` and `WEBAUTHN_RP_ID` for production domain
- [ ] Create Gmail App Password and update `SMTP_PASSWORD`

### 🟡 RECOMMENDED (For Production)
- [ ] Set up error tracking (Sentry)
- [ ] Enable database backups
- [ ] Monitor API performance
- [ ] Set up rate limiting

---

## Reference: Complete Environment Variables Map

```
┌─────────────────────────────────────────────┐
│        BACKEND ENVIRONMENT VARIABLES        │
├─────────────────────────────────────────────┤
│                                             │
│  DATABASE                                   │
│  ├─ DATABASE_URL (required)                │
│  ├─ DB_HOST (fallback)                     │
│  ├─ DB_PORT (fallback)                     │
│  ├─ DB_USER (fallback)                     │
│  ├─ DB_NAME (fallback)                     │
│  └─ DB_PASSWORD (fallback)                 │
│                                             │
│  SECURITY                                   │
│  └─ JWT_SECRET (required)                  │
│                                             │
│  COMMUNICATION                              │
│  ├─ SMTP_EMAIL (optional)                  │
│  ├─ SMTP_PASSWORD (optional)               │
│  ├─ SMTP_HOST (optional)                   │
│  └─ SMTP_PORT (optional)                   │
│                                             │
│  BUSINESS LOGIC                             │
│  ├─ COUPON_CODE (optional)                 │
│  ├─ LISTEN_ADDR (optional)                 │
│                                             │
│  WEB AUTHENTICATION                         │
│  ├─ WEBAUTHN_RP_ORIGIN (optional)          │
│  ├─ WEBAUTHN_RP_ID (optional)              │
│  └─ WEBAUTHN_DISPLAY_NAME (optional)       │
│                                             │
└─────────────────────────────────────────────┘
```

---

## Success Criteria

Your backend setup is complete when:

- ✅ `.env` file exists in `/backend/`
- ✅ All required variables are set for Vercel deployment
- ✅ `DATABASE_URL` points to working PostgreSQL database
- ✅ `JWT_SECRET` is a strong random string
- ✅ Backend builds without errors: `go build -o api .`
- ✅ `/health` endpoint returns `{"status":"ok"}`
- ✅ Frontend can connect to backend APIs

---

**Last Updated:** 2026-07-20
**Status:** ✅ Ready for Vercel Deployment
