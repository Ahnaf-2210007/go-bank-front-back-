# Environment Variables Setup Guide

## For Development

### 1. Create `.env.development.local`

Copy and fill in the following template:

```bash
# Database (required) - from Neon PostgreSQL
DATABASE_URL='postgresql://user:password@host/database?sslmode=require'

# JWT Secret (required) - generate with: openssl rand -base64 32
JWT_SECRET='your-jwt-secret-here'

# WebAuthn Configuration (optional - defaults to localhost)
WEBAUTHN_RP_ORIGIN='http://localhost:8080'
WEBAUTHN_RP_ID='localhost'
WEBAUTHN_DISPLAY_NAME='GoBank'

# SMTP Configuration (optional - for email verification)
SMTP_EMAIL='your-email@gmail.com'
SMTP_PASSWORD='your-app-password'
SMTP_HOST='smtp.gmail.com'
SMTP_PORT='587'

# Server Configuration (optional)
LISTEN_ADDR=':3000'
COUPON_CODE='OFFER1000'
```

### 2. Get Database Connection String from Neon

1. Go to https://console.neon.tech
2. Select your project
3. Copy the connection string from "Connection string" section
4. Add to `.env.development.local` as `DATABASE_URL`

Example:
```
DATABASE_URL='postgresql://neondb_owner:password@ep-example.c.us-east-1.aws.neon.tech/neondb?sslmode=require'
```

### 3. Generate JWT Secret

```bash
# On macOS/Linux
openssl rand -base64 32

# On Windows (PowerShell)
[Convert]::ToBase64String([System.Security.Cryptography.RandomNumberGenerator]::GetBytes(32))

# Or use any random string of 32+ characters
```

Copy the output and add to `.env.development.local`:
```
JWT_SECRET='your-generated-secret-here'
```

### 4. Setup Gmail SMTP (Optional - for email sending)

For sending verification emails:

1. Enable 2-Step Verification on your Google Account
2. Create an App Password:
   - Go to myaccount.google.com
   - Click "Security" in left menu
   - Scroll to "App passwords"
   - Select Mail and Windows Computer
   - Google generates a 16-character password
   
3. Add to `.env.development.local`:
```
SMTP_EMAIL='your-email@gmail.com'
SMTP_PASSWORD='your-16-char-app-password'
```

If not set, the backend will log verification codes to console instead.

### 5. WebAuthn Configuration

For local development, defaults are fine:
```
WEBAUTHN_RP_ORIGIN='http://localhost:8080'
WEBAUTHN_RP_ID='localhost'
WEBAUTHN_DISPLAY_NAME='GoBank'
```

Only change if:
- Frontend runs on different port: update WEBAUTHN_RP_ORIGIN
- Running on different domain: update WEBAUTHN_RP_ID

### 6. Run Backend

```bash
cd /vercel/share/v0-project
go run .
```

Should print:
```
JSON API server is running on :3000
```

Test it:
```bash
curl http://localhost:3000/health
# Should return: {"status":"ok"}
```

---

## For Production (Vercel)

### 1. Set Environment Variables in Vercel

Go to Vercel Dashboard → Project Settings → Environment Variables

Add the following:

#### Required

| Variable | Value | Example |
|----------|-------|---------|
| `DATABASE_URL` | PostgreSQL connection from Neon | `postgresql://user:pass@host/db?sslmode=require` |
| `JWT_SECRET` | Strong random string (32+ chars) | `abc123...xyz789` |

#### Optional (WebAuthn)

| Variable | Value | Example |
|----------|-------|---------|
| `WEBAUTHN_RP_ORIGIN` | Your app's full URL | `https://myapp.vercel.app` |
| `WEBAUTHN_RP_ID` | Domain without protocol | `myapp.vercel.app` |
| `WEBAUTHN_DISPLAY_NAME` | Display name for authenticators | `GoBank` |

#### Optional (Email)

| Variable | Value | Example |
|----------|-------|---------|
| `SMTP_EMAIL` | Gmail address | `myapp@gmail.com` |
| `SMTP_PASSWORD` | Gmail app password | `abcd efgh ijkl mnop` |
| `SMTP_HOST` | SMTP server | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP port | `587` |

#### Optional (Server)

| Variable | Value | Example |
|----------|-------|---------|
| `LISTEN_ADDR` | Server port | `:3000` |
| `COUPON_CODE` | Valid coupon code | `OFFER1000` |

### 2. Set Environment for Different Stages

Optionally set different values for Preview/Production:

**Production** (main branch):
- Use production database connection
- Use production domain for WebAuthn

**Preview** (pull requests):
- Can use same database or separate preview database
- Use preview domain for WebAuthn (if available)

### 3. Deploy

```bash
git add .
git commit -m "Backend setup complete with CORS and WebAuthn config"
git push
```

Vercel automatically deploys and sets environment variables.

### 4. Verify Deployment

```bash
curl https://your-domain.vercel.app/health
# Should return: {"status":"ok"}
```

---

## Environment Variable Reference

### DATABASE_URL (Required)

PostgreSQL connection string from Neon.

**Format:**
```
postgresql://user:password@host:port/database?sslmode=require
```

**Where to get it:**
1. Neon Console → Project → Connection string
2. Copy the full string including credentials

### JWT_SECRET (Required)

Secret key for signing JWT authentication tokens.

**Requirements:**
- Minimum 32 characters
- Should be random and unique
- Never commit to repository
- Different value for each environment

**Generate:**
```bash
openssl rand -base64 32
```

### WEBAUTHN_RP_ORIGIN (Optional)

The origin (protocol + domain) where WebAuthn credentials are registered.

