# GoBank Deployment Troubleshooting Guide

## Common Issues and Solutions

### Issue 1: Server Not Opening / 404 Errors

**Symptoms:**
- Deployment shows "Build Success" but accessing the URL returns 404
- Server appears to start but doesn't respond

**Solutions:**

1. **Check Environment Variables in Vercel Dashboard:**
   - Go to Settings → Environment Variables
   - Ensure these variables are set:
     - `DATABASE_URL` (or `DB_PASSWORD`, `DB_HOST`, etc.)
     - `JWT_SECRET`
   - Variables should NOT have quotes around them

2. **Verify Port Configuration:**
   - Vercel automatically assigns a PORT environment variable
   - The backend now uses this PORT (should work automatically)

3. **Check API Health Endpoint:**
   ```bash
   curl https://your-domain.vercel.app/health
   ```
   Should return: `{"status":"ok"}`

---

### Issue 2: Database Connection Error

**Symptoms:**
- Logs show "Database connection error"
- Cannot connect to PostgreSQL

**Causes:**
- Missing DATABASE_URL environment variable
- DATABASE_URL format is incorrect
- Database credentials are wrong
- Network/firewall blocking access

**Solutions:**

1. **Verify DATABASE_URL Format:**
   ```
   postgresql://user:password@host:port/database
   # or
   postgres://user:password@host:port/database
   ```

2. **For Neon Database:**
   ```
   Get from: https://console.neon.tech
   Format: postgresql://[user]:[password]@[host].neon.tech/[database]?sslmode=require
   ```

3. **Test Connection Locally:**
   ```bash
   export DATABASE_URL="postgresql://user:password@host/db"
   go run main.go
   ```

4. **Check Neon Connection:**
   - Ensure connection pooling is enabled (if using PgBouncer)
   - Verify IP allowlist doesn't block Vercel (usually shouldn't matter)

---

### Issue 3: JWT_SECRET Not Set

**Symptoms:**
- Error: "JWT_SECRET environment variable must be set"
- Server won't start in production

**Solutions:**

1. **Generate JWT_SECRET:**
   ```bash
   # Option 1: Using openssl
   openssl rand -base64 32
   
   # Option 2: Any random string (min 32 chars)
   ```

2. **Add to Vercel:**
   - Settings → Environment Variables
   - Key: `JWT_SECRET`
   - Value: Your generated secret

3. **Verify Local Setup:**
   - `.env.development.local` should have `JWT_SECRET` set
   - Or export as environment variable

---

### Issue 4: CORS Errors from Frontend

**Symptoms:**
- Frontend shows CORS error in console
- Cross-origin requests blocked

**Solutions:**

CORS is already configured in the backend:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, DELETE, OPTIONS, PUT, PATCH
```

If still getting errors:

1. **Check Backend Logs:**
   - Look for actual error (not just CORS)
   - Log shows detailed error messages

2. **Verify Frontend URL:**
   - Update `.env.local` in frontend with correct backend URL
   - `NEXT_PUBLIC_API_URL=https://your-backend.vercel.app`

3. **Test Directly:**
   ```bash
   curl -X OPTIONS https://your-domain.vercel.app/login \
     -H "Access-Control-Request-Method: POST"
   ```

---

### Issue 5: Build Fails on Vercel

**Symptoms:**
- Build Failed status on Vercel
- Error in build log

**Common Build Errors:**

1. **Go Version Issue:**
   - Vercel supports Go 1.18+
   - Check `go.mod` file version

2. **Missing Dependencies:**
   ```bash
   # Ensure go.mod has all dependencies
   go mod download
   go mod verify
   ```

3. **vercel.json Configuration:**
   - Should NOT have `functions` or `requiredEnvs`
   - Correct format:
     ```json
     {
       "buildCommand": "go build -o api .",
       "devCommand": "go run .",
       "framework": "go",
       "installCommand": "go mod download"
     }
     ```

---

### Issue 6: Slow Response / Timeouts

**Symptoms:**
- Requests take >30 seconds
- 504 Gateway Timeout
- Server-side processing slow

**Solutions:**

1. **Check Database Query Performance:**
   - Look for N+1 queries
   - Add indexes if needed

2. **Monitor Backend Logs:**
   - See what's taking long
   - Database queries? API calls?

3. **Increase Timeout (if needed):**
   - Frontend axios instance
   - Adjust timeout in `frontend/lib/api.ts`

---

### Issue 7: Email Verification Not Working

**Symptoms:**
- Registration succeeds but verification email not sent
- Verification endpoint returns error

**Causes:**
- SMTP credentials not configured
- Email service blocked
- Invalid email format

**Solutions:**

1. **Verify SMTP Configuration:**
   ```
   SMTP_EMAIL: your-email@gmail.com
   SMTP_PASSWORD: app-password (NOT regular password)
   SMTP_HOST: smtp.gmail.com
   SMTP_PORT: 587
   ```

2. **For Gmail:**
   - Use App Password (not regular password)
   - Enable 2-factor authentication
   - https://myaccount.google.com/apppasswords

3. **Test Email Sending:**
   - Check server logs for email service errors
   - Verify email address format

---

## Debugging Steps

### 1. Check Vercel Logs

**Command Line:**
```bash
vercel logs --follow
```

**Dashboard:**
- Go to Vercel Project
- Click "Deployments"
- Select latest deployment
- Click "View Build Logs"

### 2. Enable Debug Logging

In `main.go`, logging is already enhanced:
```go
log.Println("Starting GoBank server...")
log.Printf("Listening on %s", cfg.ListenAddr)
log.Printf("Created/verified %s table", table.name)
```

### 3. Check Environment Variables

**Vercel Dashboard:**
- Project Settings
- Environment Variables
- Verify all required variables exist

**Terminal:**
```bash
vercel env list
```

### 4. Test API Endpoints

```bash
# Health check
curl https://your-domain.vercel.app/health

# Login (test)
curl -X POST https://your-domain.vercel.app/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'
```

---

## Production Checklist

- [ ] DATABASE_URL set in Vercel
- [ ] JWT_SECRET set in Vercel
- [ ] SMTP credentials set (if using email)
- [ ] WebAuthn settings configured
- [ ] CORS headers verified
- [ ] Frontend API URL pointing to production
- [ ] SSL/HTTPS enabled
- [ ] Domain DNS configured
- [ ] Error logging enabled
- [ ] Database backups configured

---

## Getting Help

### Check These Files First

1. **Backend Setup**: `BACKEND_SETUP.md`
2. **Frontend Integration**: `FRONTEND_INTEGRATION.md`
3. **Environment Variables**: `ENV_SETUP.md`
4. **Quick Start**: `QUICK_START.md`

### Vercel Support

- https://vercel.com/help
- https://vercel.com/docs/deployments/troubleshooting

### Go Backend Resources

- https://go.dev/doc
- https://github.com/gorilla/mux (Router)
- https://golang-jwt.github.io/jwt (JWT)

---

## Quick Fix Commands

```bash
# Rebuild and redeploy
vercel deploy --prod

# Check deployment status
vercel status

# View recent deployments
vercel deployments list

# Inspect project
vercel inspect

# Set environment variable
vercel env add DATABASE_URL
```
