# Telegraph API Testing

## Postman Collection Setup

### 1. Import Collection
1. Open Postman
2. Click **Import** button
3. Select `telegraph_api_collection.json`
4. Click **Import**

### 2. Import Environment
1. Click the **Environments** icon (gear icon in top right)
2. Click **Import**
3. Select `telegraph_environment.json`
4. Click **Import**
5. Select "Telegraph Environment" from the environment dropdown

### 3. Usage Instructions

#### First Time Setup
1. **Register User**: Run "Register User" request first
   - This will automatically save the `user_id` and `user_email` to environment
2. **Login**: Run "Login" request
   - This will automatically save the `access_token` and `refresh_token`

#### Testing Flow
The collection is designed to work sequentially:

1. **Authentication**
   - Register User → Login → (Optional) Refresh Token

2. **Channels**
   - Create Group Channel → Get My Channels → Add Members → Promote/Demote

3. **Messages**
   - Send Message → Get Messages → Delete Message

#### Environment Variables
The following variables are automatically set by test scripts:
- `access_token` - JWT token for authentication
- `refresh_token` - Token for refreshing access token
- `user_id` - Current user's ID
- `user_email` - Current user's email
- `channel_id` - Last created channel ID
- `member_id` - Member ID for promote/demote operations
- `message_id` - Last sent message ID

#### Manual Variables
You may need to manually set:
- `member_id` - Copy from a user's ID when testing promote/demote
- `base_url` - Change if running on different host/port

### 4. Testing Features

#### Channel Roles (DAC)
1. Create a channel (you become Owner)
2. Add a member by email/phone
3. Copy the member's user_id to `member_id` environment variable
4. Use "Promote Member to Admin" to make them admin
5. Use "Demote Admin to Member" to remove admin privileges

#### User Lookup
- Use "Add Member by Email" with any registered user's email
- Use "Add Member by Phone" with any registered user's phone number

#### Account Limits
- Create 6 channels with a basic account to test the 5-channel limit

### 5. Expected Responses

#### Success Responses
- **Register**: 200 OK with user object
- **Login**: 200 OK with `access_token` and `refresh_token`
- **Create Channel**: 201 Created with channel object
- **Send Message**: 200/201 OK with message object

#### Error Responses
- **401 Unauthorized**: Invalid or missing token
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource doesn't exist
- **400 Bad Request**: Invalid input data

### 6. Troubleshooting

**Token not working?**
- Make sure you've run the Login request
- Check that "Telegraph Environment" is selected
- Verify `access_token` is set in environment variables

**Channel ID not found?**
- Run "Create Group Channel" first
- Check `channel_id` is set in environment variables

**Member operations failing?**
- Ensure you're the channel owner (for promote/demote)
- Verify `member_id` is set correctly
- Check the member is actually in the channel
