#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Starting Telegraph Test Suite...${NC}"

# Helper to extract JSON field
get_json_field() {
    echo "$1" | python3 -c "import sys, json; print(json.load(sys.stdin).get('$2', ''))"
}

# 1. Register User
echo -e "\n${GREEN}[1] Registering User (alice_test)...${NC}"
REGISTER_RES=$(curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice_test_'$(date +%s)'", "email": "alice_'$(date +%s)'@test.com", "password": "Password123!"}')
echo "Response: $REGISTER_RES"

EMAIL=$(get_json_field "$REGISTER_RES" "email")

# 2. Login
echo -e "\n${GREEN}[2] Logging In...${NC}"
LOGIN_RES=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"Password123!\"}")
TOKEN=$(get_json_field "$LOGIN_RES" "access_token")
echo "Token obtained."

# 3. Invalid Login
echo -e "\n${GREEN}[3] Testing Invalid Login...${NC}"
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"WrongPass\"}"
echo ""

# 4. Create Valid Channel
echo -e "\n${GREEN}[4] Creating Valid Channel...${NC}"
CHANNEL_RES=$(curl -s -X POST http://localhost:8080/api/v1/channels \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "group",
    "name": "Project Alpha",
    "description": "Top secret discussion",
    "members": [], 
    "security_label": "public"
  }')
echo "Response: $CHANNEL_RES"
CHANNEL_ID=$(get_json_field "$CHANNEL_RES" "id")

# 5. Add Member by Email
echo -e "\n${GREEN}[5] Adding Member by Email...${NC}"
# Register a second user first
USER2_RES=$(curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username": "bob_test_'$(date +%s)'", "email": "bob_'$(date +%s)'@test.com", "password": "Password123!"}')
USER2_EMAIL=$(get_json_field "$USER2_RES" "email")

curl -s -X POST http://localhost:8080/api/v1/channels/$CHANNEL_ID/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$USER2_EMAIL\"}"
echo ""

# 6. Create Invalid Channel (Bad UUID)
echo -e "\n${GREEN}[6] Testing Invalid Channel (Bad UUID)...${NC}"
curl -s -X POST http://localhost:8080/api/v1/channels \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "group",
    "name": "Broken Channel",
    "members": ["not-a-uuid"]
  }'
echo ""

# 6. Send Message
echo -e "\n${GREEN}[6] Sending Encrypted Message...${NC}"
curl -s -X POST http://localhost:8080/api/v1/channels/$CHANNEL_ID/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "SGVsbG8gV29ybGQ=", 
    "content_type": "text",
    "encryption_meta": {"alg": "AES-GCM", "iv": "MTIz"}
  }'
echo ""

# 7. Verify Audit Log
echo -e "\n${GREEN}[7] Verifying Audit Log File...${NC}"
if [ -f audit.log ]; then
    echo "Audit log exists. Last 5 lines:"
    tail -n 5 audit.log
else
    echo -e "${RED}Audit log file not found!${NC}"
fi

echo -e "\n${GREEN}Test Suite Completed.${NC}"
