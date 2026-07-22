# Vercel Frontend Deployment - Step by Step Guide

## Current Status
You are on the Vercel "New Project" page with the repository already selected:
- Repository: `Ahnaf-2210007/online-banking-system-ds`
- Branch: `main`
- Team: Ahnaf Shahriar's projects

## Step 1: Click Create Button
Click the large **"Create"** button at the bottom of the screen.

This will import your repository and take you to the configuration screen.

## Step 2: Configure Root Directory
On the next screen (Project Settings), you will see:

**Look for:** "Root Directory" or "Project Settings"
**Action:** Click the Root Directory dropdown/field and select or type: `frontend`

## Step 3: Add Environment Variables
Before clicking Deploy, add the environment variable:

**Variable 1:**
- Name: `NEXT_PUBLIC_API_URL`
- Value: `https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app`

(Replace with your actual backend Vercel URL if different)

**Why this variable?**
- It tells the frontend where to find the backend API
- `NEXT_PUBLIC_` prefix makes it available in the browser
- Frontend will use this URL for all API calls

## Step 4: Deploy
Click the **"Deploy"** button and wait 2-5 minutes for the build to complete.

## Step 5: Access Frontend
Once deployment is complete:
1. You'll get a new frontend URL like: `https://your-project-name.vercel.app`
2. Visit that URL in your browser
3. You should see the dark-themed login page

## What to Expect After Deployment

### Success Indicators
✓ Page loads without 404 errors
✓ Dark theme is visible (black background, blue accents)
✓ Login page with email and password fields
✓ Links to Register and Verify pages work
✓ Console shows no CORS errors

### Testing the Connection
1. Go to Register page
2. Enter test credentials:
   - Name: Test User
   - Email: test@example.com
   - Password: TestPassword123
3. Submit form
4. Check browser console (F12) for any errors
5. If successful, you should see a verification code request or confirmation

### Common Issues During Deployment

| Issue | Solution |
|-------|----------|
| Build fails with "module not found" | Ensure node_modules folder is in .gitignore |
| Page loads but shows 404 | Check that root directory is set to "frontend" |
| API calls fail with CORS error | Verify NEXT_PUBLIC_API_URL is correct |
| Dark theme not showing | Clear browser cache and hard refresh (Ctrl+Shift+R) |

## Reference URLs

Your deployment will create URLs:
- Frontend: `https://[your-frontend-project].vercel.app`
- Backend: `https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app`
- Backend Health Check: `https://online-banking-system-7ssih09kt-ahnaf-shahriars-projects.vercel.app/health`

## Next Steps After Frontend Deployment

1. Test the authentication flow
2. Verify API communication works
3. Proceed to Phase 2 development (Money Transfer, Transaction History, etc.)

## Need Help?

If deployment fails:
1. Check Vercel deployment logs (Deployments tab)
2. Look for specific error messages
3. Verify frontend/.env.production has correct backend URL
4. Ensure all dependencies in package.json are correct
