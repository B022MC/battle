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
import { isWeb } from '@/utils/platform';
import { useRequest } from '@/hooks/use-request';
import { useAuthStore } from '@/hooks/use-auth-store';
import { platformsList } from '@/services/platforms';
import { loginUsername } from '@/services/login';
import { basicUserMeRoles, basicUserMePerms } from '@/services/basic/user';
import z from 'zod';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ErrorMessage } from '@/components/shared/error-message';
import { Button } from '@/components/ui/button';
import { router } from 'expo-router';
import { encrypt } from '@/utils/rsa';

const loginFormSchema = z.lazy(() =>
  z.object({
    username: z.string().min(1, { message: '请输入用户名' }),
    password: z.string().min(1, { message: '请输入密码' }),
    platform: z.string().min(1, { message: '请选择登录平台' }),
  })
);

type LoginFormValues = z.infer<typeof loginFormSchema>;

export const LoginForm = () => {
  const { data: platforms } = useRequest(platformsList);

  const [showPassword, setShowPassword] = useState(false);

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
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      username: '',
      password: '',
      platform: '',
    },
    mode: 'all',
  });

  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  const { run: runUserLogin, loading } = useRequest(loginUsername, {
    manual: true,
    onSuccess: async () => {
      // 登录成功后立即获取权限和角色
      const { updateAuth } = useAuthStore.getState();
      try {
        const [rolesRes, permsRes] = await Promise.all([
          basicUserMeRoles(),
          basicUserMePerms(),
        ]);
        updateAuth({
          roles: rolesRes?.data?.role_ids,
          perms: permsRes?.data?.perms,
        });
      } catch (error) {
        console.error('获取权限失败:', error);
      }
      router.push('/(tabs)');
    },
  });

  const onSubmit = (values: LoginFormValues) => {
    const { username, password, platform } = values;
    const { updateAuth } = useAuthStore.getState();
    updateAuth({ platform });
    runUserLogin({ username, password: encrypt(password) });
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
                onPress={() => setShowPassword(!showPassword)}
                className="text-muted-foreground"
              />
            </View>
          </View>
          <ErrorMessage name="password" errors={errors} />
        </View>

        {/* Login Button */}
        <Button
          className="mb-4 w-full rounded-full"
          onPress={handleSubmit(onSubmit)}
          disabled={!isValid || loading}>
          {loading && (
            <View className="pointer-events-none animate-spin">
              <Icon as={Loader2} className="text-primary-foreground" />
            </View>
          )}
          <Text className="font-medium">登录</Text>
        </Button>
      </View>
    </KeyboardAvoidingView>
  );
};
