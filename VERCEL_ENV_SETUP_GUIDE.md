# Vercel Environment Variables Setup Guide

## Problem
The "Registration error - Load failed" is happening because environment variables are not set in your Vercel project dashboard.

## Solution: Set Environment Variables in Vercel Dashboard

### Step 1: Access Vercel Project Settings
1. Go to https://vercel.com/dashboard
2. Click on your project: `go-bank-front-back-`
3. Go to **Settings** → **Environment Variables**

### Step 2: Add These Environment Variables

Add the following variables (these are for production):

| Variable Name | Value | Type |
|---|---|---|
| `NEXT_PUBLIC_API_URL` | `https://go-bank-front-back-ivory.vercel.app` | Public |
| `NEXT_PUBLIC_BACKEND_URL` | `https://go-bank-front-back-ivory.vercel.app` | Public |
| `AI_GATEWAY_API_KEY` | `vck_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` (your actual key) | Secret |

**Important:** Make sure to select the correct environments:
- `NEXT_PUBLIC_*` variables should be available in: **Production**, **Preview**, **Development**
- `AI_GATEWAY_API_KEY` should be: **Secret** (available in all environments)

### Step 3: Verify Configuration

After setting the variables:
1. The variables should appear in your dashboard
2. Redeploy your frontend (push to main branch or click "Redeploy")
3. Test the registration page again

### Why This Happened
- `vercel.json` is for **declaring** which variables are needed
- The **actual values** must be set in the Vercel dashboard
- Without these values set in Vercel, your app can't connect to the backend

## Local Development
For local development, use `.env.development.local`:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
AI_GATEWAY_API_KEY=your-local-key-here
```

## Backend Domain Reference
- Backend URL: `https://go-bank-front-back-ivory.vercel.app`
- This is where your Go backend is deployed

## Next Steps
1. Set the environment variables in Vercel dashboard
2. Commit and push the latest code
3. Vercel will automatically redeploy
4. Test registration again
