import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { UserPlus, Mail, Lock, User, Phone, AlertCircle, CheckCircle } from 'lucide-react';
import './Auth.css';

export default function Register() {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        phone: '',
        password: '',
        confirmPassword: '',
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const register = useAuthStore((state) => state.register);
    const navigate = useNavigate();

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess(false);

        if (formData.password !== formData.confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        if (formData.password.length < 8) {
            setError('Password must be at least 8 characters long');
            return;
        }

        setIsLoading(true);

        try {
            await register({
                username: formData.username,
                email: formData.email,
                phone: formData.phone,
                password: formData.password,
            });
            setSuccess(true);
            setTimeout(() => navigate('/login'), 2000);
        } catch (err) {
            setError(err.response?.data?.error || 'Registration failed. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-card fade-in">
                <div className="auth-header">
                    <div className="auth-icon">
                        <UserPlus size={32} />
                    </div>
                    <h1>Create Account</h1>
                    <p>Join Telegraph secure messaging</p>
                </div>

                <form onSubmit={handleSubmit} className="auth-form">
                    {error && (
                        <div className="error-message">
                            <AlertCircle size={16} />
                            <span>{error}</span>
                        </div>
                    )}

                    {success && (
                        <div className="success-message">
                            <CheckCircle size={16} />
                            <span>Account created! Redirecting to login...</span>
                        </div>
                    )}

                    <div className="form-group">
                        <label htmlFor="username">
                            <User size={16} />
                            Username
                        </label>
                        <input
                            id="username"
                            name="username"
                            type="text"
                            className="input"
                            placeholder="Choose a username"
                            value={formData.username}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="email">
                            <Mail size={16} />
                            Email
                        </label>
                        <input
                            id="email"
                            name="email"
                            type="email"
                            className="input"
                            placeholder="Enter your email"
                            value={formData.email}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="phone">
                            <Phone size={16} />
                            Phone (optional)
                        </label>
                        <input
                            id="phone"
                            name="phone"
                            type="tel"
                            className="input"
                            placeholder="+1234567890"
                            value={formData.phone}
                            onChange={handleChange}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">
                            <Lock size={16} />
                            Password
                        </label>
                        <input
                            id="password"
                            name="password"
                            type="password"
                            className="input"
                            placeholder="Create a password"
                            value={formData.password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="confirmPassword">
                            <Lock size={16} />
                            Confirm Password
                        </label>
                        <input
                            id="confirmPassword"
                            name="confirmPassword"
                            type="password"
                            className="input"
                            placeholder="Confirm your password"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <button type="submit" className="btn btn-primary btn-block" disabled={isLoading}>
                        {isLoading ? (
                            <>
                                <div className="spinner-small"></div>
                                Creating account...
                            </>
                        ) : (
                            <>
                                <UserPlus size={16} />
                                Create Account
                            </>
                        )}
                    </button>

                    <div className="auth-footer">
                        <p>
                            Already have an account?{' '}
                            <Link to="/login" className="auth-link">
                                Sign in
                            </Link>
                        </p>
                    </div>
                </form>
            </div>
        </div>
    );
}
