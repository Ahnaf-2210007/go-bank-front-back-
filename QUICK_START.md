# Quick Start - Backend Integration Guide

## 30-Second Summary

Your backend had a **CORS issue** that made it crash when called from the frontend. This is now **FIXED** ✅

**What was fixed:**
- ✅ Added CORS headers to all responses
- ✅ Added health check endpoint (`/health`)
- ✅ Made WebAuthn configuration dynamic (works in production)
- ✅ Created Vercel deployment config (`vercel.json`)

**Next steps:** Deploy to Vercel and connect your frontend.

---

## Deploy to Vercel (1 minute)

```bash
# 1. Commit changes
git add .
git commit -m "Backend fixes: CORS, WebAuthn config, health check"

# 2. Push to trigger deployment
git push

# 3. Wait for Vercel to deploy...
# Check status: https://vercel.com/dashboard
```

**After deployment, test:**
```bash
curl https://your-domain.vercel.app/health
# Should return: {"status":"ok"}
```

---

## Connect Your Frontend (5 minutes)

### Update API Base URL

**Development:**
```javascript
const API_BASE_URL = 'http://localhost:3000';
```

**Production:**
```javascript
const API_BASE_URL = 'https://your-vercel-domain.vercel.app';
```

### Use This Authentication Pattern

```javascript
// 1. Login
const response = await fetch(`${API_BASE_URL}/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    number: 1234567890,
    password: 'password'
  })
});

const { token } = await response.json();
localStorage.setItem('token', token);

// 2. Use token for authenticated requests
const accountResponse = await fetch(`${API_BASE_URL}/account/123`, {
  headers: { 'Authorization': `Bearer ${token}` }
});
```

---

## All API Endpoints

### Public Endpoints (No Token Needed)

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/health` | Check if backend is running |
| POST | `/login` | Login with account number + password |
| POST | `/account` | Create new account |
| POST | `/account/verification` | Verify email with OTP |
| POST | `/webauthn/register/begin` | Start WebAuthn registration |
| POST | `/webauthn/register/finish` | Complete WebAuthn registration |
| POST | `/webauthn/login/begin` | Start WebAuthn login |
| POST | `/webauthn/login/finish/{email}` | Complete WebAuthn login |

### Protected Endpoints (Token Required)

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/account/{id}` | Get account details |
| DELETE | `/account/{id}` | Delete account |
| POST | `/account/update` | Update profile/email/password |
| GET | `/account/transactions` | Get transaction history |
| POST | `/account/{id}/offer` | Apply coupon code |
| POST | `/transfer` | Transfer funds |

---

## Example API Calls

### 1. Check Backend is Running
```bash
curl http://localhost:3000/health
# {"status":"ok"}
```

### 2. Create Account
```bash
curl -X POST http://localhost:3000/account \
  -H "Content-Type: application/json" \
  -d '{
    "firstName":"John",
    "lastName":"Doe",
    "email":"john@example.com",
    "password":"SecurePass123"
  }'
```

### 3. Login
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{
    "number":1234567890,
    "password":"SecurePass123"
  }'
# Returns: {"token":"...", "number":...}
```

### 4. Get Account (With Token)
```bash
TOKEN="your-jwt-token-here"

curl http://localhost:3000/account/123 \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Transfer Funds
```bash
TOKEN="your-jwt-token-here"

curl -X POST http://localhost:3000/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "toAccount":9876543210,
    "amount":1000
  }'
```

---

## Environment Variables

### For Vercel Project Settings

Add these in **Project Settings → Environment Variables:**

**Required:**
- `DATABASE_URL` - From Neon PostgreSQL
- `JWT_SECRET` - Run: `openssl rand -base64 32`

**Optional (WebAuthn):**
- `WEBAUTHN_RP_ORIGIN` - Your app URL (e.g., `https://myapp.vercel.app`)
- `WEBAUTHN_RP_ID` - Domain only (e.g., `myapp.vercel.app`)
- `WEBAUTHN_DISPLAY_NAME` - Display name (e.g., `GoBank`)

**Optional (Email):**
- `SMTP_EMAIL` - Your Gmail
- `SMTP_PASSWORD` - Gmail app password

---

## What Changed

### 1. CORS Middleware Added
- All responses now include CORS headers
- OPTIONS preflight requests are handled
- Frontend can call backend from any domain

### 2. Health Check Endpoint
- Added `/health` endpoint
- Returns `{"status": "ok"}`
- Useful for monitoring

