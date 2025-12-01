import { useEffect } from 'react';
import { useChatStore } from '../store/chatStore';
import { useAuthStore } from '../store/authStore';
import Sidebar from '../components/Sidebar';
import ChatArea from '../components/ChatArea';
import './Chat.css';

export default function Chat() {
    const loadChannels = useChatStore((state) => state.loadChannels);
    const initializeWebSocket = useChatStore((state) => state.initializeWebSocket);
    const disconnectWebSocket = useChatStore((state) => state.disconnectWebSocket);
    const user = useAuthStore((state) => state.user);

    useEffect(() => {
        if (user) {
            loadChannels();
            initializeWebSocket();
        }

        return () => {
            disconnectWebSocket();
        };
    }, [user, loadChannels, initializeWebSocket, disconnectWebSocket]);

    return (
        <div className="chat-container">
            <Sidebar />
            <ChatArea />
        </div>
    );
}
