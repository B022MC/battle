import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { Platform } from 'react-native';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function omitResponderProps<T extends Record<string, any>>(props: T): T {
  if (!props || Platform.OS !== 'web') return props;
  const filtered: Record<string, any> = {};
  for (const key in props) {
    if (
      key === 'onStartShouldSetResponder' ||
      key === 'onMoveShouldSetResponder' ||
      key === 'onStartShouldSetResponderCapture' ||
      key === 'onMoveShouldSetResponderCapture' ||
      key.startsWith('onResponder')
    ) {
      continue;
    }
    filtered[key] = (props as any)[key];
  }
  return filtered as T;
}
