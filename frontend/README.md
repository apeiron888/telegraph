# Telegraph Frontend

A modern, Telegram-like messaging application built with React and Vite.

## Features

- 🔐 **Authentication**: Secure login and registration
- 💬 **Real-time Messaging**: Send and receive messages in channels
- 📱 **Channel Management**: Create and manage group channels
- 🌓 **Dark/Light Theme**: Toggle between themes
- 🎨 **Telegram-like UI**: Clean, modern interface inspired by Telegram
- 🔒 **End-to-End Encryption**: Messages are encrypted before sending
- 👥 **User Management**: Add members by email or phone
- 🎭 **Role-Based Access**: Owner, Admin, and Member roles

## Getting Started

### Prerequisites

- Node.js 16+ installed
- Backend server running on `http://localhost:8080`

### Installation

```bash
cd frontend
npm install
```

### Running the Application

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## Usage Guide

### 1. Registration

1. Open `http://localhost:5173` in your browser
2. Click "Sign up" link
3. Fill in the registration form:
   - Username
   - Email
   - Phone (optional)
   - Password (min 8 characters)
4. Click "Create Account"

### 2. Login

1. Enter your email and password
2. Click "Sign In"
3. You'll be redirected to the main chat interface

### 3. Creating a Channel

1. Click the "+" button in the sidebar header
2. Choose channel type (Group or Broadcast)
3. Enter channel name (required)
4. Add description (optional)
5. Select security label
6. Click "Create Channel"

### 4. Sending Messages

1. Select a channel from the sidebar
2. Type your message in the input field at the bottom
3. Press Enter or click the send button
4. Messages are automatically encrypted before sending

### 5. Theme Toggle

- Click the Moon/Sun icon in the sidebar header to toggle between dark and light themes
- Your preference is saved in localStorage

## Project Structure

```
frontend/
├── src/
│   ├── components/          # Reusable components
│   │   ├── Sidebar.jsx     # Channel list sidebar
│   │   ├── ChatArea.jsx    # Message display and input
│   │   └── CreateChannelModal.jsx
│   ├── pages/              # Page components
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   └── Chat.jsx
│   ├── services/           # API services
│   │   └── api.js          # Axios instance and API calls
│   ├── store/              # Zustand stores
│   │   ├── authStore.js    # Authentication state
│   │   ├── chatStore.js    # Chat and channels state
│   │   └── themeStore.js   # Theme state
│   ├── App.jsx             # Main app component
│   ├── main.jsx            # Entry point
│   └── index.css           # Global styles
├── index.html
├── package.json
└── vite.config.js
```

## API Integration

The frontend communicates with the backend API at `http://localhost:8080/api/v1`. Key endpoints:

- `POST /users/register` - User registration
- `POST /auth/login` - User login
- `GET /users/me` - Get current user
- `GET /channels` - Get all channels
- `POST /channels` - Create channel
- `GET /channels/:id/messages` - Get messages
- `POST /channels/:id/messages` - Send message

## State Management

The application uses Zustand for state management:

- **authStore**: Manages authentication state and user data
- **chatStore**: Manages channels, messages, and chat operations
- **themeStore**: Manages theme (dark/light mode)

## Styling

- Custom CSS with CSS variables for theming
- Telegram-inspired design
- Responsive layout
- Smooth animations and transitions

## Security Features

- JWT token authentication with automatic refresh
- Encrypted message storage (base64 encoding for demo)
- Secure password requirements
- Protected routes
- XSS protection

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Troubleshooting

### Backend Connection Issues

If you see connection errors:
1. Ensure the backend is running on `http://localhost:8080`
2. Check CORS settings in the backend
3. Verify the API base URL in `src/services/api.js`

### Authentication Issues

If login fails:
1. Clear localStorage: `localStorage.clear()`
2. Refresh the page
3. Try registering a new account

### Theme Not Persisting

The theme is saved in localStorage. If it's not persisting:
1. Check browser console for errors
2. Ensure localStorage is enabled in your browser

## Development

### Building for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

### Preview Production Build

```bash
npm run preview
```

## License

MIT
