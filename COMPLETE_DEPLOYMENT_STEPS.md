# Complete Frontend & Backend Deployment Steps

## Current Situation

Your project has both:
- **Backend**: Go serverless handler (ready)
- **Frontend**: Next.js 16 app (ready)
- **Issue**: 404 error because frontend is not deployed yet

## What You Need To Do Now

### Step 1: Merge v0 Branch to Main (on GitHub)

1. Go to GitHub: https://github.com/Ahnaf-2210007/go-bank-front-back-
2. You should see a notification "v0/ahnafshahriar2003-5611-e28f674b had recent pushes"
3. Click "Compare & Pull Request"
4. Click "Create Pull Request"
5. Click "Merge Pull Request"
6. Confirm the merge

**OR** if you don't see the notification:
1. Go to Pull Requests tab
2. Create new PR from `v0/ahnafshahriar2003-5611-e28f674b` → `main`
3. Merge it

### Step 2: Set Frontend Environment Variables in Vercel

1. Go to **Vercel Dashboard**
2. Select your **go-bank-front-back** project
3. Click **Settings** → **Environment Variables**
4. Add these 2 variables:
   ```
   Name: NEXT_PUBLIC_API_URL
   Value: https://go-bank-front-back.vercel.app/api
   
   Name: NEXT_PUBLIC_BACKEND_URL
   Value: https://go-bank-front-back.vercel.app
   ```
   (Replace `go-bank-front-back.vercel.app` with your actual domain)

5. Click "Add" for each variable
6. Make sure they're added to "Production"

### Step 3: Wait for Vercel to Auto-Deploy

Once you merge to main, Vercel will automatically:
1. Build the frontend
2. Build the backend
3. Deploy both

This takes about 2-3 minutes.

You'll see in Vercel Dashboard:
- Green checkmark = deployment successful
- Red X = deployment failed (check logs)

### Step 4: Test Your Deployment

1. **Test Backend**:
   ```
   curl https://go-bank-front-back.vercel.app/api/health
   ```
   Should return: `{"status":"ok"}`

2. **Test Frontend**:
   - Open in browser: `https://go-bank-front-back.vercel.app`
   - You should see:
     - Loading spinner (few seconds)
     - Then either login page or dashboard redirect

3. **Test Login** (if you see login page):
   - Register a new account or use existing credentials
   - Click login
   - You should get redirected to dashboard

## If Something Goes Wrong

### Frontend shows 404

**Cause**: Frontend not deployed or build failed

**Fix**:
1. Go to Vercel Dashboard
2. Click on your project
3. Check "Deployments" tab
4. Click latest deployment
5. Check if status is ✅ (success) or ❌ (failed)
6. If failed, click to see error logs
7. Common issues:
   - Environment variables not set → go to Step 2
   - Build error → check logs for details

### Frontend connects but shows errors

**Cause**: Backend URL is wrong or backend is down

**Fix**:
1. Open browser DevTools (F12)
2. Go to Console tab
3. Look for error messages
4. Check Network tab to see API calls
5. Verify backend URL is correct in Vercel environment variables

### Backend returns 500 error

**Cause**: Missing environment variables or database connection issue

**Fix**:
1. Go to Vercel → go-bank-front-back-backend project
2. Settings → Environment Variables
3. Check these are all set:
   - `DATABASE_URL`
   - `JWT_SECRET`
   - All other DB and SMTP variables
4. If missing, add them
5. Redeploy

## All Environment Variables Checklist

### Backend (Vercel):
- ✅ DB_HOST
- ✅ DB_PORT
- ✅ DB_USER
- ✅ DB_NAME
- ✅ DB_PASSWORD
- ✅ JWT_SECRET
- ✅ LISTEN_ADDR
- ✅ SMTP_EMAIL
- ✅ SMTP_PASSWORD
- ✅ SMTP_HOST
- ✅ SMTP_PORT
- ✅ COUPON_CODE

### Frontend (Vercel):
- ⬜ NEXT_PUBLIC_API_URL (need to set)
- ⬜ NEXT_PUBLIC_BACKEND_URL (need to set)

## Success Indicators

You're done when you see:

1. ✅ GitHub shows main branch is updated
2. ✅ Vercel shows deployment successful (green checkmark)
3. ✅ Backend health check returns `{"status":"ok"}`
4. ✅ Frontend loads without 404 error
5. ✅ Can login and see dashboard

## Files to Reference

- `FRONTEND_BACKEND_DEPLOYMENT.md` - Detailed deployment guide
- `BACKEND_SETUP_SUMMARY.md` - Backend configuration details
- `BACKEND_ENV_SETUP_CHECKLIST.md` - Environment variable references

## Quick Timeline

- **Right now**: Merge PR on GitHub (2 minutes)
- **Then**: Set frontend env vars in Vercel (2 minutes)
- **Wait**: Vercel auto-deploys (2-3 minutes)
- **Finally**: Test everything (5 minutes)

**Total time: ~10-15 minutes**

## Next Steps After Deployment

Once everything is working:

1. Test all features (register, login, transfer money, etc.)
2. Test on mobile (responsive design)
3. Set up monitoring and error tracking
4. Configure domain (if you have a custom domain)
5. Consider CDN or caching for performance

---

**Questions?** Check the logs in Vercel or review the error messages in browser console.
