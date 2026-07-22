#!/bin/bash

# GoBank Backend - API Testing Script
# This script tests all backend endpoints locally
# Usage: ./test-endpoints.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:3000}"
ACCOUNT_NUMBER=""
JWT_TOKEN=""

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   GoBank Backend - API Test Suite      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if server is running
echo -e "${YELLOW}[1] Testing server connectivity...${NC}"
if ! curl -s "${API_URL}/health" > /dev/null 2>&1; then
    echo -e "${RED}✗ Server not responding at ${API_URL}${NC}"
    echo "Make sure backend is running: cd backend && go run ."
    exit 1
fi
echo -e "${GREEN}✓ Server is running${NC}"
echo ""

# Test health endpoint
echo -e "${YELLOW}[2] Testing /health endpoint...${NC}"
HEALTH=$(curl -s "${API_URL}/health")
echo "Response: $HEALTH"
if echo "$HEALTH" | grep -q "ok"; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
fi
echo ""

# Test create account
echo -e "${YELLOW}[3] Testing POST /account (Create Account)...${NC}"
CREATE_RESPONSE=$(curl -s -X POST "${API_URL}/account" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "password": "TestPassword123!"
  }')
echo "Response: $CREATE_RESPONSE"

# Extract account number from response
if echo "$CREATE_RESPONSE" | grep -q "account_number"; then
    ACCOUNT_NUMBER=$(echo "$CREATE_RESPONSE" | grep -o '"account_number":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}✓ Account created with number: $ACCOUNT_NUMBER${NC}"
else
    echo -e "${RED}✗ Account creation failed${NC}"
fi
echo ""

# Test login
echo -e "${YELLOW}[4] Testing POST /login (Login)...${NC}"
if [ -n "$ACCOUNT_NUMBER" ]; then
    LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"account_number\": $ACCOUNT_NUMBER,
        \"password\": \"TestPassword123!\"
      }")
    echo "Response: $LOGIN_RESPONSE"
    
    if echo "$LOGIN_RESPONSE" | grep -q "token"; then
        JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
        echo -e "${GREEN}✓ Login successful${NC}"
        echo "JWT Token: ${JWT_TOKEN:0:20}..."
    else
        echo -e "${RED}✗ Login failed${NC}"
    fi
else
    echo -e "${YELLOW}⊘ Skipping (no account created)${NC}"
fi
echo ""

# Test get account
echo -e "${YELLOW}[5] Testing GET /account/{id} (Get Account)...${NC}"
if [ -n "$JWT_TOKEN" ] && [ -n "$ACCOUNT_NUMBER" ]; then
    GET_RESPONSE=$(curl -s -X GET "${API_URL}/account/${ACCOUNT_NUMBER}" \
      -H "Authorization: Bearer ${JWT_TOKEN}")
    echo "Response: $GET_RESPONSE"
    
    if echo "$GET_RESPONSE" | grep -q "first_name"; then
        echo -e "${GREEN}✓ Account retrieval successful${NC}"
    else
        echo -e "${RED}✗ Account retrieval failed${NC}"
    fi
else
    echo -e "${YELLOW}⊘ Skipping (no token available)${NC}"
fi
echo ""

# Test transfer (create another account first)
echo -e "${YELLOW}[6] Testing POST /transfer (Transfer Funds)...${NC}"
CREATE_RESPONSE2=$(curl -s -X POST "${API_URL}/account" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Recipient",
    "last_name": "User",
    "email": "recipient@example.com",
    "password": "RecipientPass123!"
  }')

if echo "$CREATE_RESPONSE2" | grep -q "account_number"; then
    RECIPIENT_ACCOUNT=$(echo "$CREATE_RESPONSE2" | grep -o '"account_number":[0-9]*' | grep -o '[0-9]*')
    
    if [ -n "$JWT_TOKEN" ]; then
        TRANSFER_RESPONSE=$(curl -s -X POST "${API_URL}/transfer" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer ${JWT_TOKEN}" \
          -d "{
            \"to_account\": $RECIPIENT_ACCOUNT,
            \"amount\": 100
          }")
        echo "Response: $TRANSFER_RESPONSE"
        
        if echo "$TRANSFER_RESPONSE" | grep -q "status"; then
            echo -e "${GREEN}✓ Transfer request sent${NC}"
        else
            echo -e "${YELLOW}⊘ Transfer might need additional setup${NC}"
        fi
    fi
else
    echo -e "${YELLOW}⊘ Could not create recipient account${NC}"
fi
echo ""

# Test CORS headers
echo -e "${YELLOW}[7] Testing CORS Headers...${NC}"
CORS_RESPONSE=$(curl -s -i -X OPTIONS "${API_URL}/health" 2>&1)
if echo "$CORS_RESPONSE" | grep -q "Access-Control-Allow-Origin"; then
    echo -e "${GREEN}✓ CORS headers present${NC}"
else
    echo -e "${YELLOW}⊘ CORS headers might not be configured${NC}"
fi
echo ""

# Summary
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         Test Summary                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✓ Basic connectivity tests completed${NC}"
echo ""
echo "Next steps:"
echo "1. Verify all endpoints returned expected responses"
echo "2. Check for any error messages above"
echo "3. Test with your frontend application"
echo "4. Review backend logs for any issues"
echo ""

