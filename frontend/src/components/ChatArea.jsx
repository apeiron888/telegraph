import { useState, useRef, useEffect } from 'react';
import { useChatStore } from '../store/chatStore';
import { useAuthStore } from '../store/authStore';
import { Send, MessageSquare, Hash, MoreVertical, Check, CheckCheck } from 'lucide-react';
import ChannelMenuModal from './ChannelMenuModal';
import { api } from '../services/api';

export default function ChatArea() {
    const [messageText, setMessageText] = useState('');
    const [isSending, setIsSending] = useState(false);
    const [showMenu, setShowMenu] = useState(false);
    const messagesEndRef = useRef(null);
    const typingTimeoutRef = useRef(null);

    const activeChannel = useChatStore((state) => state.activeChannel);
    const messages = useChatStore((state) => state.messages);
    const sendMessage = useChatStore((state) => state.sendMessage);
    const typingUsersMap = useChatStore((state) => state.typingUsers);
    const user = useAuthStore((state) => state.user);

    const typingUsers = activeChannel ? (typingUsersMap[activeChannel.id] || []) : [];

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    // Mark messages as read when viewing
    useEffect(() => {
        if (activeChannel && messages.length > 0) {
            const unreadMessages = messages.filter(m =>
                m.sender_id !== user?.id &&
                !m.read_by?.includes(user?.id)
            );

            unreadMessages.forEach(async (message) => {
                try {
                    await api.post(`/channels/${activeChannel.id}/messages/${message.id}/read`);
                } catch (err) {
                    console.error('Failed to mark as read:', err);
                }
            });
        }
    }, [messages, activeChannel, user?.id]);

    const handleSend = async (e) => {
        e.preventDefault();
        if (!messageText.trim() || !activeChannel || isSending) return;

        setIsSending(true);
        try {
            await sendMessage(activeChannel.id, messageText);
            setMessageText('');
        } catch (error) {
            console.error('Failed to send message:', error);
        } finally {
            setIsSending(false);
        }
    };

    const handleTyping = () => {
        if (!activeChannel) return;

        // Send typing indicator
        api.post(`/channels/${activeChannel.id}/typing`, { typing: true })
            .catch(err => console.error('Failed to send typing:', err));

        // Clear previous timeout
        if (typingTimeoutRef.current) {
            clearTimeout(typingTimeoutRef.current);
        }

        // Stop typing after 3 seconds
        typingTimeoutRef.current = setTimeout(() => {
            api.post(`/channels/${activeChannel.id}/typing`, { typing: false })
                .catch(err => console.error('Failed to stop typing:', err));
        }, 3000);
    };

    // Check if user is owner or admin
    const isOwnerOrAdmin = () => {
        if (!activeChannel || !user) return false;
        if (activeChannel.owner_id === user.id) return true;
        const member = activeChannel.members?.find(m => m.user_id === user.id);
        return member?.role === 'admin';
    };

    // Get message status icon
    const getStatusIcon = (message) => {
        if (message.sender_id !== user?.id) return null;

        if (message.status === 'read' || message.read_by?.length > 0) {
            return <CheckCheck size={14} className="message-status read" />;
        } else if (message.status === 'delivered' || message.delivered_to?.length > 0) {
            return <CheckCheck size={14} className="message-status delivered" />;
        } else {
            return <Check size={14} className="message-status sent" />;
        }
    };

    if (!activeChannel) {
        return (
            <div className="chat-area">
                <div className="empty-state">
                    <div className="empty-state-icon">
                        <MessageSquare size={40} />
                    </div>
                    <h3>Select a channel</h3>
                    <p>Choose a channel from the sidebar to start messaging</p>
                </div>
            </div>
        );
    }

    // Get display name for private chats
    let displayName = activeChannel.name || 'Unnamed Channel';
    if (activeChannel.type === 'private' && activeChannel.members) {
        const otherMember = activeChannel.members.find(m => m.user_id !== user?.id);
        if (otherMember) {
            displayName = otherMember.username || otherMember.email || 'Unknown User';
        }
    }

    return (
        <div className="chat-area">
            <div className="chat-header" onClick={() => setShowMenu(true)} style={{ cursor: 'pointer' }}>
                <div className="chat-header-info">
                    <div className="channel-avatar">
                        <Hash size={20} />
                    </div>
                    <div>
                        <h3>{displayName}</h3>
                        <p>{activeChannel.members?.length || 0} members • Click for info</p>
                    </div>
                </div>
                <div className="chat-header-actions">
                    <button
                        className="icon-btn"
                        onClick={(e) => {
                            e.stopPropagation();
                            setShowMenu(true);
                        }}
                        title="Channel options"
                    >
                        <MoreVertical size={20} />
                    </button>
                </div>
            </div>

            <div className="messages-container">
                {messages.length === 0 ? (
                    <div className="empty-state">
                        <p>No messages yet. Start the conversation!</p>
                    </div>
                ) : (
                    messages.map((message) => {
                        const isOwn = message.sender_id === user?.id;
                        // Decode base64 content for display
                        let decodedContent = message.content;
                        try {
                            decodedContent = atob(message.content);
                        } catch (e) {
                            // If decoding fails, use original content
                        }

                        // Find sender info from channel members or use message sender info
                        let senderName = '';
                        let senderInitial = 'U';

                        if (isOwn) {
                            senderName = user?.username || 'You';
                            senderInitial = senderName[0]?.toUpperCase() || 'Y';
                        } else {
                            // Try to find from channel members first
                            const sender = activeChannel.members?.find(m => m.user_id === message.sender_id);
                            if (sender && sender.username) {
                                senderName = sender.username;
                                senderInitial = senderName[0]?.toUpperCase();
                            } else if (message.sender_username) {
                                // Fallback to message sender_username if available
                                senderName = message.sender_username;
                                senderInitial = senderName[0]?.toUpperCase();
                            } else {
                                senderName = 'User';
                                senderInitial = 'U';
                            }
                        }

                        return (
                            <div key={message.id} className={`message ${isOwn ? 'own' : ''}`}>
                                <div className="message-avatar" title={senderName}>
                                    {senderInitial}
                                </div>
                                <div className="message-content">
                                    {!isOwn && <div className="message-sender">{senderName}</div>}
                                    <div className="message-bubble">
                                        {decodedContent}
                                        {message.edited && <span className="message-edited"> (edited)</span>}
                                    </div>
                                    <div className="message-time">
                                        {new Date(message.timestamp || message.created_at).toLocaleTimeString([], {
                                            hour: '2-digit',
                                            minute: '2-digit',
                                        })}
                                        {getStatusIcon(message)}
                                    </div>
                                </div>
                            </div>
                        );
                    })
                )}
                {typingUsers.length > 0 && (
                    <div className="typing-indicator">
                        <span>{typingUsers.join(', ')} {typingUsers.length === 1 ? 'is' : 'are'} typing...</span>
                    </div>
                )}
                <div ref={messagesEndRef} />
            </div>

            <div className="message-input-container">
                <form onSubmit={handleSend} className="message-input-wrapper">
                    <textarea
                        className="message-input"
                        placeholder="Type a message..."
                        value={messageText}
                        onChange={(e) => {
                            setMessageText(e.target.value);
                            handleTyping();
                        }}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' && !e.shiftKey) {
                                e.preventDefault();
                                handleSend(e);
                            }
                        }}
                        rows={1}
                    />
                    <button type="submit" className="send-btn" disabled={!messageText.trim() || isSending}>
                        <Send size={20} />
                    </button>
                </form>
            </div>

            {showMenu && (
                <ChannelMenuModal
                    channel={activeChannel}
                    onClose={() => setShowMenu(false)}
                    isOwnerOrAdmin={isOwnerOrAdmin()}
                />
            )}
        </div>
    );
}


