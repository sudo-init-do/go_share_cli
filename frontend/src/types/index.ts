export interface FileItem {
  name: string;
  path: string;
  size: number;
  isDir: boolean;
  modTime: string;
  downloadCount: number;
}

export interface PageData {
  title: string;
  currentPath: string;
  parentPath: string;
  hasParent: boolean;
  files: FileItem[];
  serverURL: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  sessionToken?: string;
}

export interface UploadProgress {
  fileName: string;
  progress: number;
  status: 'uploading' | 'completed' | 'error';
}
