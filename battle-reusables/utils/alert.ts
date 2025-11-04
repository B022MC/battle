import { type LucideIcon } from 'lucide-react-native';

export type AlertOptions = {
  title: string;
  description?: string;
  icon?: LucideIcon;
  duration?: number; // 毫秒
  showClose?: boolean;
  confirmText?: string;
  cancelText?: string;
  onConfirm?: () => void;
  onCancel?: () => void;
};

type AlertAPI = { show: (options: AlertOptions) => void };

let apiRef: AlertAPI | null = null;

export const bindAlert = (api: AlertAPI) => {
  apiRef = api;
};

export const unbindAlert = () => {
  apiRef = null;
};

export const alert = {
  show(options: AlertOptions) {
    apiRef?.show(options);
  },
};
