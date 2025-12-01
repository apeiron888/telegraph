import { create } from 'zustand';
import { authAPI } from '../services/api';

export const useAuthStore = create((set) => ({
    user: null,
    isAuthenticated: false,
    isLoading: true,

    login: async (credentials) => {
        const { data } = await authAPI.login(credentials);
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        set({ isAuthenticated: true });
        return data;
    },

    register: async (userData) => {
        const { data } = await authAPI.register(userData);
        return data;
    },

    logout: () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        set({ user: null, isAuthenticated: false });
    },

    loadUser: async () => {
        try {
            const token = localStorage.getItem('access_token');
            if (!token) {
                set({ isLoading: false, isAuthenticated: false });
                return;
            }
            const { data } = await authAPI.getCurrentUser();
            set({ user: data, isAuthenticated: true, isLoading: false });
        } catch (error) {
            localStorage.clear();
            set({ user: null, isAuthenticated: false, isLoading: false });
        }
    },
}));
