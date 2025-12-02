declare namespace API {
  /**
   * 费用规则
   */
  type FeeRule = {
    threshold: number; // 分数阈值
    fee: number;       // 费用（分）
    kind?: string;     // 游戏类型（不设置表示全局规则）
    base?: number;     // 底分（不设置表示全局规则）
  };

  /**
   * 费用配置
   */
  type FeesConfig = {
    rules: FeeRule[];
  };

  /**
   * 查询店铺费用 - 请求参数
   */
  type GetShopFeesParams = {
    house_gid: number;
  };

  /**
   * 店铺费用配置 - 响应数据
   */
  type ShopFeesResult = {
    house_gid: number;
    fees_json: string;   // JSON 字符串
    share_fee: boolean;  // 分运开关
    push_credit: number; // 推送额度（分）
  };

  /**
   * 设置店铺费用 - 请求参数
   */
  type SetShopFeesParams = {
    house_gid: number;
    fees_json: string; // JSON 字符串，格式：{"rules": [...]}
  };

  /**
   * 设置分运开关 - 请求参数
   */
  type SetShareFeeParams = {
    house_gid: number;
    share: boolean;
  };

  /**
   * 设置推送额度 - 请求参数
   */
  type SetPushCreditParams = {
    house_gid: number;
    credit: number;
  };
}
