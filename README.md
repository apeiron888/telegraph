# 🔐 Telegraph - Secure Messaging Platform

> **High-security end-to-end encrypted messaging with enterprise-grade access control**

Telegraph is a secure messaging platform built with privacy-first architecture, implementing RBAC, MAC, ABAC access controls, and true end-to-end encryption.

## ✨ Features

- 🔐 **End-to-End Encryption** - AES-256-GCM encryption for all messages
- 👤 **Multi-Factor Authentication** - Email OTP verification
- 🛡️ **Three-Layer Access Control**:
  - **RBAC** - Role-based permissions (Admin/Moderator/Member)
  - **MAC** - Mandatory security labels (Public/Internal/Confidential)
  - **ABAC** - Attribute-based policies (MFA, Premium, Region, etc.)
- 💬 **Channel Types** - Private chats, Group chats, Broadcast channels
- 📝 **Audit Logging** - Complete trail of all security events
- 🔑 **JWT Authentication** - Secure token-based auth with refresh
- 🔒 **Argon2id Hashing** - Industry-standard password security

## 🏗️ Architecture

```
User → API Gateway → Middleware (JWT → RBAC → MAC → ABAC) → Services → PostgreSQL
```

**Backend**: Golang + Chi + PostgreSQL  
**Encryption**: AES-256-GCM (client-side)  
**Auth**: JWT + Refresh Tokens + MFA  

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL 14+
- SMTP server (for MFA emails)

### Backend Setup

```bash
# 1. Database
createdb telegraph
cd backend/migrations
# Run migrations with your tool (goose, migrate, etc.)

# 2. Configuration
cp .env.example backend/.env
# Edit .env with your database and SMTP credentials

# 3. Run server
cd backend
go run cmd/api/main.go
```

**Expected output**:
```
✓ Database connected
✓ Telegraph server running at :8080
✓ Access Control: RBAC + MAC + ABAC enabled
✓ E2EE: Message encryption active
✓ Audit Logging: Enabled
✓ All systems operational
```

## 📚 API Documentation

### Authentication

```bash
# Register
POST /api/v1/users/register
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "SecurePass123!"
}

# Login
POST /api/v1/auth/login
{
  "email": "alice@example.com",
  "password": "SecurePass123!"
}

# MFA Verification
POST /api/v1/auth/mfa/verify
{
  "email": "alice@example.com",
  "code": "123456"
}
```

### Channels

```bash
# Create channel
POST /api/v1/channels
Authorization: Bearer <token>
{
  "type": "group",
  "name": "Project Team",
  "security_label": "internal"
}

# List my channels
GET /api/v1/channels
Authorization: Bearer <token>
```

### Messages

```bash
# Send encrypted message
POST /api/v1/channels/{channelId}/messages
Authorization: Bearer <token>
{
  "content": "<base64-encrypted-blob>",
  "content_type": "text",
  "encryption_meta": {
    "algorithm": "AES-256-GCM",
    "iv": "<base64-iv>"
  }
}

# Get messages (paginated)
GET /api/v1/channels/{channelId}/messages?limit=50&offset=0
Authorization: Bearer <token>
```

## 🗂️ Project Structure

```
telegraph/
├── backend/
│   ├── cmd/api/           # Application entry point
│   ├── internal/
│   │   ├── acl/           # Access Control Layer (RBAC/MAC/ABAC)
│   │   ├── audit/         # Audit logging
│   │   ├── auth/          # Authentication & MFA
│   │   ├── channels/      # Channel management
│   │   ├── config/        # Configuration
│   │   ├── database/      # DB connection
│   │   ├── messages/      # Message handling + E2EE
│   │   ├── middleware/    # HTTP middleware (JWT, ACL)
│   │   └── users/         # User management
│   ├── migrations/        # Database migrations
│   └── .env               # Environment config
└── README.md
```

## 🔒 Security Model

### Access Control Layers

1. **JWT Authentication** - Validates user identity
2. **RBAC** - Role-based permissions (member → moderator → admin)
3. **MAC** - Security clearance levels (public → internal → confidential)
4. **ABAC** - Attribute policies (MFA required, premium only, etc.)

### Encryption

- **At Rest**: PostgreSQL with encrypted message blobs
- **In Transit**: HTTPS/TLS 1.3
- **E2EE**: Client-side AES-256-GCM encryption
- **Keys**: Client-managed, server never sees plaintext

## 📊 Database Schema

**Core Tables**:
- `users` - With RBAC/MAC/ABAC fields
- `channels` - With UUID[] members array
- `messages` - BYTEA encrypted content
- `refresh_tokens` - Session management
- `otps` - MFA codes
- `audit_logs` - Security event trail

## 🧪 Testing

```bash
# Run tests
cd backend
go test ./...

# Test specific module
go test ./internal/acl/...
go test ./internal/messages/...
```

**Security Tests**:
- SQL Injection protection
- XSS prevention
- Token expiry
- MAC bypass attempts
- Permission escalation

## 📈 Next Steps

- [ ] Build React frontend with Web Crypto API
- [ ] Add WebSocket for real-time messaging
- [ ] Implement file upload with encryption
- [ ] Create Flutter mobile app
- [ ] Add rate limiting middleware
- [ ] Deploy to production (Docker/K8s)

## 📖 Documentation

- [Implementation Plan](./implementation_plan.md) - Detailed technical plan
- [Walkthrough](./walkthrough.md) - Complete implementation guide
- [SRADD](./SRADD.md) - Original requirements document

## 🤝 Contributing

This project follows secure coding practices:
- All DB queries use parameterized statements
- Input validation on all endpoints
- Least privilege principle enforced
- Audit logging for all actions

## 📝 License

Apache 2.0 - See LICENSE file

---

**Built with ❤️ for privacy and security**
