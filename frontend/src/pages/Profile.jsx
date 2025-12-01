import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { User, Mail, Phone, ArrowLeft, Save, Shield } from 'lucide-react';
import './Profile.css';

export default function Profile() {
    const navigate = useNavigate();
    const user = useAuthStore((state) => state.user);
    const [formData, setFormData] = useState({
        username: user?.username || '',
        email: user?.email || '',
        phone: user?.phone || '',
        bio: user?.bio || '',
    });
    const [isLoading, setIsLoading] = useState(false);
    const [message, setMessage] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setIsLoading(true);

        try {
            // TODO: Implement update user API call
            setMessage('Profile updated successfully!');
            setTimeout(() => setMessage(''), 3000);
        } catch (error) {
            setMessage('Failed to update profile');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="profile-container">
            <div className="profile-header">
                <button className="icon-btn" onClick={() => navigate('/')}>
                    <ArrowLeft size={20} />
                </button>
                <h1>Profile & Settings</h1>
                <div></div>
            </div>

            <div className="profile-content">
                <div className="profile-card">
                    <div className="profile-avatar-large">
                        {user?.username?.[0]?.toUpperCase() || 'U'}
                    </div>
                    <h2>{user?.username}</h2>
                    <p className="profile-role">
                        <Shield size={14} />
                        {user?.account_type || 'basic'}
                    </p>
                </div>

                <div className="profile-section">
                    <h3>Account Information</h3>
                    <form onSubmit={handleSubmit}>
                        {message && (
                            <div className={`message ${message.includes('success') ? 'success-message' : 'error-message'}`}>
                                {message}
                            </div>
                        )}

                        <div className="form-group">
                            <label htmlFor="username">
                                <User size={16} />
                                Username
                            </label>
                            <input
                                id="username"
                                type="text"
                                className="input"
                                value={formData.username}
                                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="email">
                                <Mail size={16} />
                                Email
                            </label>
                            <input
                                id="email"
                                type="email"
                                className="input"
                                value={formData.email}
                                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                disabled
                            />
                            <small style={{ color: 'var(--text-tertiary)', fontSize: '0.75rem' }}>
                                Email cannot be changed
                            </small>
                        </div>

                        <div className="form-group">
                            <label htmlFor="phone">
                                <Phone size={16} />
                                Phone
                            </label>
                            <input
                                id="phone"
                                type="tel"
                                className="input"
                                value={formData.phone}
                                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="bio">Bio</label>
                            <textarea
                                id="bio"
                                className="input"
                                rows={4}
                                value={formData.bio}
                                onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                                placeholder="Tell us about yourself..."
                            />
                        </div>

                        <button type="submit" className="btn btn-primary btn-block" disabled={isLoading}>
                            <Save size={16} />
                            {isLoading ? 'Saving...' : 'Save Changes'}
                        </button>
                    </form>
                </div>

                <div className="profile-section">
                    <h3>Account Details</h3>
                    <div className="info-grid">
                        <div className="info-item">
                            <span className="info-label">Account Type</span>
                            <span className="info-value">{user?.account_type || 'basic'}</span>
                        </div>
                        <div className="info-item">
                            <span className="info-label">Security Label</span>
                            <span className="info-value">{user?.security_label || 'public'}</span>
                        </div>
                        <div className="info-item">
                            <span className="info-label">Member Since</span>
                            <span className="info-value">
                                {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
