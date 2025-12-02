import { useEffect, useState, useMemo } from 'react';
import { useRequest } from '@/hooks/use-request';
import { usePermission } from '@/hooks/use-permission';
import { shopsAdminsMe } from '@/services/shops/admins';
import { shopsHousesOptions } from '@/services/shops/houses';
import type { MobileSelectOption } from '@/components/ui/mobile-select';

/**
 * 店铺选择Hook
 * - 店铺管理员：自动加载并使用其管理的店铺
 * - 超级管理员：提供下拉选择所有店铺
 */
export function useHouseSelector() {
  const { isSuperAdmin, isStoreAdmin } = usePermission();
  const [houseGid, setHouseGid] = useState('');
  const [isReady, setIsReady] = useState(false);

  // 获取店铺管理员的店铺信息 - 手动加载避免无限循环
  const { data: adminInfo, loading: loadingAdmin, run: loadAdminInfo } = useRequest(
    shopsAdminsMe,
    {
      manual: true,
      onSuccess: (data) => {
        if (data?.house_gid) {
          setHouseGid(String(data.house_gid));
          setIsReady(true);
        }
      },
    }
  );

  // 获取超级管理员的店铺选项 - 手动加载避免无限循环
  const { data: houseOptions, loading: loadingOptions, run: loadHouseOptions } = useRequest(
    shopsHousesOptions,
    {
      manual: true,
      onSuccess: () => {
        setIsReady(true);
      },
    }
  );

  // 当角色确定后加载对应数据 - 只加载一次
  useEffect(() => {
    if (isStoreAdmin && !adminInfo && !loadingAdmin) {
      loadAdminInfo();
    } 
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isStoreAdmin, adminInfo, loadingAdmin]);

  useEffect(() => {
    if (isSuperAdmin && !houseOptions && !loadingOptions) {
      loadHouseOptions();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isSuperAdmin, houseOptions, loadingOptions]);

  // 更新店铺ID和准备状态
  useEffect(() => {
    if (isStoreAdmin && adminInfo?.house_gid) {
      setHouseGid(String(adminInfo.house_gid));
      setIsReady(true);
    } else if (isSuperAdmin && houseOptions) {
      setIsReady(true);
    }
  }, [isStoreAdmin, isSuperAdmin, adminInfo, houseOptions]);

  // 转换店铺选项为 MobileSelectOption 格式
  const formattedHouseOptions: MobileSelectOption[] = useMemo(() => {
    if (!houseOptions) return [];
    // houseOptions 本身就是数字数组
    return houseOptions.map((houseId: number) => ({
      label: `店铺 ${houseId}`,
      value: String(houseId),
    }));
  }, [houseOptions]);

  return {
    // 当前选中的店铺号
    houseGid,
    // 设置店铺号（仅超级管理员可用）
    setHouseGid: isSuperAdmin ? setHouseGid : () => {},
    // 是否为超级管理员
    isSuperAdmin,
    // 是否为店铺管理员
    isStoreAdmin,
    // 店铺选项列表（超级管理员）
    houseOptions: formattedHouseOptions,
    // 是否正在加载
    loading: loadingAdmin || loadingOptions,
    // 是否已准备好（已获取到店铺信息）
    isReady,
    // 是否可以手动输入/选择店铺
    canSelectHouse: isSuperAdmin,
  };
}
