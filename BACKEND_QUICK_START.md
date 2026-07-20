# Backend Quick Start Guide 🚀

## In 5 Minutes

### 1. Set Database URL (CRITICAL)
```
Vercel Dashboard → Settings → Environment Variables → Add:
DATABASE_URL = postgresql://user:pass@host:port/gobank
```

Choose one:
- **Neon** (Recommended): https://console.neon.tech → Create project → Copy connection string
- **Local**: `postgresql://postgres:gobank@localhost:5432/gobank?sslmode=disable`

### 2. Generate JWT Secret
```bash
openssl rand -base64 32
# Copy output to Vercel: JWT_SECRET = <output>
```

### 3. Add Other Variables to Vercel
```
COUPON_CODE=OFFER1000
SMTP_EMAIL=ahnaf.shahriar2003@gmail.com
SMTP_PASSWORD=oiyo jwqw bflo cgwa
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
WEBAUTHN_RP_ORIGIN=https://your-frontend-domain.com
WEBAUTHN_RP_ID=your-frontend-domain.com
WEBAUTHN_DISPLAY_NAME=GoBank
```

### 4. Deploy
```bash
git add .
git commit -m "Backend setup complete"
git push origin main
# Vercel auto-deploys in ~2 minutes
```

---

## Development Locally

### Start Backend
```bash
cd backend
go run .
# Server runs on http://localhost:3000
```

### Test It Works
```bash
curl http://localhost:3000/health
# Response: {"status":"ok"}
```

### Run Full Test Suite
```bash
bash backend/test-endpoints.sh
```

---

## Connect Frontend

Update frontend `.env`:
```
NEXT_PUBLIC_API_URL=https://your-backend.vercel.app
NEXT_PUBLIC_WEBAUTHN_RP_ID=your-domain.com
NEXT_PUBLIC_WEBAUTHN_ORIGIN=https://your-domain.com
```

---

## Environment Variables Status

✅ **Created:** `/backend/.env`

| Variable | Status | Value |
|----------|--------|-------|
| DATABASE_URL | 🔴 **TODO** | Set in Vercel |
| JWT_SECRET | ✅ Set | (for local dev) |
| COUPON_CODE | ✅ Set | OFFER1000 |
| SMTP_EMAIL | ✅ Set | ahnaf.shahriar2003@gmail.com |
| SMTP_PASSWORD | ✅ Set | oiyo jwqw bflo cgwa |
| SMTP_HOST | ✅ Set | smtp.gmail.com |
| SMTP_PORT | ✅ Set | 587 |
| WEBAUTHN_* | ✅ Set | localhost (dev) |

---

## Key Files

```
backend/.env                              # ← Configuration (created)
backend/.env.example                      # Template
BACKEND_DEPLOYMENT_COMPLETE.md            # Full guide
BACKEND_ENV_SETUP_CHECKLIST.md            # Variable details
BACKEND_SETUP_SUMMARY.md                  # Complete overview
```

---

## Top 3 Common Issues

| Issue | Solution |
|-------|----------|
| "DATABASE_URL not set" | Add to Vercel settings, redeploy |
| CORS errors from frontend | Verify `/health` works, check API_URL |
| WebAuthn not working | Check RP_ID (no https), RP_ORIGIN (with https) |

---

## API Endpoints

```
Public:
GET    /health                    # Health check
POST   /account                   # Create account
POST   /login                     # Login
POST   /account/verification      # Verify email
POST   /webauthn/register/begin
POST   /webauthn/login/begin

Protected (need JWT):
GET    /account/{id}              # Get account
POST   /account/update            # Update profile
POST   /transfer                  # Send money
GET    /account/transactions      # History
```

---

## Quick Deploy Checklist

- [ ] Set DATABASE_URL in Vercel
- [ ] Generate new JWT_SECRET
- [ ] Add WEBAUTHN variables for production domain
- [ ] `git push origin main`
- [ ] Wait 2 minutes for deployment
- [ ] Test `/health` endpoint
- [ ] Update frontend API_URL

---

## Support Docs

1. **Full Details:** `BACKEND_DEPLOYMENT_COMPLETE.md`
2. **Setup Steps:** `BACKEND_ENV_SETUP_CHECKLIST.md`
3. **Summary:** `BACKEND_SETUP_SUMMARY.md`

---

**Status:** ✅ Ready for deployment!
