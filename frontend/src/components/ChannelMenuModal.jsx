import { useState } from 'react';
import { useChatStore } from '../store/chatStore';
import { useAuthStore } from '../store/authStore';
import { X, UserPlus, Trash2, Shield, Edit, Users, ShieldCheck, ShieldOff, Save } from 'lucide-react';

export default function ChannelMenuModal({ channel, onClose, isOwnerOrAdmin }) {
    const [activeTab, setActiveTab] = useState('info');
    const [memberEmail, setMemberEmail] = useState('');
    const [newChannelName, setNewChannelName] = useState(channel.name || '');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const addMember = useChatStore((state) => state.addMember);
    const updateChannelName = useChatStore((state) => state.updateChannelName);
    const deleteChannel = useChatStore((state) => state.deleteChannel);
    const promoteMember = useChatStore((state) => state.promoteMember);
    const demoteMember = useChatStore((state) => state.demoteMember);
    const removeMember = useChatStore((state) => state.removeMember);
    const user = useAuthStore((state) => state.user);

    const isOwner = channel.owner_id === user?.id;

    const handleAddMember = async (e) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            await addMember(channel.id, { email: memberEmail });
            setMemberEmail('');
            // Reload channel to get updated members
            await useChatStore.getState().loadChannels();
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to add member');
        } finally {
            setIsLoading(false);
        }
    };

    const handleUpdateName = async (e) => {
        e.preventDefault();
        if (!newChannelName.trim()) return;

        setIsLoading(true);
        try {
            await updateChannelName(channel.id, newChannelName);
            onClose();
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to update channel name');
        } finally {
            setIsLoading(false);
        }
    };

    const handleDeleteChannel = async () => {
        if (!confirm('Are you sure you want to delete this channel? This action cannot be undone.')) {
            return;
        }

        setIsLoading(true);
        try {
            await deleteChannel(channel.id);
            onClose();
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to delete channel');
            setIsLoading(false);
        }
    };

    const handlePromote = async (memberId) => {
        try {
            await promoteMember(channel.id, memberId);
            await useChatStore.getState().loadChannels();
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to promote member');
        }
    };

    const handleDemote = async (memberId) => {
        try {
            await demoteMember(channel.id, memberId);
            await useChatStore.getState().loadChannels();
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to demote member');
        }
    };

    const handleRemove = async (memberId) => {
        if (!confirm('Are you sure you want to remove this member?')) return;

        try {
            await removeMember(channel.id, memberId);
            await useChatStore.getState().loadChannels();
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to remove member');
        }
    };

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '600px' }}>
                <div className="modal-header">
                    <h2>Channel Settings</h2>
                    <button className="icon-btn" onClick={onClose}>
                        <X size={20} />
                    </button>
                </div>

                <div className="modal-tabs">
                    <button
                        className={`tab-btn ${activeTab === 'info' ? 'active' : ''}`}
                        onClick={() => setActiveTab('info')}
                    >
                        <Users size={16} />
                        Info
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'members' ? 'active' : ''}`}
                        onClick={() => setActiveTab('members')}
                    >
                        <Users size={16} />
                        Members
                    </button>
                    {isOwnerOrAdmin && (
                        <>
                            <button
                                className={`tab-btn ${activeTab === 'settings' ? 'active' : ''}`}
                                onClick={() => setActiveTab('settings')}
                            >
                                <Edit size={16} />
                                Settings
                            </button>
                            {isOwner && (
                                <button
                                    className={`tab-btn ${activeTab === 'danger' ? 'active' : ''}`}
                                    onClick={() => setActiveTab('danger')}
                                >
                                    <Trash2 size={16} />
                                    Delete
                                </button>
                            )}
                        </>
                    )}
                </div>

                <div className="modal-body">
                    {error && (
                        <div className="error-message" style={{ marginBottom: '1rem' }}>
                            {error}
                        </div>
                    )}

                    {activeTab === 'info' && (
                        <div>
                            <div className="channel-info-section">
                                <h4>Channel Information</h4>
                                <div className="info-grid">
                                    <div className="info-item">
                                        <span className="info-label">Name</span>
                                        <span className="info-value">{channel.name || 'Unnamed'}</span>
                                    </div>
                                    <div className="info-item">
                                        <span className="info-label">Type</span>
                                        <span className="info-value">{channel.type}</span>
                                    </div>
                                    <div className="info-item">
                                        <span className="info-label">Members</span>
                                        <span className="info-value">{channel.members?.length || 0}</span>
                                    </div>
                                    <div className="info-item">
                                        <span className="info-label">Security</span>
                                        <span className="info-value">{channel.security_label || 'public'}</span>
                                    </div>
                                </div>
                                {channel.description && (
                                    <div style={{ marginTop: '1rem' }}>
                                        <span className="info-label">Description</span>
                                        <p style={{ marginTop: '0.5rem', color: 'var(--text-primary)' }}>
                                            {channel.description}
                                        </p>
                                    </div>
                                )}
                            </div>

                            {isOwnerOrAdmin && (
                                <div style={{ marginTop: '1.5rem', display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
                                    <button
                                        className="btn btn-primary"
                                        onClick={() => setActiveTab('members')}
                                    >
                                        <UserPlus size={16} />
                                        Manage Members
                                    </button>
                                    <button
                                        className="btn btn-secondary"
                                        onClick={() => setActiveTab('settings')}
                                    >
                                        <Edit size={16} />
                                        Update Channel
                                    </button>
                                    {isOwner && (
                                        <button
                                            className="btn btn-danger"
                                            onClick={() => setActiveTab('danger')}
                                        >
                                            <Trash2 size={16} />
                                            Delete Channel
                                        </button>
                                    )}
                                </div>
                            )}
                        </div>
                    )}
                    {activeTab === 'members' && (
                        <div>
                            {isOwnerOrAdmin && (
                                <form onSubmit={handleAddMember} style={{ marginBottom: '1.5rem' }}>
                                    <div className="form-group">
                                        <label>Add Member</label>
                                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                                            <input
                                                type="email"
                                                className="input"
                                                placeholder="Enter email address"
                                                value={memberEmail}
                                                onChange={(e) => setMemberEmail(e.target.value)}
                                                required
                                            />
                                            <button type="submit" className="btn btn-primary" disabled={isLoading}>
                                                <UserPlus size={16} />
                                                Add
                                            </button>
                                        </div>
                                    </div>
                                </form>
                            )}

                            <div className="members-list">
                                <h4 style={{ marginBottom: '0.75rem', fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
                                    {channel.members?.length || 0} Members
                                </h4>
                                {channel.members?.map((member) => (
                                    <div key={member.user_id} className="member-item">
                                        <div className="member-avatar">
                                            {member.username?.[0]?.toUpperCase() || 'U'}
                                        </div>
                                        <div className="member-info">
                                            <div className="member-name">{member.username || 'Unknown User'}</div>
                                            <div className="member-role">
                                                {member.role === 'owner' && <><ShieldCheck size={12} /> Owner</>}
                                                {member.role === 'admin' && <><Shield size={12} /> Admin</>}
                                                {member.role === 'member' && 'Member'}
                                            </div>
                                        </div>
                                        {isOwner && member.user_id !== user?.id && (
                                            <div className="member-actions">
                                                {member.role === 'member' && (
                                                    <button
                                                        className="btn btn-secondary"
                                                        style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem' }}
                                                        onClick={() => handlePromote(member.user_id)}
                                                    >
                                                        Promote
                                                    </button>
                                                )}
                                                {member.role === 'admin' && (
                                                    <button
                                                        className="btn btn-secondary"
                                                        style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem' }}
                                                        onClick={() => handleDemote(member.user_id)}
                                                    >
                                                        Demote
                                                    </button>
                                                )}
                                                <button
                                                    className="btn btn-danger"
                                                    style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem' }}
                                                    onClick={() => handleRemove(member.user_id)}
                                                >
                                                    Remove
                                                </button>
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    {activeTab === 'settings' && isOwnerOrAdmin && (
                        <form onSubmit={handleUpdateName}>
                            <div className="form-group">
                                <label htmlFor="channelName">Channel Name</label>
                                <input
                                    id="channelName"
                                    type="text"
                                    className="input"
                                    value={newChannelName}
                                    onChange={(e) => setNewChannelName(e.target.value)}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label htmlFor="channelDesc">Description</label>
                                <textarea
                                    id="channelDesc"
                                    className="input"
                                    rows={3}
                                    placeholder="Enter channel description..."
                                    defaultValue={channel.description || ''}
                                />
                            </div>
                            <button type="submit" className="btn btn-primary" disabled={isLoading}>
                                <Save size={16} />
                                Save Changes
                            </button>
                        </form>
                    )}

                    {activeTab === 'danger' && isOwner && (
                        <div>
                            <h4 style={{ marginBottom: '0.75rem', color: 'var(--danger)' }}>Delete Channel</h4>
                            <p style={{ marginBottom: '1rem', fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
                                Once you delete a channel, there is no going back. Please be certain.
                            </p>
                            <button
                                className="btn btn-danger"
                                onClick={handleDeleteChannel}
                                disabled={isLoading}
                            >
                                <Trash2 size={16} />
                                Delete Channel
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
