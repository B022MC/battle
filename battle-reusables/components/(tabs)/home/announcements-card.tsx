import React from 'react';
import { View, TouchableOpacity, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Card } from '@/components/ui/card';
import { Megaphone, ChevronRight } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';

interface Announcement {
  id: number;
  title: string;
  content: string;
  priority: 'normal' | 'important' | 'urgent';
  created_at: string;
}

interface AnnouncementsCardProps {
  announcements: Announcement[];
  onViewAll?: () => void;
}

export function AnnouncementsCard({ announcements, onViewAll }: AnnouncementsCardProps) {
  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'urgent':
        return 'text-red-600';
      case 'important':
        return 'text-orange-600';
      default:
        return 'text-gray-700';
    }
  };

  const getPriorityBadge = (priority: string) => {
    switch (priority) {
      case 'urgent':
        return { bg: 'bg-red-100', text: 'text-red-700', label: '紧急' };
      case 'important':
        return { bg: 'bg-orange-100', text: 'text-orange-700', label: '重要' };
      default:
        return null;
    }
  };

  return (
    <Card className="p-4">
      {/* 标题栏 */}
      <View className="flex-row items-center justify-between mb-4">
        <View className="flex-row items-center gap-2">
          <Icon as={Megaphone} className="size-5 text-blue-600" />
          <Text className="text-lg font-semibold">系统公告</Text>
        </View>
        {onViewAll && announcements.length > 3 && (
          <TouchableOpacity
            onPress={onViewAll}
            className="flex-row items-center gap-1"
          >
            <Text className="text-sm text-blue-600">查看全部</Text>
            <Icon as={ChevronRight} className="size-4 text-blue-600" />
          </TouchableOpacity>
        )}
      </View>

      {/* 公告列表 */}
      {announcements.length === 0 ? (
        <View className="py-8 items-center">
          <Text className="text-sm text-muted-foreground">暂无公告</Text>
        </View>
      ) : (
        <View className="gap-3">
          {announcements.slice(0, 3).map((announcement) => {
            const badge = getPriorityBadge(announcement.priority);
            return (
              <TouchableOpacity
                key={announcement.id}
                className="p-3 bg-secondary rounded-lg active:bg-secondary/80"
              >
                <View className="flex-row items-start gap-2">
                  {badge && (
                    <View className={`px-2 py-0.5 rounded ${badge.bg}`}>
                      <Text className={`text-xs font-medium ${badge.text}`}>
                        {badge.label}
                      </Text>
                    </View>
                  )}
                  <View className="flex-1">
                    <Text 
                      className={`text-sm font-medium ${getPriorityColor(announcement.priority)}`}
                      numberOfLines={1}
                    >
                      {announcement.title}
                    </Text>
                    <Text 
                      className="text-xs text-muted-foreground mt-1"
                      numberOfLines={2}
                    >
                      {announcement.content}
                    </Text>
                    <Text className="text-xs text-muted-foreground mt-1">
                      {new Date(announcement.created_at).toLocaleDateString('zh-CN', {
                        month: '2-digit',
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </Text>
                  </View>
                </View>
              </TouchableOpacity>
            );
          })}
        </View>
      )}
    </Card>
  );
}
