# GoBank Backend - Complete Deployment Setup ✅

## Overview

Your Go backend is now ready for deployment! The `.env` file has been created with all required configuration values. Follow this guide to deploy to Vercel and integrate with your frontend.

---

## 1. Environment Variables Setup

### Files Created
- ✅ `/backend/.env` - Local development configuration

### Required Environment Variables for Vercel Deployment

You need to set these in your Vercel project settings:

#### **Critical Variables (Required)**

```
DATABASE_URL=postgresql://user:password@host:port/dbname
JWT_SECRET=<your-32-char-random-secret>
```

#### **Optional Variables (with defaults)**

```
COUPON_CODE=OFFER1000
SMTP_EMAIL=ahnaf.shahriar2003@gmail.com
SMTP_PASSWORD=oiyo jwqw bflo cgwa
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
WEBAUTHN_RP_ORIGIN=https://your-frontend-domain.com
WEBAUTHN_RP_ID=your-domain.com
WEBAUTHN_DISPLAY_NAME=GoBank
```

---

## 2. Database Setup (Neon PostgreSQL)

### Option A: Using Neon PostgreSQL (Recommended)

1. **Create Neon Account**
   - Go to [neon.tech](https://neon.tech)
   - Sign up and create a project

2. **Create Database**
   - Create a new database named `gobank`
   - Copy your connection string

3. **Get Connection String**
   - Format: `postgresql://user:password@host:5432/gobank?sslmode=require`
   - This is your `DATABASE_URL`

4. **Create Required Tables** (auto-created on first run)
   - The backend will automatically create tables on startup:
     - `accounts` - User account data
     - `pending_accounts` - Email verification
     - `transfers` - Transaction records
     - `coupon_redemptions` - Offers tracking
     - `pending_profile_updates` - Profile change verification
     - `webauthn_credentials` - Passwordless auth data

### Option B: Local PostgreSQL (Development)

```bash
# Install PostgreSQL
# macOS
brew install postgresql

# Start PostgreSQL
brew services start postgresql

# Create database
createdb gobank

# Connection string
DATABASE_URL="postgresql://postgres:gobank@localhost:5432/gobank?sslmode=disable"
```

---

## 3. JWT Secret Generation

**Create a strong random secret:**

```bash
# On macOS/Linux
openssl rand -base64 32

# Or use this Python one-liner
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

**Current Secret (for local dev):**
```
NhDmIdl2BeGZbwklk8pBv2yIBZWsfrC8tI4x3o993aWwo=
```

⚠️ **For Production:** Generate a NEW secret and update in Vercel settings.

---

## 4. SMTP Configuration (Optional - for Emails)

### Gmail Setup

1. **Enable 2-Factor Authentication** in Google Account
2. **Generate App Password**
   - Go to [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)
   - Select "Mail" and "Windows Computer" (or your OS)
   - Copy the generated password
3. **Set Environment Variables**
   ```
   SMTP_EMAIL=your-email@gmail.com
   SMTP_PASSWORD=your-app-specific-password
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   ```

---

## 5. Vercel Deployment Steps

### Step 1: Connect Repository to Vercel

```bash
# Push your code to GitHub
git add .
git commit -m "Backend setup complete with .env configuration"
git push origin main
```

### Step 2: Add Project to Vercel

1. Go to [vercel.com/dashboard](https://vercel.com/dashboard)
2. Click "Add New..." → "Project"
3. Select your GitHub repository: `Ahnaf-2210007/go-bank-front-back-`
4. Select "Root Directory" as your backend folder

### Step 3: Set Environment Variables in Vercel

In Vercel Project Settings → Environment Variables, add:

```
DATABASE_URL=postgresql://user:password@host:5432/gobank?sslmode=require
JWT_SECRET=<your-new-32-char-secret>
COUPON_CODE=OFFER1000
SMTP_EMAIL=ahnaf.shahriar2003@gmail.com
SMTP_PASSWORD=oiyo jwqw bflo cgwa
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
WEBAUTHN_RP_ORIGIN=https://your-frontend-domain.vercel.app
WEBAUTHN_RP_ID=your-frontend-domain.vercel.app
WEBAUTHN_DISPLAY_NAME=GoBank
```

### Step 4: Configure Build Settings

Vercel should auto-detect from `vercel.json`:

- **Build Command:** `go build -o api .`
- **Install Command:** `go mod download`
- **Output Directory:** `.`
- **Framework:** Go

### Step 5: Deploy

Click "Deploy" and wait for completion. Your backend will be live at:
```
https://your-project-name.vercel.app
```

---

## 6. Frontend Integration

### Update Frontend Environment Variables

In your **frontend** `.env.local` or Vercel settings:

```
NEXT_PUBLIC_API_URL=https://your-project-name.vercel.app
```

### Update WebAuthn Configuration

In `frontend/.env`:
```
NEXT_PUBLIC_WEBAUTHN_RP_ID=your-frontend-domain.vercel.app
NEXT_PUBLIC_WEBAUTHN_ORIGIN=https://your-frontend-domain.vercel.app
```

### Example Frontend API Call

```typescript
// Authenticate with backend
const response = await fetch(
  `${process.env.NEXT_PUBLIC_API_URL}/login`,
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ account_number: 123, password: 'pass' })
  }
);
const data = await response.json();
const token = data.token; // Use for Authorization header
```

---

## 7. Testing Deployment

### Health Check Endpoint

```bash
curl https://your-project-name.vercel.app/health
# Expected response: {"status":"ok"}
```

### Test API Endpoints

```bash
# Create account
curl -X POST https://your-project-name.vercel.app/account \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "SecurePassword123"
  }'

