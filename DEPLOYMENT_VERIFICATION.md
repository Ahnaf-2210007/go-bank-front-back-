# Deployment Verification Checklist

## Environment Variables Status: ✅ CONFIGURED

All required environment variables have been set in Vercel:

### Database Configuration
- ✅ DB_HOST - PostgreSQL host
- ✅ DB_PORT - Port 5432
- ✅ DB_USER - Database user
- ✅ DB_NAME - Database name
- ✅ DB_PASSWORD - Database password
- ⚠️ DATABASE_URL - Alternative (not needed if individual vars set)

### Authentication & Security
- ✅ JWT_SECRET - Long random secret for token signing

### Server Configuration
- ✅ LISTEN_ADDR - Server listening address

### Email Configuration (Optional)
- ✅ SMTP_EMAIL - Gmail or SMTP email
- ✅ SMTP_PASSWORD - Email password/app password
- ✅ SMTP_HOST - SMTP host (default: smtp.gmail.com)
- ✅ SMTP_PORT - SMTP port (default: 587)

### Business Logic
- ✅ COUPON_CODE - Promotional coupon code

---

## How to Test Deployment

### 1. Test Health Endpoint
After redeployment, test if the backend is running:

```bash
curl https://your-domain.vercel.app/health
```

Expected response:
```json
{"status":"ok"}
```

### 2. Test Login Endpoint (Optional)
```bash
curl -X POST https://your-domain.vercel.app/login \
  -H "Content-Type: application/json" \
  -d '{"number":"1234567890","password":"wrongpassword"}'
```

Expected response:
```json
{"error":"invalid account number or password"}
```

This confirms the API is responding.

### 3. Check Vercel Logs
```bash
vercel logs --follow
```

You should see:
```
Starting GoBank server...
Loading configuration...
Listening on :xxxxx
Connecting to database...
Database connection established
Creating database tables...
Created/verified accounts table
Created/verified coupon_redemptions table
...
JSON API server is running on :xxxxx
```

---

## Frontend Configuration

Update your frontend `.env.local`:

**Development:**
```
NEXT_PUBLIC_API_URL=http://localhost:3000
```

**Production:**
```
NEXT_PUBLIC_API_URL=https://your-backend-domain.vercel.app
```

---

## Testing Frontend-Backend Connection

### 1. Start Frontend
```bash
cd frontend
npm run dev
```

### 2. Go to Login Page
```
http://localhost:8080/login
```

### 3. Test Registration
- Click "Create Account"
- Enter test credentials
- Submit form
- Check browser console for API response

### 4. Check API Connection
In browser DevTools Console:
```javascript
fetch('http://localhost:3000/health')
  .then(r => r.json())
  .then(d => console.log(d))
```

Should log: `{status: "ok"}`

---

## Common Issues & Solutions

### Issue: Health endpoint returns 404 or 503
**Solution**: 
- Verify all environment variables are set in Vercel
- Check DATABASE_URL format if using that instead of individual vars
- Check Vercel build logs for errors

### Issue: Frontend can't connect to backend
**Solution**:
- Verify `NEXT_PUBLIC_API_URL` is set correctly in frontend
- Check CORS is enabled (should be by default)
- Verify backend is actually running (check health endpoint)

### Issue: Database connection fails
**Solution**:
- Verify DB_HOST, DB_USER, DB_PASSWORD are correct
- Ensure database exists and is accessible
- Check if PostgreSQL server is running
- Verify firewall/security group allows connections

### Issue: JWT errors
**Solution**:
- Verify JWT_SECRET is set and is at least 32 characters
- Ensure token is being sent in Authorization header as `Bearer {token}`

### Issue: Emails not sending
**Solution**:
- SMTP is optional - not critical for core functionality
- If needed, verify SMTP credentials
- For Gmail, use "App Passwords" not regular password
- Check spam folder

---

## Redeployment Steps

If you need to redeploy after any changes:

```bash
# Make changes to code
# Commit to git
git add .
git commit -m "your commit message"
git push origin main

# Vercel will auto-deploy on push to main branch
# Or manually trigger:
vercel deploy --prod
```

---

## Next Steps

1. ✅ Environment variables configured
2. ⏳ Redeploy to Vercel (should auto-deploy on push)
3. ✅ Test health endpoint
4. ✅ Test frontend-backend connection
5. ⏳ Proceed to Phase 2: Money Transfer Features

---

## Support

If you encounter any issues:

1. Check DEPLOYMENT_TROUBLESHOOTING.md for detailed solutions
2. Review Vercel logs: `vercel logs --follow`
3. Check browser console for frontend errors
4. Verify all environment variables are set correctly

