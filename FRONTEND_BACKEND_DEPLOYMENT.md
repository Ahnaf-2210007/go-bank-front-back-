# Frontend & Backend Deployment Guide

## Current Status

- **Backend**: Go serverless handler at `/api/handler.go` ✅
- **Frontend**: Next.js 16 at `/frontend` ✅
- **Connection**: Ready (API URL configured)

## Complete Deployment Steps

### Phase 1: Verify Backend is Working (5 minutes)

1. **Check if backend is deployed**:
   ```bash
   curl https://your-backend-domain.vercel.app/api/health
   # Should return: {"status":"ok"}
   ```

2. **If it returns 404**, your backend hasn't been deployed yet:
   - Go to Vercel → Your Project → Deployments
   - Check if the latest deployment is successful
   - If failed, check logs and environment variables

### Phase 2: Set Frontend Environment Variables (5 minutes)

1. **Get your backend URL**:
   - Go to Vercel Dashboard
   - Find your GoBank project
   - Copy the production domain from the URL bar
   - Format: `https://go-bank-front-back.vercel.app` (example, replace with yours)

2. **Add to Vercel Environment Variables**:
   - Settings → Environment Variables
   - Add these variables:
     ```
     NEXT_PUBLIC_API_URL=https://your-backend-domain.vercel.app/api
     NEXT_PUBLIC_BACKEND_URL=https://your-backend-domain.vercel.app
     ```

3. **Click "Add to Production"** to apply to production builds

### Phase 3: Deploy Frontend (5 minutes)

1. **Ensure changes are committed**:
   ```bash
   cd /vercel/share/v0-project
   git add .
   git commit -m "Frontend environment configuration"
   git push origin main
   ```

2. **Watch Vercel deploy**:
   - Go to Vercel Dashboard
   - Wait for deployment to complete (2-3 minutes)
   - You'll see a checkmark when done

3. **Test the deployment**:
   - Visit your frontend URL: `https://your-frontend-domain.vercel.app`
   - You should see the login page or dashboard redirect
   - If you see 404, the build may have failed

### Phase 4: Test End-to-End (10 minutes)

1. **Test health check**:
   ```bash
   curl https://your-backend-domain.vercel.app/api/health
   ```

2. **Test frontend loads**:
   - Open `https://your-frontend-domain.vercel.app` in browser
   - You should see GoBank login page or dashboard

3. **Test login flow**:
   - Create a test account (register)
   - Login with credentials
   - You should see dashboard
   - Check if API calls work in browser console

### Phase 5: Troubleshoot (if needed)

**Frontend shows 404**:
- Check Vercel deployment status
- Check build logs: Vercel → Deployments → Click latest → Build logs
- Verify environment variables are set

**Frontend won't connect to backend**:
- Check NEXT_PUBLIC_API_URL in Vercel environment variables
- Open browser DevTools → Network tab
- Look for failed API calls
- Check the URL being called matches your backend

**Backend returns 500 error**:
- Check Vercel backend deployment logs
- Verify all backend environment variables are set:
  - DATABASE_URL
  - JWT_SECRET
  - SMTP_EMAIL, SMTP_PASSWORD
  - etc.

## Environment Variables Summary

### Backend (already set):
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_NAME`
- `DB_PASSWORD`
- `JWT_SECRET`
- `LISTEN_ADDR`
- `SMTP_EMAIL`
- `SMTP_PASSWORD`
- `SMTP_HOST`
- `SMTP_PORT`
- `COUPON_CODE`

### Frontend (need to set):
- `NEXT_PUBLIC_API_URL` → `https://your-backend.vercel.app/api`
- `NEXT_PUBLIC_BACKEND_URL` → `https://your-backend.vercel.app`

## Quick Commands

```bash
# Commit and push changes
git add .
git commit -m "Deploy frontend and backend"
git push origin main

# Check backend health
curl https://your-backend.vercel.app/api/health

# Check logs
vercel logs --tail
```

## API Endpoints Available

Once connected, frontend can call:

**Public endpoints** (no auth needed):
- `POST /api/register` - Register new user
- `POST /api/login` - Login
- `POST /api/webauthn/login/begin` - Start passkey login
- `POST /api/webauthn/login/finish` - Complete passkey login
- `POST /api/verify/email` - Verify email with OTP
- `GET /api/health` - Health check

**Protected endpoints** (JWT token required):
- `GET /api/accounts` - Get account info
- `POST /api/transfer` - Transfer funds
- `GET /api/transactions` - Get transaction history
- `POST /api/offers` - Apply coupon
- And 11 more...

## Success Criteria

You'll know it's working when:

✅ Frontend deployed and accessible
✅ Backend health endpoint returns `{"status":"ok"}`
✅ Login page loads in frontend
✅ Can register and login successfully
✅ Dashboard shows account info from backend
✅ Transfers work between accounts

## Support

If something goes wrong:

1. **Check browser console** for error messages
2. **Check Vercel deployment logs** for build/runtime errors
3. **Verify environment variables** are correctly set
4. **Test health endpoint** to confirm backend is working
5. **Review this guide** for missed steps
