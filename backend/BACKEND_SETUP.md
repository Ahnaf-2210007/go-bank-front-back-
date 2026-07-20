# Online Banking System - Backend Setup Guide

## Overview

This is a Go backend for an online banking system using PostgreSQL (Neon) database. The backend provides RESTful APIs for user authentication, account management, WebAuthn support, and fund transfers.

## Architecture

- **Language**: Go
- **Database**: PostgreSQL (via Neon)
- **Authentication**: JWT tokens and WebAuthn
- **Email**: SMTP (Gmail)
- **Deployment**: Vercel

## Environment Variables

All configuration is done through environment variables. Here's what you need:

### Required Environment Variables

```
DATABASE_URL                  # PostgreSQL connection string (from Neon)
JWT_SECRET                    # Secret key for signing JWT tokens (should be 32+ chars)
```

### Optional Environment Variables (with defaults)

```
LISTEN_ADDR                   # Server listen address (default: :3000)
COUPON_CODE                   # Coupon code for account offers (default: OFFER1000)
WEBAUTHN_RP_ORIGIN           # WebAuthn origin URL (default: http://localhost:8080)
WEBAUTHN_RP_ID               # WebAuthn relying party ID (default: localhost)
WEBAUTHN_DISPLAY_NAME        # WebAuthn display name (default: GoBank)
SMTP_EMAIL                   # SMTP sender email address
SMTP_PASSWORD                # SMTP password
SMTP_HOST                    # SMTP host (default: smtp.gmail.com)
SMTP_PORT                    # SMTP port (default: 587)
```

## API Endpoints

### Health Check
- `GET /health` - Returns server status

### Authentication & Account Management
- `POST /login` - Login with account number and password
- `POST /account` - Create new account
- `POST /account/verification` - Verify email with OTP
- `GET /account/{id}` - Get account details (requires JWT)
- `DELETE /account/{id}` - Delete account (requires JWT)
- `POST /account/update` - Update account profile/email/password (requires JWT)
- `GET /account/transactions` - Get transaction history (requires JWT)
- `POST /account/{id}/offer` - Apply coupon offer (requires JWT)

### Transfers
- `POST /transfer` - Transfer funds between accounts (requires JWT)

### WebAuthn (Passwordless Authentication)
- `POST /webauthn/register/begin` - Start WebAuthn registration
- `POST /webauthn/register/finish` - Complete WebAuthn registration
- `POST /webauthn/login/begin` - Start WebAuthn login
- `POST /webauthn/login/finish/{email}` - Complete WebAuthn login

## CORS Configuration

The backend automatically handles CORS requests with the following headers:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, DELETE, OPTIONS, PUT, PATCH`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With`
- `Access-Control-Max-Age: 86400`

Preflight (OPTIONS) requests are automatically handled.

## Database Schema

The backend automatically creates the following tables on startup:

1. **accounts** - User account data
2. **pending_accounts** - Accounts pending email verification
3. **transactions** - Fund transfer records
4. **coupon_redemptions** - Coupon usage tracking
5. **pending_profile_updates** - Email change verification
6. **webauthn_credentials** - WebAuthn authenticator data

## Running Locally

1. Set up environment variables in `.env.development.local`
2. Run the application: `go run .`
3. Server listens on `localhost:3000` by default

## Deployment to Vercel

1. Make sure all required environment variables are set in Vercel project settings
2. Vercel will automatically detect `vercel.json` configuration
3. Deployment: `git push` to trigger Vercel deployment

The `vercel.json` file configures:
- Build command: `go build -o api .`
- Development command: `go run .`
- Function timeout: 60 seconds
- Memory: 1024 MB

## Frontend Integration

### Base URL
- Development: `http://localhost:3000`
- Production: `https://your-vercel-domain.vercel.app`

### Request Format
All requests use JSON:
```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer YOUR_JWT_TOKEN"  // For protected endpoints
}
```

### Response Format
Successful responses:
```json
{
  "status": "ok",
  // response data
}
```

Error responses:
```json
{
  "error": "error message"
}
```

### Authentication
Protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <token>
```

## Fixed Issues

### Issue 1: Missing CORS Headers
**Problem**: Serverless function was crashing when called from frontend due to missing CORS headers.

**Solution**: Added CORS middleware that:
- Adds proper CORS headers to all responses
- Handles preflight OPTIONS requests
- Allows cross-origin requests from any domain

### Issue 2: Hardcoded WebAuthn Config
**Problem**: WebAuthn was hardcoded to `localhost`, failing in production.

**Solution**: Moved configuration to environment variables:
- `WEBAUTHN_RP_ORIGIN` - Full origin URL
- `WEBAUTHN_RP_ID` - Domain for relying party
- `WEBAUTHN_DISPLAY_NAME` - Display name for authenticators

### Issue 3: Missing Health Check
**Problem**: No way to verify if backend is running.

**Solution**: Added `/health` endpoint that returns `{"status": "ok"}` for monitoring.

### Issue 4: Missing Production Configuration
**Problem**: No deployment configuration for Vercel.

**Solution**: Created `vercel.json` with:
- Proper build and dev commands for Go
- Function memory and timeout settings
- Environment variable declarations

## Security Considerations

1. **JWT_SECRET**: Must be a strong random string (32+ characters)
2. **CORS**: Currently allows all origins. Consider restricting to your frontend domain in production:
   - Modify `corsMiddleware` in api.go to check `r.Header.Get("Origin")`
3. **Database**: Use environment variables from Neon, never commit credentials
4. **SMTP**: Use app-specific passwords for Gmail, not account passwords
5. **WebAuthn**: Ensure `WEBAUTHN_RP_ORIGIN` matches exactly (protocol + domain)

## Troubleshooting

### Server won't start
- Check all required environment variables are set
- Verify DATABASE_URL is correct and accessible
- Check JWT_SECRET is set

### CORS errors from frontend
- Verify backend is running
- Check `/health` endpoint responds
- Ensure CORS middleware is applied to router

### WebAuthn not working
- Verify WEBAUTHN_RP_ID matches your domain (no protocol)
- Verify WEBAUTHN_RP_ORIGIN includes full URL with protocol
- Check browser console for error details

### Email not sending
- Verify SMTP_EMAIL and SMTP_PASSWORD are set
- Check Gmail app-specific password is correct
- Verify SMTP_HOST and SMTP_PORT are correct

## Next Steps

1. **Frontend Integration**: Connect your frontend to these endpoints
2. **Testing**: Test all endpoints with Postman or similar
3. **Monitoring**: Set up error tracking (e.g., Sentry)
4. **Rate Limiting**: Consider adding rate limiting for production
5. **Validation**: Add more comprehensive input validation
6. **Logging**: Enhance server-side logging for debugging

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review environment variable configuration
3. Check backend logs during deployment
4. Review frontend console for error messages
