// hooks/use-auth-store.ts
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import AsyncStorage from '@react-native-async-storage/async-storage';

type AuthState = {
  isAuthenticated: boolean;
  isLoading: boolean;
  accessToken?: string;
  refreshToken?: string;
  expiresAt?: number; // 绝对过期时间
  platform?: string;
  role?: string; // 用户角色：super_admin, store_admin, user
  roles?: number[];
  perms?: string[];
  user?: {
    id?: number;
    username?: string;
    nick_name?: string;
    avatar?: string;
  };
};

type AuthActions = {
  hydrate: () => void; // 持久化恢复后收尾
  updateAuth: (p: {
    accessToken?: string;
    refreshToken?: string;
    expiresIn?: number;
    platform?: string;
    role?: string;
    roles?: number[];
    perms?: string[];
    user?: {
      id?: number;
      username?: string;
      nick_name?: string;
      avatar?: string;
    };
  }) => void;
  clearAuth: () => void;
  refreshFromStorage: () => Promise<void>; // 可选：按需手动刷新
};

const getIsAuthenticated = (state: Partial<AuthState>) => {
  const { accessToken, expiresAt } = state;
  return accessToken && expiresAt && expiresAt > Date.now() ? true : false;
};


 export const useAuthStore = create<AuthState & AuthActions>()(
   persist(
     (set, get) => ({
       isAuthenticated: false,
       isLoading: true,

       hydrate: () => {
         set({ isLoading: false, isAuthenticated: getIsAuthenticated(get()) });
       },

      updateAuth: ({ accessToken, refreshToken, expiresIn, platform, role, roles, perms, user }) => {
         set((state) => {
           const newState = {
             accessToken: accessToken ?? state.accessToken,
             refreshToken: refreshToken ?? state.refreshToken,
             expiresAt: expiresIn ? Date.now() + expiresIn * 1000 : state.expiresAt,
             platform: platform ?? state.platform,
             role: role ?? state.role,
             roles: roles ?? state.roles,
             perms: perms ?? state.perms,
            user: user ?? state.user,
           } as AuthState;
           return {
             ...newState,
             isLoading: false,
             isAuthenticated: getIsAuthenticated(newState),
           };
         });
       },

       clearAuth: () => {
         set({
           accessToken: undefined,
           refreshToken: undefined,
           expiresAt: undefined,
           platform: undefined,
           role: undefined,
           roles: undefined,
           perms: undefined,
           isAuthenticated: false,
           isLoading: false,
         });
       },

       // 可选：如果你有复杂刷新逻辑，也可以在这里做
       refreshFromStorage: async () => {
         const token = get().accessToken;
         set({ isAuthenticated: !!token, isLoading: false });
       },
     }),
     {
       name: 'auth',
       storage: createJSONStorage(() => AsyncStorage),
       // 只持久化 session 字段（可选）
       partialize: (state) => ({
         accessToken: state.accessToken,
         refreshToken: state.refreshToken,
         expiresAt: state.expiresAt,
         platform: state.platform,
         role: state.role,
         roles: state.roles,
         perms: state.perms,
          user: state.user,
       }),
       // rehydrate 完成后，统一收尾
       onRehydrateStorage: () => (state) => {
         state?.hydrate?.();
       },
     }
   )
 );


