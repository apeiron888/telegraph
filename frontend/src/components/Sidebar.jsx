import { useState } from 'react';
import { useChatStore } from '../store/chatStore';
import { useAuthStore } from '../store/authStore';
import { useThemeStore } from '../store/themeStore';
import { useNavigate } from 'react-router-dom';
import { Plus, LogOut, Moon, Sun, Hash, User as UserIcon } from 'lucide-react';
import CreateChannelModal from './CreateChannelModal';

export default function Sidebar() {
    const [showCreateModal, setShowCreateModal] = useState(false);
    const navigate = useNavigate();
    const channels = useChatStore((state) => state.channels);
    const unreadCounts = useChatStore((state) => state.unreadCounts);
    const activeChannel = useChatStore((state) => state.activeChannel);
    const setActiveChannel = useChatStore((state) => state.setActiveChannel);
    const logout = useAuthStore((state) => state.logout);
    const user = useAuthStore((state) => state.user);
    const theme = useThemeStore((state) => state.theme);
    const toggleTheme = useThemeStore((state) => state.toggleTheme);

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <div className="sidebar">
            <div className="sidebar-header">
                <h2>Telegraph</h2>
                <div className="sidebar-actions">
                    <button className="icon-btn" onClick={() => navigate('/profile')} title="Profile">
                        <UserIcon size={20} />
                    </button>
                    <button className="icon-btn" onClick={toggleTheme} title="Toggle theme">
                        {theme === 'light' ? <Moon size={20} /> : <Sun size={20} />}
                    </button>
                    <button className="icon-btn" onClick={() => setShowCreateModal(true)} title="New conversation">
                        <Plus size={20} />
                    </button>
                    <button className="icon-btn" onClick={handleLogout} title="Logout">
                        <LogOut size={20} />
                    </button>
                </div>
            </div>

            <div className="channels-list">
                {channels.length === 0 ? (
                    <div className="empty-state" style={{ padding: '2rem 1rem' }}>
                        <p>No channels yet</p>
                        <button className="btn btn-primary" onClick={() => setShowCreateModal(true)}>
                            <Plus size={16} />
                            Start Conversation
                        </button>
                    </div>
                ) : (
                    channels.map((channel) => {
                        // For private chats, show the other person's name
                        let displayName = channel.name || 'Unnamed Channel';
                        if (channel.type === 'private' && channel.members) {
                            const otherMember = channel.members.find(m => m.user_id !== user?.id);
                            if (otherMember) {
                                displayName = otherMember.username || otherMember.email || 'Unknown User';
                            }
                        }

                        return (
                            <div
                                key={channel.id}
                                className={`channel-item ${activeChannel?.id === channel.id ? 'active' : ''}`}
                                onClick={() => setActiveChannel(channel)}
                            >
                                <div className="channel-avatar">
                                    {channel.type === 'group' ? <Hash size={24} /> : displayName[0]?.toUpperCase() || 'C'}
                                </div>
                                <div className="channel-info">
                                    <div className="channel-name">
                                        {displayName}
                                        {unreadCounts[channel.id] > 0 && (
                                            <span className="unread-badge">{unreadCounts[channel.id]}</span>
                                        )}
                                    </div>
                                    <div className="channel-preview">
                                        {channel.type === 'private' ? 'Private Chat' : `${channel.members?.length || 0} members`}
                                    </div>
                                </div>
                            </div>
                        );
                    })
                )}
            </div>

            {showCreateModal && <CreateChannelModal onClose={() => setShowCreateModal(false)} />}
        </div>
    );
}
