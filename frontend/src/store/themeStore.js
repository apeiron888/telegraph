import { create } from 'zustand';

export const useThemeStore = create((set) => ({
    theme: localStorage.getItem('theme') || 'light',

    toggleTheme: () => set((state) => {
        const newTheme = state.theme === 'light' ? 'dark' : 'light';
        localStorage.setItem('theme', newTheme);
        document.documentElement.classList.toggle('dark', newTheme === 'dark');
        return { theme: newTheme };
    }),

    initTheme: () => {
        const theme = localStorage.getItem('theme') || 'light';
        document.documentElement.classList.toggle('dark', theme === 'dark');
        set({ theme });
    },
}));