# Login
curl -X POST https://your-project-name.vercel.app/login \
  -H "Content-Type: application/json" \
  -d '{
    "account_number": 123,
    "password": "SecurePassword123"
  }'
```

---

## 8. Troubleshooting

### ❌ Build Fails: "go: command not found"
- **Fix:** Vercel should auto-detect Go. Ensure `vercel.json` exists in backend root.

### ❌ Database Connection Error
- Check `DATABASE_URL` format is correct
- Verify database is accessible from Vercel
- For Neon: Check IP whitelist settings

### ❌ CORS Errors from Frontend
- Ensure backend is deployed and running
- Verify `/health` endpoint returns `{"status":"ok"}`
- Check frontend is calling correct `API_URL`

### ❌ WebAuthn Not Working
- Verify `WEBAUTHN_RP_ID` matches your domain (no `https://`)
- Verify `WEBAUTHN_RP_ORIGIN` includes protocol

### ❌ Email Not Sending
- Check SMTP credentials are correct
- Verify Gmail has 2FA enabled and app password created
- Check email address in `SMTP_EMAIL` is correct

---

## 9. Security Checklist

- ✅ `.env` file created with secrets
- ✅ `.gitignore` prevents `.env` from being committed
- ⚠️ **TODO:** Generate new `JWT_SECRET` for production
- ⚠️ **TODO:** Use app-specific password for Gmail
- ⚠️ **TODO:** Restrict CORS to your frontend domain (update in `api.go`)
- ⚠️ **TODO:** Enable HTTPS everywhere
- ⚠️ **TODO:** Set up error tracking (Sentry, Rollbar)

---

## 10. File Structure

```
backend/
├── .env                    ✅ Environment variables (created)
├── .env.example           📋 Template
├── .gitignore             🔒 Prevents .env commit
├── vercel.json            ⚙️  Deployment config
├── go.mod                 📦 Dependencies
├── main.go                🚀 Entry point
├── config.go              ⚙️  Configuration
├── api.go                 🔌 API endpoints
├── storage.go             💾 Database layer
├── types.go               📝 Data models
└── webauthn.go            🔐 Passwordless auth
```

---

## 11. Quick Deployment Command

```bash
# From project root
cd backend
git add .env
git commit -m "Add backend environment configuration"
git push origin main

# Then deploy via Vercel dashboard or CLI:
# vercel deploy --prod
```

---

## 12. Next Steps

1. **Set Database:** Add `DATABASE_URL` to Vercel settings
2. **Generate JWT Secret:** Create new secret for production
3. **Deploy Backend:** Push to GitHub and trigger Vercel deployment
4. **Update Frontend:** Point to backend domain
5. **Test APIs:** Verify endpoints work
6. **Monitor:** Set up error tracking

---

## 13. API Endpoints Reference

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/health` | ❌ | Health check |
| POST | `/account` | ❌ | Create account |
| POST | `/login` | ❌ | Login |
| POST | `/account/verification` | ❌ | Verify email |
| GET | `/account/{id}` | ✅ | Get account details |
| DELETE | `/account/{id}` | ✅ | Delete account |
| POST | `/account/update` | ✅ | Update profile |
| POST | `/transfer` | ✅ | Transfer funds |
| GET | `/account/transactions` | ✅ | Transaction history |
| POST | `/webauthn/register/begin` | ❌ | WebAuthn registration |
| POST | `/webauthn/login/begin` | ❌ | WebAuthn login |

---

## 14. Support & Resources

- **Vercel Go Support:** https://vercel.com/docs/frameworks/go
- **Neon Database:** https://neon.tech/docs
- **Go Module Management:** https://golang.org/doc/modules/managing-dependencies
- **JWT Guide:** https://jwt.io
- **WebAuthn Spec:** https://www.w3.org/TR/webauthn-2/

---

**Status:** ✅ Backend configuration complete and ready for deployment!

Generated: 2026-07-20
