import React, { useRef, useState } from 'react';
import { KeyboardAvoidingView, Platform, View } from 'react-native';
import { Text } from '@/components/ui/text';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
  SelectGroup,
  SelectLabel,
} from '@/components/ui/select';
import { Icon } from '@/components/ui/icon';
import { Eye, EyeOff, Loader2 } from 'lucide-react-native';
import { TriggerRef } from '@rn-primitives/select';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { encrypt } from '@/utils/rsa';
import { isWeb } from '@/utils/platform';
import { useRequest } from '@/hooks/use-request';
import { useAuthStore } from '@/hooks/use-auth-store';
import { platformsList } from '@/services/platforms';
import { loginRegister } from '@/services/login';
import z from 'zod';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@/components/ui/button';
import { ErrorMessage } from '@/components/shared/error-message';
import { router } from 'expo-router';

const registerFormSchema = z.lazy(() =>
  z
    .object({
      username: z.string().min(1, { message: '请输入用户名' }),
      password: z.string().min(1, { message: '请输入密码' }),
      platform: z.string().min(1, { message: '请选择登录平台' }),
      nickname: z.string().optional(),
      confirmPassword: z.string().min(1, { message: '请确认密码' }),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: '两次输入的密码不一致',
      path: ['confirmPassword'],
    })
);

type RegisterFormValues = z.infer<typeof registerFormSchema>;

export const RegisterForm = () => {
  const { data: platforms } = useRequest(platformsList);

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const ref = useRef<TriggerRef>(null);
  const insets = useSafeAreaInsets();
  const contentInsets = {
    top: insets.top,
    bottom: Platform.select({ ios: insets.bottom, android: insets.bottom + 24 }),
    left: 12,
    right: 12,
  };

  const {
    control,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerFormSchema),
    defaultValues: {
      username: '',
      password: '',
      platform: '',
      nickname: '',
      confirmPassword: '',
    },
    mode: 'all',
  });

  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  const { run: runUserRegister, loading } = useRequest(loginRegister, {
    manual: true,
    onSuccess: () => router.replace('/'),
  });

  const onSubmit = (values: RegisterFormValues) => {
    const { username, password, platform, nickname } = values;
    const { updateAuth } = useAuthStore.getState();
    updateAuth({ platform });
    runUserRegister({ username, password: encrypt(password), nick_name: nickname });
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={{ flex: 1 }}>
      <View className="gap-6">
        {/* 平台选择 */}
        <View className="gap-1">
          <Controller
            control={control}
            name="platform"
            render={({ field: { onChange } }) => (
              <Select onValueChange={(option) => onChange(option?.value)}>
                <SelectTrigger ref={ref} className="w-full" onTouchStart={onTouchStart}>
                  <SelectValue placeholder="请选择平台" />
                </SelectTrigger>
                <SelectContent insets={contentInsets} className="w-full">
                  <SelectGroup>
                    <SelectLabel>平台</SelectLabel>
                    {platforms?.map(({ name, platform }) => (
                      <SelectItem label={name!} key={platform} value={platform!}>
                        {name}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            )}
          />
          <ErrorMessage name="platform" errors={errors} />
        </View>

        {/* 用户名输入 */}
        <View className="gap-1">
          <Controller
            control={control}
            name="username"
            render={({ field: { onChange, value } }) => (
              <Input value={value} onChangeText={onChange} placeholder="请输入用户名" />
            )}
          />
          <ErrorMessage name="username" errors={errors} />
        </View>

        {/* 密码输入 */}
        <View className="gap-1">
          <View className="relative flex-row items-center">
            <Controller
              control={control}
              name="password"
              render={({ field: { onChange, value } }) => (
                <Input
                  textContentType="password"
                  className="flex-1 pr-10"
                  value={value}
                  onChangeText={onChange}
                  placeholder="请输入密码"
                  secureTextEntry={!showPassword}
                />
              )}
            />
            <View className="absolute right-3">
              <Icon
                as={showPassword ? Eye : EyeOff}
                size={20}
                className="text-muted-foreground"
                onPress={() => setShowPassword(!showPassword)}
              />
            </View>
          </View>
          <ErrorMessage name="password" errors={errors} />
        </View>

        {/* 确认密码输入 */}
        <View className="gap-1">
          <View className="relative flex-row items-center">
            <Controller
              control={control}
              name="confirmPassword"
              render={({ field: { onChange, value } }) => (
                <Input
                  textContentType="password"
                  className="flex-1 pr-10"
                  value={value}
                  onChangeText={onChange}
                  placeholder="请再次输入密码"
                  secureTextEntry={!showConfirmPassword}
                />
              )}
            />
            <View className="absolute right-3">
              <Icon
                as={showConfirmPassword ? Eye : EyeOff}
                size={20}
                className="text-muted-foreground"
                onPress={() => setShowConfirmPassword(!showConfirmPassword)}
              />
            </View>
          </View>
          <ErrorMessage name="confirmPassword" errors={errors} />
        </View>

        {/* 昵称输入 */}
        <View className="gap-1">
          <Controller
            control={control}
            name="nickname"
            render={({ field: { onChange, value } }) => (
              <Input value={value} onChangeText={onChange} placeholder="请输入昵称" />
            )}
          />
          <ErrorMessage name="nickname" errors={errors} />
        </View>

        {/* Register Button */}
        <Button
          className="mb-4 w-full rounded-full"
          onPress={handleSubmit(onSubmit)}
          disabled={!isValid || loading}>
          {loading && (
            <View className="pointer-events-none animate-spin">
              <Icon as={Loader2} className="text-primary-foreground" />
            </View>
          )}
          <Text className="font-medium">注册</Text>
        </Button>
      </View>
    </KeyboardAvoidingView>
  );
};
