# Frontend Integration Guide

## Quick Start

### Backend API Base URL

**Development:**
```
http://localhost:3000
```

**Production:**
```
https://online-banking-system-go.vercel.app  // Replace with actual domain
```

### CORS is Enabled

The backend already handles CORS, so you can call it directly from your frontend without proxy issues.

## Example API Calls

### 1. Create Account

```javascript
const response = await fetch('http://localhost:3000/account', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    firstName: 'John',
    lastName: 'Doe',
    email: 'john@example.com',
    password: 'SecurePassword123'
  })
});

const data = await response.json();
// Returns: { message: "verification code sent to email" }
```

### 2. Verify Email

```javascript
const response = await fetch('http://localhost:3000/account/verification', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    code: '123456'  // 6-digit code from email
  })
});

const account = await response.json();
// Returns account object with created account details
```

### 3. Login

```javascript
const response = await fetch('http://localhost:3000/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    number: 1234567890,  // Account number
    password: 'SecurePassword123'
  })
});

const { token, number } = await response.json();
// Store token in localStorage or sessionStorage
localStorage.setItem('authToken', token);
```

### 4. Get Account Details (Authenticated)

```javascript
const token = localStorage.getItem('authToken');

const response = await fetch('http://localhost:3000/account/123', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const account = await response.json();
```

### 5. Transfer Funds

```javascript
const token = localStorage.getItem('authToken');

const response = await fetch('http://localhost:3000/transfer', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    toAccount: 9876543210,  // Recipient account number
    amount: 1000
  })
});

const result = await response.json();
```

### 6. Get Transaction History

```javascript
const token = localStorage.getItem('authToken');

const response = await fetch('http://localhost:3000/account/transactions?limit=20&offset=0', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const transactions = await response.json();
```

### 7. Update Account Profile

```javascript
const token = localStorage.getItem('authToken');

const response = await fetch('http://localhost:3000/account/update', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    action: 'profile',
    firstName: 'Jane',
    lastName: 'Smith'
  })
});

const updated = await response.json();
```

### 8. Health Check

```javascript
const response = await fetch('http://localhost:3000/health');
const { status } = await response.json();
// Returns: { status: "ok" }
```

## WebAuthn Integration

WebAuthn provides passwordless authentication. Here's how to use it:

### Register with WebAuthn

```javascript
// Step 1: Begin registration
const response = await fetch('http://localhost:3000/webauthn/register/begin', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const options = await response.json();

// Step 2: Use WebAuthn API (browser)
const credential = await navigator.credentials.create(options);

// Step 3: Finish registration
const finishResponse = await fetch('http://localhost:3000/webauthn/register/finish', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify(credential)
});

const result = await finishResponse.json();
```

### Login with WebAuthn

```javascript
// Step 1: Begin login
const response = await fetch('http://localhost:3000/webauthn/login/begin', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'john@example.com'
  })
});

const options = await response.json();

// Step 2: Use WebAuthn API (browser)
const assertion = await navigator.credentials.get(options);

// Step 3: Finish login
const finishResponse = await fetch(
  'http://localhost:3000/webauthn/login/finish/john@example.com',
  {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(assertion)
  }
);

const { token, number } = await finishResponse.json();
localStorage.setItem('authToken', token);
```

## Error Handling

All errors return JSON with an `error` field:

```javascript
const response = await fetch('http://localhost:3000/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ number: 123, password: 'wrong' })
});

if (!response.ok) {
  const { error } = await response.json();
  console.error('Error:', error);
  // Returns: { error: "not authorized" }
}
```

## Common HTTP Status Codes

- **200 OK** - Request successful
- **400 Bad Request** - Invalid request data
- **401 Unauthorized** - Missing or invalid JWT token
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error

## Authentication Pattern

Use this pattern for all authenticated requests:

```javascript
class APIClient {
  constructor(baseURL = 'http://localhost:3000') {
    this.baseURL = baseURL;
  }

  async getAuthToken() {
    return localStorage.getItem('authToken');
  }

  async request(endpoint, options = {}) {
    const token = await this.getAuthToken();
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers
    };

    if (token && options.authenticated !== false) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers
    });

    if (!response.ok) {
      const { error } = await response.json();
      throw new Error(error || 'Request failed');
    }

    return response.json();
  }

  // Convenience methods
  async get(endpoint) {
    return this.request(endpoint, { method: 'GET' });
  }

  async post(endpoint, body) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(body)
    });
  }

  async delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }
}

// Usage
const api = new APIClient();

// Login
const { token } = await api.post('/login', {
  number: 1234567890,
  password: 'password'
});
localStorage.setItem('authToken', token);

// Get account
const account = await api.get('/account/123');

// Transfer
await api.post('/transfer', {
  toAccount: 9876543210,
  amount: 1000
});
```

## Environment Variables for Frontend

Add these to your frontend `.env` file:

```
VITE_API_BASE_URL=http://localhost:3000          # Development
VITE_API_BASE_URL=https://your-domain.vercel.app # Production
```

Then use:
```javascript
const baseURL = import.meta.env.VITE_API_BASE_URL;
const api = new APIClient(baseURL);
```

## Testing with curl

Test the backend without frontend:

```bash
# Health check
curl http://localhost:3000/health

# Create account
curl -X POST http://localhost:3000/account \
  -H "Content-Type: application/json" \
  -d '{"firstName":"John","lastName":"Doe","email":"john@example.com","password":"SecurePassword123"}'

# Login
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"number":1234567890,"password":"SecurePassword123"}'

# Get account (replace TOKEN with actual token)
curl http://localhost:3000/account/123 \
  -H "Authorization: Bearer TOKEN"
```

## Common Integration Issues

### CORS Error
**Solution**: Backend already has CORS enabled. If error persists:
1. Check backend is running (`/health` should work)
2. Verify the URL is exactly right (no trailing slash issues)
3. Check browser console for actual error

### JWT Token Expired
**Solution**: Implement token refresh or re-login flow:
```javascript
if (error === 'not authorized') {
  // Redirect to login
  localStorage.removeItem('authToken');
  window.location.href = '/login';
}
```

### WebAuthn Not Supported
**Solution**: Add fallback to password authentication:
```javascript
if (!window.PublicKeyCredential) {
  // Show password login form instead
  showPasswordLogin();
}
```

### Email Not Sending
**Solution**: Check backend logs and SMTP configuration:
1. Verify email credentials in backend
2. Check if 2FA is blocking SMTP access
3. Use app-specific password instead of account password

## Performance Tips

1. **Cache Tokens**: Store JWT in localStorage/sessionStorage
2. **Batch Requests**: Combine multiple API calls when possible
3. **Pagination**: Use `limit` and `offset` for transaction history
4. **Debounce**: Debounce search/filter operations
5. **Error Boundaries**: Wrap API calls in try-catch
6. **Loading States**: Show loading indicators during requests

## Security Tips

1. **Never log tokens**: Don't print tokens in console
2. **HTTPS only**: Always use HTTPS in production
3. **Secure storage**: Use httpOnly cookies if possible instead of localStorage
4. **Input validation**: Validate all user input before sending to API
5. **Rate limiting**: Implement client-side rate limiting
6. **CSRF protection**: Include origin checks for sensitive operations

## Support

For backend-related issues:
1. Check backend logs: `vercel logs --tail`
2. Review error messages in response
3. Verify environment variables are set correctly
4. Check BACKEND_SETUP.md for detailed documentation
