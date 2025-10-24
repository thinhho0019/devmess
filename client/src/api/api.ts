import axios, { type AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';

// Types for API responses
interface AuthResponse {
  access_token: string;
  refresh_token?: string;
}

interface QueueItem {
  resolve: (token: string) => void;
  // eslint-disable-next-line
  reject: (error: any) => void;
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
});

console.log("🔧 API baseURL:", import.meta.env.VITE_API_URL);

// Request interceptor - automatically add Authorization header
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('access_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

// Response interceptor - handle token refresh
let isRefreshing = false;
let failedQueue: QueueItem[] = [];
// eslint-disable-next-line
const processQueue = (error: any, token: string | null = null): void => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else if (token) {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

api.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Handle 401 errors with token refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Queue the request if token refresh is already in progress
        return new Promise<string>((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token: string) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            }
            return api(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refresh_token = localStorage.getItem('refresh_token');
        if (!refresh_token) {
          throw new Error('No refresh token available');
        }

        const res = await axios.post<AuthResponse>(`${import.meta.env.VITE_API_URL}/auth/refresh`, { 
          refresh_token 
        });

        const newAccessToken = res.data.access_token;
        localStorage.setItem('access_token', newAccessToken);
        
        // Update default header for future requests
        if (api.defaults.headers) {
          api.defaults.headers.Authorization = `Bearer ${newAccessToken}`;
        }
        
        processQueue(null, newAccessToken);
        
        // Retry original request with new token
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }
        return api(originalRequest);
      } catch (err) {
        processQueue(err, null);
        
        // Clear all localStorage and redirect to login
        localStorage.clear();
        window.location.href = '/l';
        return Promise.reject(err);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

export default api;
