# Backend Connection Fix

## Problem
You were seeing "Registration error - Load failed" on the signup page because the frontend was trying to call the wrong URL.

## Root Cause
The `NEXT_PUBLIC_API_URL` environment variable was incorrectly pointing to:
```
https://go-bank-front-back.vercel.app/api  ❌ (FRONTEND - WRONG)
```

Instead of:
```
https://go-bank-front-back-ivory.vercel.app  ✅ (BACKEND - CORRECT)
```

## What Was Happening
1. Frontend receives registration form submission
2. Frontend tries to POST to `/account` endpoint
3. Frontend constructs URL: `https://go-bank-front-back.vercel.app/api/account`
4. This URL points to the FRONTEND itself, not the backend
5. Frontend doesn't have an `/account` endpoint, so it fails with "Load failed"

## Solution Applied
Updated both configuration files to point to the correct backend:

### vercel.json
```json
{
  "env": {
    "NEXT_PUBLIC_API_URL": "https://go-bank-front-back-ivory.vercel.app",
    "NEXT_PUBLIC_BACKEND_URL": "https://go-bank-front-back-ivory.vercel.app"
  }
}
```

### How It Works Now
1. Frontend POST to registration: `POST https://go-bank-front-back-ivory.vercel.app/account`
2. Request reaches the **backend** server
3. Backend processes account creation
4. Backend returns success response
5. Frontend redirects to verification page

## Verification
After redeploying, you should be able to:
- Register new accounts successfully
- See the account creation success message
- Be redirected to email verification

## Domain Reference
- **Frontend**: https://go-bank-front-back.vercel.app
- **Backend**: https://go-bank-front-back-ivory.vercel.app
- **API Endpoint**: https://go-bank-front-back-ivory.vercel.app/account (for registration)
