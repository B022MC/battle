## 脱微信依赖迁移记录与功能梳理（仅非微信）

### 术语澄清
- **“值班” = 成为管理员**：在本项目语义中，“值班”指将某个用户设为管理员（具有圈子管理与运营权限）。

### 圈子（Group）功能设计要点
- **约束**
  - 一名管理员只管理一个圈子（管理员:圈子 = 1:1）。
  - 普通用户任意时刻至多属于一个圈子（用户:活跃圈子 = N:1）。
  - 加入前置校验：用户当前不在其它活跃圈子。
- **数据模型建议**
  - groups：id, admin_user_id(唯一), name, created_at, updated_at。
  - group_members：id, group_id, user_id(唯一-仅针对left_at IS NULL), joined_at, left_at(NULL表示活跃), status。
  - 约束：group_members 对 (user_id, left_at IS NULL) 建部分唯一索引，防并发重复入圈。
- **核心接口（REST，幂等）**
  - POST /groups → 管理员创建自己的圈子（若已存在返回已有）。
  - POST /groups/{groupId}/members { user_id } → 将用户加入圈；若用户在其它活跃圈子返回冲突错误。
  - DELETE /groups/{groupId}/members/{userId} → 成员退圈（left_at=now）。
  - GET /groups/{groupId}/members?keyword=&page= → 成员搜索/分页。
- **权限/审计**
  - 仅圈子管理员/平台管理员可操作；关键操作写审计（操作者、对象、原因）。
- **前端交互（最小可行）**
  - 管理页：搜索用户→加入；成员列表→移除/搜索/分页；提示“用户已在其它圈子”。

### commands 对照（仅非微信相关）
- **店铺/圈子**
  - 管理店铺(CmdBindHouse)/替换店铺(CmdReplaceHouse)/退出店铺(CmdUnBindHouse)：后端已有店铺/中控/会话基础；前端入口待统一。
  - 绑圈(CmdBindGroup)/退圈(CmdUnBindGroup)/删圈(CmdDeleteGroup)/禁圈(CmdFreezeGroup)/解圈(CmdReleaseGroup)：需按上文圈子设计补全接口与界面。
- **费率与额度（运营参数）**
  - 运费(CmdSetFee)/分运(CmdSetShareFee)/取消分运(CmdUnsetShareFee)：需参数模型+设置接口+审计。
  - 额度(CmdSetCredit)/推送额度(CmdGetPushCredit)/查额度调整(CmdGetCredits)：钱包/资金查询部分已具，额度配置与变更记录待补。
  - 开多号(CmdOpenMultiGIDs)/关多号(CmdCloseMultiGIDs)，关闭解禁(CmdCloseFreeze)/打开解禁(CmdOpenFreeze)：需模型与接口落地。
- **运营与统计**
  - 值班(CmdManagerOnDuty)：即“设为管理员”；实现管理员角色授予/回收接口与界面。
  - 刷新(CmdRefresh)：复用现有拉取/同步接口即可。
  - 今日/昨日/上周统计(CmdGetTodayStat/CmdGetYesterdayStat/CmdGetLastweekStat)：后端 stats 已有，前端页已接。
  - 查分/查大/查小/查消分(CmdGetBalances/CmdGetBalancesBiggerThan/CmdGetBalancesLessThan/CmdGetDestroiedBalances)：钱包查询有基础，阈值筛选/删除账目查询需补API。
- **账号/申请流转**
  - 申请/通过/拒绝(CmdListApplications/CmdAggreeApplication/CmdRefuseApplication)：application 服务与前端列表/表单已接，完善流程/权限即可。
  - 删除用户(CmdDeletePlayer)：需提供软删/审计。
- **管理员运行控制（非微信调用层）**
  - 全禁/恢复(CmdFreezeAll/CmdRestoreAll)、查多号(CmdGetMultiGids)、踢(CmdKickOffGamer)：基于会话/成员维度提供统一控制接口与UI。
  - 房间/解散/查房/查桌/解桌(CmdListRooms/CmdDismissRoom/CmdCheckRoom/CmdCheckTable/CmdDismissTable)：已有桌台/会话/成员骨架，需补具体控制API与界面。
  - 退出登录(CmdLogoutUser)：会话管理已具，补统一登出流程。

> 已剔除：所有“微信相关”交互（例如 CmdSearchGID / CmdSearchWx 及“微信好友”分组命令）。

### 已有能力（非微信、非游戏层）
- 认证/RBAC、平台多租切库（base_platform + ConnPool）、HTTP装配与中间件、请求封装、通用UI库。
- 店铺/中控/会话基础服务与监控任务框架、申请/统计/钱包查询等基础能力。

### 优先路线图（建议）
1) 圈子管理端到端（表→仓储→用例→接口→前端页）。
2) 运营参数（运费/分运/额度等）模型、接口、审计、灰度。
3) 管理员运行控制（踢/全禁/恢复/解桌等）统一API与UI。
4) 可观测性与审计（操作日志、metrics、追踪）。

### 附：接口契约草案（圈子）
- POST /groups → 200 { id, admin_user_id, name }
- POST /groups/{id}/members { user_id } → 200/409
- DELETE /groups/{id}/members/{user_id} → 200
- GET /groups/{id}/members?keyword=&page=&page_size=


