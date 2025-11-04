import React, { useState } from 'react';
import { View, TouchableOpacity, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { LoginForm } from '@/components/auth/index/login-form';
import { RegisterForm } from '@/components/auth/index/register-form';

export default function AuthScreen() {
  const [activeTab, setActiveTab] = useState('login');

  return (
    <ScrollView className="flex-1 bg-background">
      <View className="flex-1 px-4 py-24">
        {/* Header */}
        <View className="mb-8 items-center">
          <Text variant="h1" className="mb-2 font-bold text-primary">
            Hello!
          </Text>
          <Text variant="p" className="text-muted-foreground">
            欢迎来到 battle
          </Text>
        </View>

        {/* TODO Tab Selector */}
        <View className="mb-8 flex-row border-b border-border">
          <TouchableOpacity
            className={`flex-1 items-center py-3 ${activeTab === 'login' ? 'border-b-2 border-primary' : ''}`}
            onPress={() => setActiveTab('login')}>
            <Text
              className={`font-medium ${activeTab === 'login' ? 'text-primary' : 'text-muted-foreground'}`}>
              登录
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            className={`flex-1 items-center py-3 ${activeTab === 'register' ? 'border-b-2 border-primary' : ''}`}
            onPress={() => setActiveTab('register')}>
            <Text
              className={`font-medium ${activeTab === 'register' ? 'text-primary' : 'text-muted-foreground'}`}>
              注册
            </Text>
          </TouchableOpacity>
        </View>

        {/* Form Container */}
        <View className="mb-6">{activeTab === 'login' ? <LoginForm /> : <RegisterForm />}</View>
      </View>
    </ScrollView>
  );
}
