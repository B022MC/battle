import React, { useState, useEffect } from 'react';
import { Modal, View, Pressable, Platform, ScrollView } from 'react-native';
import { Text } from '@/components/ui/text';
import { Button } from '@/components/ui/button';
import { useRequest } from '@/hooks/use-request';
import { shopsHousesOptions } from '@/services/shops/houses';
import { X, ChevronDown } from 'lucide-react-native';
import { Icon } from '@/components/ui/icon';

interface SetShopAdminModalProps {
  visible: boolean;
  onClose: () => void;
  onConfirm: (houseGid: number) => void;
  userName: string;
  loading?: boolean;
}

export const SetShopAdminModal: React.FC<SetShopAdminModalProps> = ({
  visible,
  onClose,
  onConfirm,
  userName,
  loading = false,
}) => {
  const [selectedHouse, setSelectedHouse] = useState<string>('');
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const { data: houseOptions, run: loadHouses, loading: loadingHouses } = useRequest(shopsHousesOptions, { manual: true });

  useEffect(() => {
    if (visible) {
      loadHouses();
      setSelectedHouse('');
      setDropdownOpen(false);
    }
  }, [visible]);

  const handleConfirm = () => {
    if (!selectedHouse) {
      return;
    }
    const houseGid = Number(selectedHouse);
    if (isNaN(houseGid) || houseGid <= 0) {
      return;
    }
    onConfirm(houseGid);
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
    >
      <Pressable
        className="flex-1 bg-black/50 justify-center items-center"
        onPress={() => {
          setDropdownOpen(false);
          onClose();
        }}
      >
        <Pressable
          className="bg-card rounded-lg p-6 w-[90%] max-w-md shadow-lg"
          onPress={(e) => {
            e.stopPropagation();
            setDropdownOpen(false);
          }}
        >
          {/* 标题栏 */}
          <View className="flex-row justify-between items-center mb-4">
            <Text className="text-lg font-semibold">设置店铺管理员</Text>
            <Pressable onPress={onClose} className="p-1">
              <Icon as={X} className="text-muted-foreground size-5" />
            </Pressable>
          </View>

          {/* 用户信息 */}
          <View className="mb-4">
            <Text className="text-muted-foreground text-sm mb-1">用户</Text>
            <Text className="font-medium">{userName}</Text>
          </View>

          {/* 店铺选择 */}
          <View className="mb-6">
            <Text className="text-muted-foreground text-sm mb-2">
              选择店铺 {houseOptions && `(共 ${houseOptions.length} 个)`}
            </Text>
            {loadingHouses ? (
              <View className="p-4 border border-border rounded-md">
                <Text className="text-muted-foreground text-center text-sm">加载中...</Text>
              </View>
            ) : !houseOptions || houseOptions.length === 0 ? (
              <View className="p-4 border border-border rounded-md">
                <Text className="text-muted-foreground text-center text-sm">暂无可用店铺</Text>
              </View>
            ) : (
              <View>
                {/* 自定义下拉框触发器 */}
                <Pressable
                  onPress={() => setDropdownOpen(!dropdownOpen)}
                  className="flex-row items-center justify-between border border-border rounded-md px-3 py-2 bg-background"
                >
                  <Text className={selectedHouse ? 'text-foreground' : 'text-muted-foreground'}>
                    {selectedHouse ? `店铺 ${selectedHouse}` : '请选择店铺'}
                  </Text>
                  <Icon as={ChevronDown} className="text-muted-foreground size-4" />
                </Pressable>

                {/* 下拉选项列表 */}
                {dropdownOpen && (
                  <View
                    className="mt-2 bg-popover border border-border rounded-md shadow-lg"
                    style={{ maxHeight: houseOptions.length > 3 ? 150 : undefined }}
                  >
                    <ScrollView
                      nestedScrollEnabled
                      showsVerticalScrollIndicator={houseOptions.length > 3}
                    >
                      {houseOptions.map((houseGid) => (
                        <Pressable
                          key={houseGid}
                          onPress={() => {
                            setSelectedHouse(String(houseGid));
                            setDropdownOpen(false);
                          }}
                          className="px-3 py-2 hover:bg-accent active:bg-accent border-b border-border last:border-b-0"
                        >
                          <Text className="text-foreground">店铺 {houseGid}</Text>
                        </Pressable>
                      ))}
                    </ScrollView>
                  </View>
                )}
              </View>
            )}
          </View>

          {/* 操作按钮 */}
          <View className="flex-row gap-3">
            <Button
              variant="outline"
              onPress={onClose}
              disabled={loading}
              className="flex-1"
            >
              <Text>取消</Text>
            </Button>
            <Button
              onPress={handleConfirm}
              disabled={!selectedHouse || loading}
              className="flex-1"
            >
              <Text>{loading ? '设置中...' : '确认'}</Text>
            </Button>
          </View>
        </Pressable>
      </Pressable>
    </Modal>
  );
};

