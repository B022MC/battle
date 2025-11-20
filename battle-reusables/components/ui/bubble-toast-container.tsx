import React, { useState, useCallback, useEffect } from 'react';
import { View } from 'react-native';
import { BubbleToast, BubbleToastProps, BubbleToastType } from './bubble-toast';

interface BubbleToastItem extends BubbleToastProps {
  id: string;
}

type BubbleToastAPI = {
  show: (props: Omit<BubbleToastProps, 'onClose'>) => void;
  success: (title: string, description?: string) => void;
  error: (title: string, description?: string) => void;
  warning: (title: string, description?: string) => void;
  info: (title: string, description?: string) => void;
};

let apiRef: BubbleToastAPI | null = null;

export const bindBubbleToast = (api: BubbleToastAPI) => {
  apiRef = api;
};

export const unbindBubbleToast = () => {
  apiRef = null;
};

export const bubbleToast = {
  show(props: Omit<BubbleToastProps, 'onClose'>) {
    apiRef?.show(props);
  },
  success(title: string, description?: string) {
    apiRef?.success(title, description);
  },
  error(title: string, description?: string) {
    apiRef?.error(title, description);
  },
  warning(title: string, description?: string) {
    apiRef?.warning(title, description);
  },
  info(title: string, description?: string) {
    apiRef?.info(title, description);
  },
};

export function BubbleToastContainer() {
  const [toasts, setToasts] = useState<BubbleToastItem[]>([]);

  const show = useCallback((props: Omit<BubbleToastProps, 'onClose'>) => {
    const id = `${Date.now()}-${Math.random()}`;
    setToasts((prev) => [...prev, { ...props, id }]);
  }, []);

  const success = useCallback((title: string, description?: string) => {
    show({ type: 'success', title, description, duration: 3000 });
  }, [show]);

  const error = useCallback((title: string, description?: string) => {
    show({ type: 'error', title, description, duration: 4000 });
  }, [show]);

  const warning = useCallback((title: string, description?: string) => {
    show({ type: 'warning', title, description, duration: 3500 });
  }, [show]);

  const info = useCallback((title: string, description?: string) => {
    show({ type: 'info', title, description, duration: 3000 });
  }, [show]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  useEffect(() => {
    bindBubbleToast({ show, success, error, warning, info });
    return () => unbindBubbleToast();
  }, [show, success, error, warning, info]);

  return (
    <View
      className="absolute top-0 left-0 right-0 z-[100]"
      pointerEvents="box-none"
      style={{ paddingTop: 60 }} // 留出状态栏空间
    >
      {toasts.map((toast) => (
        <BubbleToast
          key={toast.id}
          {...toast}
          onClose={() => removeToast(toast.id)}
        />
      ))}
    </View>
  );
}

