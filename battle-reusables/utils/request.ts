import { isWeb } from './platform';
import { Platform } from 'react-native';
import { alert } from './alert';
import { useAuthStore } from '@/hooks/use-auth-store';

type RequestOptions = RequestInit & {
  url?: string; // 支持传入 url 字段
  params?: Record<string, any>;
  data?: Record<string, any>;
  showError?: boolean; // 是否展示全局错误弹窗，默认 false
};

export type ResponseStructure<T> = {
  code?: number;
  data?: T;
  msg?: string;
};

let lastAuthErrorAt = 0;

// 支持两种调用方式：
// 1. request(path, options) - 传统方式
// 2. request({ url, method, data, ... }) - 对象方式
export const request = async <T>(
  pathOrOptions: string | RequestOptions,
  optionsParam?: RequestOptions
): Promise<ResponseStructure<T>> => {
  // 判断调用方式并规范化参数
  let path: string;
  let options: RequestOptions;

  if (typeof pathOrOptions === 'string') {
    // 传统调用方式: request(path, options)
    path = pathOrOptions;
    options = optionsParam || {};
  } else {
    // 对象调用方式: request({ url, method, data, ... })
    const { url, ...rest } = pathOrOptions;
    path = url || '';
    options = rest;
  }

  const { headers, params, data, showError = true, ...rest } = options;

  const hostFromEnv = process.env.EXPO_PUBLIC_API_HOST;
  const nativeApiUrl = hostFromEnv
    ? `http://${hostFromEnv}:8000`
    : process.env.EXPO_PUBLIC_DEV_API_URL ||
      (Platform.OS === 'android' ? 'https://10.0.2.2:8000' : 'https://127.0.0.1:8000');

  // Web 环境下，如果有自定义 host 则使用，否则使用开发环境的 API URL
  const apiUrl = isWeb
    ? (hostFromEnv ? `http://${hostFromEnv}:8000` : process.env.EXPO_PUBLIC_DEV_API_URL || 'http://127.0.0.1:8000')
    : nativeApiUrl;

  const query = params ? `?${new URLSearchParams(params).toString()}` : '';
  const url = `${apiUrl}${path}${query}`;

  const config: RequestOptions = { ...rest };

  const { platform, accessToken, updateAuth } = useAuthStore.getState();

  config.headers = {
    Accept: 'application/json',
    ...(data && { 'Content-Type': 'application/json' }),
    ...(platform && { Platform: platform }),
    ...(accessToken && { Authorization: `${accessToken}` }),
    ...headers,
  };

  if (data) config.body = JSON.stringify(data);

  const response = await fetch(url, config);
  let res: ResponseStructure<T>;

  try {
    res = await response.json();
  } catch (error) {
    const msg = `Failed to parse response as JSON: ${error}`;
    if (showError) alert.show({ title: '请求失败', description: msg });
    throw new Error(msg);
  }

  if (!response.ok || res.code !== 0) {
    const msg = res.msg ?? response.statusText ?? `Request Failed with status ${response.status}`;
    // 401 统一处理：清除会话，触发布局跳转到登录页（避免在登录页重复弹窗）
    if (response.status === 401 || res.code === 401) {
      const now = Date.now();
      const { clearAuth } = useAuthStore.getState();
      clearAuth();
      // 防抖：1.5s 内不重复弹错误
      if (showError && now - lastAuthErrorAt > 1500) {
        alert.show({ title: '登录已过期', description: '请重新登录' });
        lastAuthErrorAt = now;
      }
    } else if (showError) {
      alert.show({ title: '请求失败', description: msg });
    }
    throw new Error(msg);
  }

  // 登录/注册响应附带令牌、平台、角色、权限、用户信息时，更新到会话
  const d: any = res.data;
  if (d && (d.access_token || d.refresh_token)) {
    updateAuth({
      accessToken: d.access_token,
      refreshToken: d.refresh_token,
      expiresIn: d.expires_in,
      platform: d.platform,
      roles: d.roles,
      perms: d.perms,
      user: d.user,
    });
  }

  return res;
};

export const get = <T>(path: string, params?: Record<string, any>) =>
  request<T>(path, { method: 'GET', params });

export const post = <T>(path: string, data?: Record<string, any>) =>
  request<T>(path, { method: 'POST', data });

export const patch = <T>(path: string, data?: Record<string, any>) =>
  request<T>(path, { method: 'PATCH', data });

export const put = <T>(path: string, data?: Record<string, any>) =>
  request<T>(path, { method: 'PUT', data });

export const del = <T>(path: string) => request<T>(path, { method: 'DELETE' });
