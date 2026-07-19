import axios, { type InternalAxiosRequestConfig } from 'axios';
import router from '@/router';

interface RetryableConfig extends InternalAxiosRequestConfig {
  _retryCount?: number;
}

const MAX_RETRIES = 3;

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const config = error.config as RetryableConfig | undefined;
    if (
      error.response?.status === 429 &&
      config &&
      (config._retryCount ?? 0) < MAX_RETRIES
    ) {
      config._retryCount = (config._retryCount ?? 0) + 1;
      const retryAfter = parseInt(error.response.headers['retry-after'], 10);
      const delay = (Number.isFinite(retryAfter) ? retryAfter : 2) * 1000;
      await new Promise((resolve) => setTimeout(resolve, delay));
      return api.request(config);
    }
    if (
      error.response?.status === 401 &&
      config &&
      !config.url?.includes('/admin/verify') &&
      !config.url?.includes('/admin/session')
    ) {
      router.push({ name: 'login' });
    }
    return Promise.reject(error);
  },
);

export default api;
