import { create } from 'zustand';
import { channelsAPI, messagesAPI } from '../services/api';
import { wsService } from '../services/websocket';

export const useChatStore = create((set, get) => ({
    channels: [],
    activeChannel: null,
    messages: [],
    typingUsers: {},
    unreadCounts: {},
    isLoadingChannels: false,
    isLoadingMessages: false,

    initializeWebSocket: () => {
        // Minimal implementation to allow build
    },

    disconnectWebSocket: () => {
        wsService.disconnect();
    },

    loadChannels: async () => {
        set({ isLoadingChannels: true });
        try {
            const { data } = await channelsAPI.getAll();
            set({ channels: data || [], isLoadingChannels: false });
        } catch (error) {
            set({ isLoadingChannels: false });
        }
    },

    setActiveChannel: async (channel) => { },
    sendMessage: async (channelId, content) => { },
    createChannel: async (channelData) => { },
    addMember: async (channelId, memberData) => { },
    removeMember: async (channelId, userId) => { },
    promoteMember: async (channelId, userId) => { },
    demoteMember: async (channelId, userId) => { },
    updateChannelName: async (channelId, name) => { },
    deleteChannel: async (channelId) => { },
}));
