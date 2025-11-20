/**
 * 气泡提示工具
 * 用于操作成功后的提示
 */

import { bubbleToast } from '@/components/ui/bubble-toast-container';

/**
 * 显示操作成功的气泡提示
 */
export function showSuccessBubble(title: string, description?: string) {
  bubbleToast.success(title, description);
}

/**
 * 显示操作失败的气泡提示
 */
export function showErrorBubble(title: string, description?: string) {
  bubbleToast.error(title, description);
}

/**
 * 显示警告的气泡提示
 */
export function showWarningBubble(title: string, description?: string) {
  bubbleToast.warning(title, description);
}

/**
 * 显示信息的气泡提示
 */
export function showInfoBubble(title: string, description?: string) {
  bubbleToast.info(title, description);
}

// 导出 bubbleToast 对象以便直接使用
export { bubbleToast };

