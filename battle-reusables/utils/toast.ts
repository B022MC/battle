import { type LucideIcon } from 'lucide-react-native';

/**
 * Toast 提示类型
 */
export type ToastType = 'success' | 'error' | 'warning' | 'info';

/**
 * Toast 选项
 */
export type ToastOptions = {
  /** 标题 */
  title: string;
  /** 描述文本 */
  description?: string;
  /** 提示类型 */
  type?: ToastType;
  /** 自定义图标 */
  icon?: LucideIcon;
  /** 显示时长(毫秒), 0表示不自动关闭 */
  duration?: number;
  /** 是否显示关闭按钮 */
  showClose?: boolean;
  /** 位置 */
  position?: 'top' | 'center' | 'bottom';
};

/**
 * 确认框选项
 */
export type ConfirmOptions = {
  /** 标题 */
  title: string;
  /** 描述文本 */
  description?: string;
  /** 提示类型 */
  type?: ToastType;
  /** 自定义图标 */
  icon?: LucideIcon;
  /** 确认按钮文本 */
  confirmText?: string;
  /** 取消按钮文本 */
  cancelText?: string;
  /** 确认按钮样式 */
  confirmVariant?: 'default' | 'destructive';
  /** 确认回调 */
  onConfirm?: () => void | Promise<void>;
  /** 取消回调 */
  onCancel?: () => void;
};

/**
 * Toast API 接口
 */
type ToastAPI = {
  show: (options: ToastOptions) => void;
  success: (title: string, description?: string) => void;
  error: (title: string, description?: string) => void;
  warning: (title: string, description?: string) => void;
  info: (title: string, description?: string) => void;
  confirm: (options: ConfirmOptions) => void;
  close: () => void;
};

let apiRef: ToastAPI | null = null;

/**
 * 绑定 Toast API
 */
export const bindToast = (api: ToastAPI) => {
  apiRef = api;
};

/**
 * 解绑 Toast API
 */
export const unbindToast = () => {
  apiRef = null;
};

/**
 * 全局 Toast 实例
 */
export const toast = {
  /**
   * 显示自定义 Toast
   */
  show(options: ToastOptions) {
    apiRef?.show(options);
  },

  /**
   * 显示成功提示
   */
  success(title: string, description?: string) {
    apiRef?.success(title, description);
  },

  /**
   * 显示错误提示
   */
  error(title: string, description?: string) {
    apiRef?.error(title, description);
  },

  /**
   * 显示警告提示
   */
  warning(title: string, description?: string) {
    apiRef?.warning(title, description);
  },

  /**
   * 显示信息提示
   */
  info(title: string, description?: string) {
    apiRef?.info(title, description);
  },

  /**
   * 显示确认框
   */
  confirm(options: ConfirmOptions) {
    apiRef?.confirm(options);
  },

  /**
   * 关闭当前 Toast
   */
  close() {
    apiRef?.close();
  },
};

/**
 * 显示 Toast 提示（兼容旧版 API）
 * @param message 提示消息
 * @param type 提示类型
 */
export function showToast(message: string, type: ToastType = 'info') {
  switch (type) {
    case 'success':
      toast.success(message);
      break;
    case 'error':
      toast.error(message);
      break;
    case 'warning':
      toast.warning(message);
      break;
    case 'info':
    default:
      toast.info(message);
      break;
  }
}
