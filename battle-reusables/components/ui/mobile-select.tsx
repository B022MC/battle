import React, { useState } from 'react';
import { View, Platform, Pressable, Modal, TouchableOpacity, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Icon } from '@/components/ui/icon';
import { ChevronDown, Check } from 'lucide-react-native';
import { cn } from '@/lib/utils';

export interface MobileSelectOption {
  label: string;
  value: string;
}

export interface MobileSelectProps {
  value?: string;
  placeholder?: string;
  options: MobileSelectOption[];
  onValueChange?: (value: string) => void;
  className?: string;
  disabled?: boolean;
}

/**
 * 移动端友好的 Select 组件
 * 在移动端浏览器使用原生模态框，在桌面端使用常规下拉
 */
export function MobileSelect({
  value,
  placeholder = '请选择',
  options,
  onValueChange,
  className,
  disabled = false,
}: MobileSelectProps) {
  const [isOpen, setIsOpen] = useState(false);

  const selectedOption = options.find(opt => opt.value === value);
  const displayText = selectedOption?.label || placeholder;

  const handleSelect = (optionValue: string) => {
    onValueChange?.(optionValue);
    setIsOpen(false);
  };

  // 检测是否为移动端（包括平板）
  const isMobile = Platform.OS === 'web' && typeof window !== 'undefined' && 
    ('ontouchstart' in window || navigator.maxTouchPoints > 0);

  if (Platform.OS === 'web' && !isMobile) {
    // 桌面端使用原生 select
    return (
      <select
        value={value || ''}
        onChange={(e) => onValueChange?.(e.target.value)}
        disabled={disabled}
        className={cn(
          'flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm shadow-black/5',
          'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring',
          'disabled:cursor-not-allowed disabled:opacity-50',
          !value && 'text-muted-foreground',
          className
        )}
      >
        <option value="" disabled>
          {placeholder}
        </option>
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    );
  }

  // 移动端使用自定义模态框
  return (
    <>
      <Pressable
        onPress={() => !disabled && setIsOpen(true)}
        className={cn(
          'flex h-10 w-full flex-row items-center justify-between rounded-md border border-input bg-background px-3 py-2 shadow-sm shadow-black/5',
          disabled && 'opacity-50',
          className
        )}
        disabled={disabled}
      >
        <Text
          className={cn(
            'text-sm',
            !value && 'text-muted-foreground'
          )}
        >
          {displayText}
        </Text>
        <Icon as={ChevronDown} className="text-muted-foreground size-4" />
      </Pressable>

      <Modal
        visible={isOpen}
        transparent
        animationType="slide"
        onRequestClose={() => setIsOpen(false)}
      >
        <Pressable
          className="flex-1 bg-black/50"
          onPress={() => setIsOpen(false)}
        >
          <View className="flex-1" />
          <View className="bg-background rounded-t-3xl border-t border-border max-h-[70%]">
            {/* 标题栏 */}
            <View className="flex-row items-center justify-between px-4 py-3 border-b border-border">
              <Text className="text-base font-semibold">{placeholder}</Text>
              <TouchableOpacity onPress={() => setIsOpen(false)}>
                <Text className="text-primary text-base font-medium">完成</Text>
              </TouchableOpacity>
            </View>

            {/* 选项列表 */}
            <ScrollView className="flex-1">
              {options.map((option) => {
                const isSelected = option.value === value;
                return (
                  <TouchableOpacity
                    key={option.value}
                    onPress={() => handleSelect(option.value)}
                    className={cn(
                      'flex-row items-center justify-between px-4 py-3 border-b border-border/50',
                      isSelected && 'bg-accent'
                    )}
                  >
                    <Text
                      className={cn(
                        'text-base',
                        isSelected && 'font-semibold text-primary'
                      )}
                    >
                      {option.label}
                    </Text>
                    {isSelected && (
                      <Icon as={Check} className="text-primary size-5" />
                    )}
                  </TouchableOpacity>
                );
              })}
            </ScrollView>
          </View>
        </Pressable>
      </Modal>
    </>
  );
}