### 3. WebAuthn Configurable
- Before: Hardcoded to `localhost:8080`
- Now: Uses environment variables
- Works in any environment (local, production, etc.)

### 4. Production Config Added
- Created `vercel.json`
- Tells Vercel how to build and run the backend
- Sets up timeout and memory limits

---

## Files to Read

| File | Read This If... |
|------|-----------------|
| `README_BACKEND_FIXES.md` | You want complete overview of all changes |
| `BACKEND_SETUP.md` | You need detailed backend documentation |
| `FRONTEND_INTEGRATION.md` | You're building the frontend and need examples |
| `ENV_SETUP.md` | You need to set up environment variables |
| `CHANGES.md` | You want technical details of what changed |

---

## Common Issues

### CORS Error in Frontend
```
Access to XMLHttpRequest has been blocked by CORS policy
```

**Solution:**
1. Check backend is running: `curl http://localhost:3000/health`
2. Check frontend is using correct API URL
3. Verify backend is deployed (if using production URL)

### WebAuthn Not Working
**Solution:**
1. Check `WEBAUTHN_RP_ORIGIN` matches your domain exactly
2. Check `WEBAUTHN_RP_ID` is domain only (no protocol)

### "Not Authorized" Error
**Solution:**
1. Make sure you're including the token in Authorization header
2. Format: `Authorization: Bearer <your-token>`
3. Copy entire token exactly (no extra spaces)

### Database Connection Error
**Solution:**
1. Verify `DATABASE_URL` is set in environment
2. Check connection string is from Neon (not localhost)
3. Verify Neon database is running

---

## Deployment Steps

### Step 1: Make Sure Everything is Committed

```bash
git status
# Should show nothing to commit
```

### Step 2: Push to Repository

```bash
git push origin main
```

### Step 3: Vercel Deploys Automatically

- Vercel detects changes
- Runs build command from `vercel.json`
- Deploys new version
- Takes 1-2 minutes

### Step 4: Verify Deployment

```bash
curl https://your-domain.vercel.app/health
```

---

## Quick Test Checklist

- [ ] Backend runs locally: `go run .`
- [ ] Health check works: `curl http://localhost:3000/health`
- [ ] Can create account: POST `/account`
- [ ] Can login: POST `/login`
- [ ] Can use token: GET `/account/{id}` with Authorization header
- [ ] Deployed to Vercel: `git push`
- [ ] Production health check works: `curl https://your-domain/health`
- [ ] Frontend can call backend (no CORS errors)

---

## Frontend Template

```javascript
// api.js or api.ts

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

export async function api(endpoint, options = {}) {
  const token = localStorage.getItem('token');
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  };
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers
  });
  
  if (!response.ok) {
    const { error } = await response.json();
    throw new Error(error || 'Request failed');
  }
  
  return response.json();
}

// Usage
export const login = (number, password) => 
  api('/login', {
    method: 'POST',
    body: JSON.stringify({ number, password })
  });

export const getAccount = (id) => 
  api(`/account/${id}`);

export const transfer = (toAccount, amount) =>
  api('/transfer', {
    method: 'POST',
    body: JSON.stringify({ toAccount, amount })
  });
```

---

## Production Checklist

- [ ] All required environment variables set in Vercel
- [ ] Database connection verified
- [ ] JWT_SECRET set and strong (32+ chars)
- [ ] WebAuthn configured for production domain
- [ ] Backend deployed and health check working
- [ ] Frontend API URL updated for production
- [ ] CORS origin allowed (or restricted appropriately)
- [ ] SMTP configured (optional, for emails)
- [ ] Error logging set up (optional)
- [ ] Monitoring enabled (optional)

---

## Support

**Something not working?**

1. Check the health endpoint: `curl https://your-domain.vercel.app/health`
2. Check Vercel logs: `vercel logs --tail`
3. Review error messages in response
4. Read the appropriate documentation file (see "Files to Read" section above)

**Still stuck?**

1. Review `BACKEND_SETUP.md` troubleshooting section
2. Check `FRONTEND_INTEGRATION.md` for example API calls
3. Verify environment variables in `ENV_SETUP.md`

---

## That's It!

Your backend is ready. Now:
1. ✅ Deploy to Vercel (`git push`)
2. ✅ Connect your frontend (use examples above)
3. ✅ Test the integration
4. ✅ You're done! 🚀

---

**Need detailed info?** See the full documentation files in the project root.
