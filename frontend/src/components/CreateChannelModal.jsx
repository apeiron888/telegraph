import { useState, useEffect, useCallback } from 'react';
import { useChatStore } from '../store/chatStore';
import { useAuthStore } from '../store/authStore';
import { api } from '../services/api';
import { X, Hash, Users, AlertCircle, Search } from 'lucide-react';

export default function CreateChannelModal({ onClose }) {
    const [formData, setFormData] = useState({
        type: 'private',
        name: '',
        description: '',
        searchQuery: '',
        selectedUser: null,
    });
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [isSearching, setIsSearching] = useState(false);

    const createChannel = useChatStore((state) => state.createChannel);
    const addMember = useChatStore((state) => state.addMember);
    const deleteChannel = useChatStore((state) => state.deleteChannel);
    const user = useAuthStore((state) => state.user);

    // Debounced search for users
    useEffect(() => {
        if (formData.type !== 'private') return;

        const searchTimeout = setTimeout(async () => {
            if (formData.searchQuery.length < 2) {
                setSearchResults([]);
                return;
            }

            setIsSearching(true);
            try {
                const { data } = await api.get(`/users/search?q=${encodeURIComponent(formData.searchQuery)}`);
                // Filter out the current user from results
                const filteredResults = data.filter(u => u.id !== user?.id);
                setSearchResults(filteredResults);
            } catch (err) {
                console.error('Search error:', err);
                setSearchResults([]);
            } finally {
                setIsSearching(false);
            }
        }, 300);

        return () => clearTimeout(searchTimeout);
    }, [formData.searchQuery, formData.type, user?.id]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            if (formData.type === 'private') {
                if (!formData.selectedUser) {
                    setError('Please select a user to chat with');
                    setIsLoading(false);
                    return;
                }

                // Create private channel with the selected user's name
                const channelData = {
                    type: 'private',
                    name: formData.selectedUser.username,
                    description: 'Private chat',
                    security_label: 'public',
                };

                const channel = await createChannel(channelData);

                try {
                    // Add the selected user as a member
                    await addMember(channel.id, { user_id: formData.selectedUser.id });
                    onClose();
                } catch (addMemberError) {
                    // If adding member fails, delete the channel
                    await deleteChannel(channel.id);
                    throw new Error(addMemberError.response?.data?.error || 'Failed to add user to chat');
                }
            } else {
                // For groups, name is required
                if (!formData.name) {
                    setError('Channel name is required for groups');
                    setIsLoading(false);
                    return;
                }

                const channelData = {
                    type: formData.type,
                    name: formData.name,
                    description: formData.description,
                    security_label: 'public',
                };

                await createChannel(channelData);
                onClose();
            }
        } catch (err) {
            setError(err.message || err.response?.data?.error || 'Failed to create channel');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
                <div className="modal-header">
                    <h2>New Conversation</h2>
                    <button className="icon-btn" onClick={onClose}>
                        <X size={20} />
                    </button>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="modal-body">
                        {error && (
                            <div className="error-message" style={{ marginBottom: '1rem', display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                                <AlertCircle size={16} />
                                <span>{error}</span>
                            </div>
                        )}

                        <div className="form-group">
                            <label>Type</label>
                            <div style={{ display: 'flex', gap: '0.75rem', marginTop: '0.5rem' }}>
                                <button
                                    type="button"
                                    className={`btn ${formData.type === 'private' ? 'btn-primary' : 'btn-secondary'}`}
                                    onClick={() => setFormData({ ...formData, type: 'private', selectedUser: null, searchQuery: '' })}
                                >
                                    <Users size={16} />
                                    Private Chat
                                </button>
                                <button
                                    type="button"
                                    className={`btn ${formData.type === 'group' ? 'btn-primary' : 'btn-secondary'}`}
                                    onClick={() => setFormData({ ...formData, type: 'group' })}
                                >
                                    <Hash size={16} />
                                    Group
                                </button>
                            </div>
                        </div>

                        {formData.type === 'private' ? (
                            <div className="form-group">
                                <label htmlFor="searchQuery">Find User *</label>
                                <div style={{ position: 'relative' }}>
                                    <input
                                        id="searchQuery"
                                        type="text"
                                        className="input"
                                        placeholder="Search by username or email..."
                                        value={formData.selectedUser ? formData.selectedUser.username : formData.searchQuery}
                                        onChange={(e) => {
                                            if (!formData.selectedUser) {
                                                setFormData({ ...formData, searchQuery: e.target.value });
                                            }
                                        }}
                                        onFocus={() => {
                                            if (formData.selectedUser) {
                                                setFormData({ ...formData, selectedUser: null, searchQuery: '' });
                                            }
                                        }}
                                        required
                                        autoComplete="off"
                                    />
                                    <Search size={18} style={{ position: 'absolute', right: '0.75rem', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-tertiary)' }} />
                                </div>

                                {/* Search Results */}
                                {!formData.selectedUser && formData.searchQuery && (
                                    <div style={{
                                        marginTop: '0.5rem',
                                        maxHeight: '200px',
                                        overflowY: 'auto',
                                        border: '1px solid var(--border-color)',
                                        borderRadius: '0.5rem',
                                        backgroundColor: 'var(--bg-secondary)',
                                    }}>
                                        {isSearching ? (
                                            <div style={{ padding: '1rem', textAlign: 'center', color: 'var(--text-secondary)' }}>
                                                Searching...
                                            </div>
                                        ) : searchResults.length === 0 ? (
                                            <div style={{ padding: '1rem', textAlign: 'center', color: 'var(--text-secondary)' }}>
                                                No users found
                                            </div>
                                        ) : (
                                            searchResults.map((result) => (
                                                <div
                                                    key={result.id}
                                                    onClick={() => setFormData({ ...formData, selectedUser: result, searchQuery: result.username })}
                                                    style={{
                                                        padding: '0.75rem',
                                                        cursor: 'pointer',
                                                        borderBottom: '1px solid var(--border-color)',
                                                        transition: 'background-color 0.2s',
                                                    }}
                                                    onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--bg-hover)'}
                                                    onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                                                >
                                                    <div style={{ fontWeight: '500' }}>{result.username}</div>
                                                    <div style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>{result.email}</div>
                                                </div>
                                            ))
                                        )}
                                    </div>
                                )}

                                <small style={{ color: 'var(--text-tertiary)', fontSize: '0.75rem', marginTop: '0.25rem', display: 'block' }}>
                                    Search for a user to start a private conversation
                                </small>
                            </div>
                        ) : (
                            <>
                                <div className="form-group">
                                    <label htmlFor="name">Group Name *</label>
                                    <input
                                        id="name"
                                        type="text"
                                        className="input"
                                        placeholder="Enter group name"
                                        value={formData.name}
                                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                        required
                                    />
                                </div>

                                <div className="form-group">
                                    <label htmlFor="description">Description (optional)</label>
                                    <textarea
                                        id="description"
                                        className="input"
                                        placeholder="What's this group about?"
                                        value={formData.description}
                                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                        rows={3}
                                        style={{ resize: 'vertical' }}
                                    />
                                </div>
                            </>
                        )}
                    </div>

                    <div className="modal-footer">
                        <button type="button" className="btn btn-secondary" onClick={onClose}>
                            Cancel
                        </button>
                        <button type="submit" className="btn btn-primary" disabled={isLoading}>
                            {isLoading ? 'Creating...' : formData.type === 'private' ? 'Start Chat' : 'Create Group'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
