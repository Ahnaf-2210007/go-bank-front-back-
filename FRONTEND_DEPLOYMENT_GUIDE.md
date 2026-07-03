# Frontend Deployment Guide

## Overview
The frontend (Next.js) and backend (Go) are deployed separately on Vercel:
- **Backend**: Already deployed at `https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app`
- **Frontend**: Needs separate deployment

## Step-by-Step Frontend Deployment

### Option 1: Deploy from GitHub (Recommended)

1. **Create a new Vercel project for frontend**
   - Go to vercel.com/new
   - Import your GitHub repository
   - Select "Next.js" as the framework
   - Root directory: `frontend`

2. **Set environment variables in Vercel**
   - Go to Settings → Environment Variables
   - Add: `NEXT_PUBLIC_API_URL=https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app`
   - This URL is your backend domain

3. **Deploy**
   - Click "Deploy"
   - Vercel will build and deploy automatically

### Option 2: Deploy using Vercel CLI

```bash
cd frontend
npm i -g vercel
vercel deploy --prod --env NEXT_PUBLIC_API_URL=https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app
```

## Configuration Files

### `.env.local` (Development)
```
NEXT_PUBLIC_API_URL=http://localhost:3000
```
Use this when running frontend and backend locally.

### `.env.production` (Production)
```
NEXT_PUBLIC_API_URL=https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app
```
Use this for production deployment.

## Testing the Connection

### 1. Local Testing
```bash
# Terminal 1: Backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```
- Backend: http://localhost:3000
- Frontend: http://localhost:8080
- Try registering an account to test connection

### 2. Production Testing
1. Go to your frontend URL
2. Try the login/register flow
3. Check browser DevTools (F12) → Network tab
4. Verify API calls go to backend domain
5. Check Console tab for any errors

## Troubleshooting

### Frontend shows 404
- Check that frontend is deployed (not the backend root)
- Verify environment variables are set

### API calls fail (CORS error)
- Confirm backend has CORS enabled (it does by default)
- Check `NEXT_PUBLIC_API_URL` matches backend domain
- Verify backend is accessible at health endpoint

### Frontend can't reach backend
- Test backend health: `curl https://your-backend-url/health`
- Check environment variable in Vercel settings
- Verify both apps are deployed and running

## Architecture

```
┌─────────────────────────────────┐
│   Frontend (Next.js)            │
│   - vercel.app domain 1         │
│   - Serves UI                   │
│   - Makes API calls to backend  │
└────────────────┬────────────────┘
                 │
                 │ HTTPS API calls
                 ↓
┌─────────────────────────────────┐
│   Backend (Go)                  │
│   - vercel.app domain 2         │
│   - Handles business logic      │
│   - Database operations         │
│   - Returns JSON responses      │
└─────────────────────────────────┘
```

## Deployment Checklist

- [ ] Backend deployed and accessible
- [ ] Frontend created as separate Vercel project
- [ ] Environment variables set in frontend project
- [ ] Frontend deployed successfully
- [ ] Health check works: `curl https://backend-url/health`
- [ ] Frontend loads without 404 errors
- [ ] API calls work from frontend
- [ ] Authentication flows work end-to-end

## Next Steps

Once both are deployed and working:
1. Proceed with Phase 2 development (Money Transfers, Transaction History)
2. Add more features to dashboard
3. Optimize performance
4. Set up monitoring and logging

## Support

For issues:
1. Check `DEPLOYMENT_TROUBLESHOOTING.md` for backend issues
2. Check frontend logs in Vercel dashboard
3. Verify environment variables are correct
4. Test individual endpoints with curl
