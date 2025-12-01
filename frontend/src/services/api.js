import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token expired, try to refresh
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          });
          localStorage.setItem('access_token', data.access_token);
          error.config.headers.Authorization = `Bearer ${data.access_token}`;
          return axios(error.config);
        } catch (refreshError) {
          localStorage.clear();
          window.location.href = '/login';
        }
      } else {
        localStorage.clear();
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  register: (data) => api.post('/users/register', data),
  login: (data) => api.post('/auth/login', data),
  getCurrentUser: () => api.get('/users/me'),
};

export const channelsAPI = {
  getAll: () => api.get('/channels'),
  getById: (id) => api.get(`/channels/${id}`),
  create: (data) => api.post('/channels', data),
  addMember: (channelId, data) => api.post(`/channels/${channelId}/members`, data),
  removeMember: (channelId, userId) => api.delete(`/channels/${channelId}/members/${userId}`),
  promoteMember: (channelId, userId) => api.post(`/channels/${channelId}/members/${userId}/promote`),
  demoteMember: (channelId, userId) => api.post(`/channels/${channelId}/members/${userId}/demote`),
  delete: (channelId) => api.delete(`/channels/${channelId}`),
};

export const messagesAPI = {
  getMessages: (channelId, params) => api.get(`/channels/${channelId}/messages`, { params }),
  send: (channelId, data) => api.post(`/channels/${channelId}/messages`, data),
  delete: (messageId) => api.delete(`/messages/${messageId}`),
  getUnreadCounts: () => api.get('/unread'),
};

export { api };
export default api;
