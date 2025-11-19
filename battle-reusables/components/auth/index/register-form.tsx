import React, { useRef, useState } from 'react';
import { KeyboardAvoidingView, Platform, View, ActivityIndicator } from 'react-native';
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
import { Eye, EyeOff, Loader2, CheckCircle2, XCircle, AlertCircle } from 'lucide-react-native';
import { TriggerRef } from '@rn-primitives/select';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { encrypt } from '@/utils/rsa';
import { md5Upper } from '@/utils/md5';
import { isWeb } from '@/utils/platform';
import { useRequest } from '@/hooks/use-request';
import { useAuthStore } from '@/hooks/use-auth-store';
import { platformsList } from '@/services/platforms';
import { loginRegister } from '@/services/login';
import { gameAccountVerify } from '@/services/game/account';
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
      // Game account fields
      gameAccountMode: z.enum(['account', 'mobile'], { message: '请选择游戏账号类型' }),
      gameAccount: z.string().min(1, { message: '请输入游戏账号' }),
      gamePassword: z.string().min(1, { message: '请输入游戏密码' }),
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
  const [showGamePassword, setShowGamePassword] = useState(false);
  const [gameAccountVerified, setGameAccountVerified] = useState<boolean | null>(null);

  const ref = useRef<TriggerRef>(null);
  const gameModeRef = useRef<TriggerRef>(null);
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
    watch,
    formState: { errors, isValid },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerFormSchema),
    defaultValues: {
      username: '',
      password: '',
      platform: '',
      nickname: '',
      confirmPassword: '',
      gameAccountMode: 'account',
      gameAccount: '',
      gamePassword: '',
    },
    mode: 'all',
  });

  const gameAccountMode = watch('gameAccountMode');
  const gameAccount = watch('gameAccount');
  const gamePassword = watch('gamePassword');

  function onTouchStart() {
    isWeb && ref.current?.open();
  }

  function onGameModeTouchStart() {
    isWeb && gameModeRef.current?.open();
  }

  // Verify game account (manual trigger only, no auto-verify)
  const { run: runVerifyAccount, loading: verifying } = useRequest(gameAccountVerify, {
    manual: true,
    onSuccess: () => {
      setGameAccountVerified(true);
    },
    onError: () => {
      setGameAccountVerified(false);
    },
  });

  const { run: runUserRegister, loading: registering } = useRequest(loginRegister, {
    manual: true,
    onSuccess: () => router.replace('/(tabs)'),
  });

  const onSubmit = async (values: RegisterFormValues) => {
    const { username, password, platform, nickname, gameAccountMode, gameAccount, gamePassword } =
      values;

    // 先验证游戏账号
    if (gameAccount && gamePassword) {
      setGameAccountVerified(null); // 重置验证状态
      const pwdMD5 = md5Upper(gamePassword);

      try {
        // 调用验证接口
        await runVerifyAccount({
          mode: gameAccountMode,
          account: gameAccount,
          pwd_md5: pwdMD5,
        });

        // 验证成功，继续注册
        const { updateAuth } = useAuthStore.getState();
        updateAuth({ platform });
        runUserRegister({
          username,
          password: encrypt(password),
          nick_name: nickname,
          game_account_mode: gameAccountMode,
          game_account: gameAccount,
          game_password: pwdMD5,
        });
      } catch (error) {
        // 验证失败，不继续注册
        setGameAccountVerified(false);
        return;
      }
    } else {
      // 没有填写游戏账号，直接注册
      const { updateAuth } = useAuthStore.getState();
      updateAuth({ platform });
      runUserRegister({ username, password: encrypt(password), nick_name: nickname });
    }
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

        {/* Divider */}
        <View className="my-2 flex-row items-center">
          <View className="h-px flex-1 bg-border" />
          <Text className="mx-4 text-sm text-muted-foreground">游戏账号绑定</Text>
          <View className="h-px flex-1 bg-border" />
        </View>

        {/* Game Account Mode Selection */}
        <View className="gap-1">
          <Controller
            control={control}
            name="gameAccountMode"
            render={({ field: { onChange, value } }) => (
              <Select 
                onValueChange={(option) => onChange(option?.value)} 
                value={{ value, label: value === 'mobile' ? '手机号' : '游戏账号' }}
              >
                <SelectTrigger
                  ref={gameModeRef}
                  className="w-full"
                  onTouchStart={onGameModeTouchStart}>
                  <SelectValue placeholder="请选择账号类型" />
                </SelectTrigger>
                <SelectContent insets={contentInsets} className="w-full">
                  <SelectGroup>
                    <SelectLabel>账号类型</SelectLabel>
                    <SelectItem label="游戏账号" value="account">
                      游戏账号
                    </SelectItem>
                    <SelectItem label="手机号" value="mobile">
                      手机号
                    </SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            )}
          />
          <ErrorMessage name="gameAccountMode" errors={errors} />
        </View>

        {/* Game Account Input */}
        <View className="gap-1">
          <Controller
            control={control}
            name="gameAccount"
            render={({ field: { onChange, value } }) => (
              <Input
                value={value}
                onChangeText={onChange}
                placeholder={gameAccountMode === 'mobile' ? '请输入手机号' : '请输入游戏账号'}
              />
            )}
          />
          <ErrorMessage name="gameAccount" errors={errors} />
          {gameAccountVerified === false && (
            <View className="flex-row items-center gap-1">
              <Icon as={AlertCircle} size={14} className="text-red-500" />
              <Text className="text-xs text-red-500">游戏账号验证失败，请检查账号和密码</Text>
            </View>
          )}
        </View>

        {/* Game Password Input */}
        <View className="gap-1">
          <View className="relative flex-row items-center">
            <Controller
              control={control}
              name="gamePassword"
              render={({ field: { onChange, value } }) => (
                <Input
                  textContentType="password"
                  className="flex-1 pr-10"
                  value={value}
                  onChangeText={onChange}
                  placeholder="请输入游戏密码"
                  secureTextEntry={!showGamePassword}
                />
              )}
            />
            <View className="absolute right-3">
              <Icon
                as={showGamePassword ? Eye : EyeOff}
                size={20}
                className="text-muted-foreground"
                onPress={() => setShowGamePassword(!showGamePassword)}
              />
            </View>
          </View>
          <ErrorMessage name="gamePassword" errors={errors} />
        </View>

        {/* Register Button */}
        <Button
          className="mb-4 mt-4 w-full rounded-full"
          onPress={handleSubmit(onSubmit)}
          disabled={!isValid || registering || verifying}>
          {(registering || verifying) && (
            <View className="pointer-events-none animate-spin">
              <Icon as={Loader2} className="text-primary-foreground" />
            </View>
          )}
          <Text className="font-medium">
            {verifying ? '验证中...' : registering ? '注册中...' : '注册'}
          </Text>
        </Button>
      </View>
    </KeyboardAvoidingView>
  );
};
