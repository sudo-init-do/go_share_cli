import axios from 'axios';
import { PageData } from '../types';

const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8081';

const api = axios.create({
  baseURL: API_BASE,
  withCredentials: true,
});

export const authService = {
  async login(password: string): Promise<boolean> {
    try {
      const response = await api.post('/login', 
        new URLSearchParams({
          password,
          redirect: '/'
        }),
        {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          maxRedirects: 0,
          validateStatus: (status) => status < 400,
        }
      );
      return response.status === 303 || response.status === 200;
    } catch (error) {
      return false;
    }
  },

  async checkAuth(): Promise<boolean> {
    try {
      const response = await api.get('/api/auth/check');
      return response.status === 200;
    } catch (error) {
      return false;
    }
  },

  async logout(): Promise<void> {
    try {
      await api.post('/logout');
    } catch (error) {
      console.error('Logout error:', error);
    }
  }
};

export const fileService = {
  async getFiles(path: string = '/'): Promise<PageData> {
    try {
      const response = await api.get(`/api/files`, {
        params: { path }
      });
      return response.data;
    } catch (error) {
      throw new Error('Failed to fetch files');
    }
  },

  async uploadFiles(files: FileList, directory: string = '/'): Promise<void> {
    const formData = new FormData();
    formData.append('directory', directory);
    
    for (let i = 0; i < files.length; i++) {
      formData.append('files', files[i]);
    }

    await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  getDownloadUrl(path: string): string {
    return `${API_BASE}${path}?download=1`;
  },

  getPreviewUrl(path: string): string {
    return `${API_BASE}${path}`;
  }
};

export default api;