**Default:** `http://localhost:8080`

**For production:**
- Must match your actual frontend domain
- Must include protocol (http/https)
- Example: `https://myapp.vercel.app`

**Important:** If mismatch, WebAuthn will fail silently

### WEBAUTHN_RP_ID (Optional)

The relying party ID for WebAuthn (domain without protocol).

**Default:** `localhost`

**For production:**
- Must be the domain portion only (no http/https)
- Must match domain in WEBAUTHN_RP_ORIGIN
- Example: `myapp.vercel.app`

### WEBAUTHN_DISPLAY_NAME (Optional)

Display name shown in authenticator apps.

**Default:** `GoBank`

**Example:** `My Banking App`

### SMTP_EMAIL (Optional)

Email address to send from (Gmail).

**Format:** `your-email@gmail.com`

**Note:** Use app-specific password, not account password

### SMTP_PASSWORD (Optional)

App-specific password from Gmail.

**Note:** 
- 16 characters with spaces
- Remove spaces when adding to env var
- Required for email verification to work

### SMTP_HOST (Optional)

SMTP server host.

**Default:** `smtp.gmail.com`

**For Gmail:** Always `smtp.gmail.com`

### SMTP_PORT (Optional)

SMTP server port.

**Default:** `587`

**For Gmail:** Always `587`

### LISTEN_ADDR (Optional)

Server listen address and port.

**Default:** `:3000`

**Format:** `:port` or `host:port`

**Examples:**
- `:3000` - Listen on all interfaces, port 3000
- `localhost:8000` - Listen only on localhost, port 8000

### COUPON_CODE (Optional)

Valid coupon code for account offers.

**Default:** `OFFER1000`

**Used by:** POST /account/{id}/offer endpoint

---

## Checking Variables

### Development

List loaded variables:
```bash
# Show all env vars
env | grep -E "DATABASE_URL|JWT_SECRET|WEBAUTHN|SMTP|LISTEN_ADDR|COUPON_CODE"

# Check if required vars are set
if [ -z "$JWT_SECRET" ]; then echo "JWT_SECRET not set"; fi
if [ -z "$DATABASE_URL" ]; then echo "DATABASE_URL not set"; fi
```

### Production

View in Vercel Dashboard:
1. Project Settings → Environment Variables
2. Can see which variables are set for each environment

Or via CLI:
```bash
vercel env ls

# View specific variable (shows if set)
vercel env pull
```

---

## Troubleshooting

### Missing DATABASE_URL

**Error:** `DB_PASSWORD (or DATABASE_URL) environment variable must be set`

**Solution:**
1. Get connection string from Neon
2. Add as `DATABASE_URL` to environment
3. Verify connection string is valid

### Missing JWT_SECRET

**Error:** `JWT_SECRET environment variable must be set`

**Solution:**
1. Generate: `openssl rand -base64 32`
2. Add to `.env.development.local` or Vercel project
3. Restart backend

### WebAuthn Not Working

**Common causes:**
1. WEBAUTHN_RP_ORIGIN doesn't match frontend URL
2. WEBAUTHN_RP_ID doesn't match domain

**Solution:**
1. Frontend at `http://localhost:3000`?
   - Set `WEBAUTHN_RP_ORIGIN='http://localhost:3000'`
   
2. Frontend at `https://myapp.vercel.app`?
   - Set `WEBAUTHN_RP_ORIGIN='https://myapp.vercel.app'`
   - Set `WEBAUTHN_RP_ID='myapp.vercel.app'`

### Email Not Sending

**If emails not sending:**
1. Check `SMTP_EMAIL` and `SMTP_PASSWORD` are set
2. Backend logs verification code instead (check logs)
3. For Gmail:
   - Use 16-character app password (not account password)
   - Enable 2-Step Verification first
   - Check "Allow less secure apps" if using full password

### Connection Refused

**Error:** `connect: connection refused`

**Solution:**
1. Verify DATABASE_URL is correct
2. Check Neon database is online
3. Verify network access is allowed
4. Test connection: `psql "postgresql://..."` (if psql installed)

---

## Security Best Practices

1. **Never commit .env files**
   - Add `.env.development.local` to `.gitignore`
   - Never push production secrets

2. **Use strong JWT_SECRET**
   - Minimum 32 characters
   - Generate with cryptographic randomness
   - Different per environment

3. **Rotate secrets regularly**
   - Change JWT_SECRET periodically
   - Use different passwords for staging/production

4. **Restrict database access**
   - Use different database user per environment
   - Set appropriate database permissions
   - Use IP whitelisting if available

5. **Monitor access**
   - Review environment variables access logs
   - Check production deployments for secrets
   - Use tools like git-secrets to prevent commits

---

## Quick Setup Checklist

- [ ] Get DATABASE_URL from Neon
- [ ] Generate JWT_SECRET with openssl
- [ ] Create `.env.development.local`
- [ ] Add DATABASE_URL and JWT_SECRET
- [ ] Add WebAuthn vars (optional, use defaults)
- [ ] Add SMTP vars (optional, for emails)
- [ ] Test with `go run .`
- [ ] Test with `curl http://localhost:3000/health`
- [ ] Add to Vercel project settings for production
- [ ] Deploy and verify with `curl https://your-domain/health`

---

## Questions?

If environment variables aren't loading:
1. Check file path is correct
2. Verify variable names are exact (case-sensitive)
3. Check for typos in variable names
4. Reload/restart the application
5. Check logs for error messages

For production issues:
1. Check Vercel build logs
2. Verify all required variables are set
3. Check for typos in variable names
4. Check values are correct for production domain
5. Review deployment logs: `vercel logs --tail`
