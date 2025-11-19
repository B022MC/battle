# 老项目命令 vs 新项目实现对比分析

## 📋 命令列表

| 命令 | 老项目命令 | 功能说明 | 是否需要游戏接口 | 新项目API |
|------|-----------|---------|----------------|----------|
| 额度 | `CmdSetCredit` | 查询/设置游戏额度 | ❌ 否 | `/shops/fees/*` |
| 分运 | `CmdSetShareFee` | 开启分运费 | ❌ 否 | `/shops/sharefee/set` |
| 取消分运 | `CmdUnsetShareFee` | 关闭分运费 | ❌ 否 | `/shops/sharefee/set` |
| 申请 | `CmdListApplications` | 查看申请列表 | ❌ 否（内存读取） | `/shops/game-applications/list` |
| 通过 | `CmdAggreeApplication` | 通过申请 | ✅ **是** | `/shops/game-applications/approve` |
| 拒绝 | `CmdRefuseApplication` | 拒绝申请 | ✅ **是** | `/shops/game-applications/reject` |
| 上 | `CmdAddUserScore` | 上分 | ❌ 否 | `/members/credit/deposit` |
| 下 | `CmdReduceUserScore` | 下分 | ❌ 否 | `/members/credit/withdraw` |

---

## 🔍 详细分析

### 1. 额度（CmdSetCredit）❌ 不需要游戏接口

#### 老项目实现
**文件**: `service_house_command_handler.go:497-562`

```go
func (that *Service) handleCmdSetCredit(manager *Manager, txt string) {
    // 解析参数
    ps := that.parseCmdParams(txt, CmdSetCredit)
    
    switch len(ps) {
    case 0:
        // 查询：从本地数据读取
        for _, credit := range manager.house.GetGameCredits(manager.model.PlayGroup) {
            lines = append(lines, fmt.Sprintf("🈲 %d/%s/%d", 
                credit.Credit/100, 
                plaza.GetKindName(credit.Kind), 
                credit.BaseScore))
        }
        
    case 3:
        // 设置：保存到本地数据库
        credit, _ := strconv.Atoi(ps[0])
        kindName := ps[1]
        base, _ := strconv.Atoi(ps[2])
        manager.house.SetGameCredit(manager.model.PlayGroup, credit, 
            plaza.GetKindID(kindName), base)
    }
}
```

**特点**：
- ✅ 纯本地数据库操作
- ✅ 不需要请求游戏服务器
- ✅ 只是存储和查询配置

#### 新项目实现
- **查询**: `GET /shops/fees/get?house_gid={gid}`
- **设置**: `POST /shops/fees/update` 
- **数据表**: `game_house_settings.fees_json`

---

### 2. 分运/取消分运（CmdSetShareFee/CmdUnsetShareFee）❌ 不需要游戏接口

#### 老项目实现
**文件**: `service_house_command_handler.go:577-595`

```go
func (that *Service) handleCmdSetShareFee(manager *Manager) {
    // 直接更新数据库字段
    manager.house.model.ShareFee = true
    if err := that.db.Save(manager.house.model).Error; err != nil {
        hlogger.Error("Failed to set share fee", err)
        rob.SendText(robot.FileHelper, "失败")
        return
    }
    rob.SendText(robot.FileHelper, "成功")
}
```

**特点**：
- ✅ 纯数据库操作
- ✅ 只是修改 `share_fee` 字段
- ✅ 不涉及游戏接口

#### 新项目实现
- **API**: `POST /shops/sharefee/set`
- **参数**: `{"house_gid": 123, "share_fee": true}`
- **数据表**: `game_house_settings.share_fee`

---

### 3. 申请（CmdListApplications）❌ 不需要游戏接口

#### 老项目实现
**文件**: `service_house_command_handler.go:918-930`

```go
func (that *Service) handleCmdListApplications(manager *Manager) {
    // 从内存中读取申请列表
    lines := manager.house.GetApplicationsList()
    if len(lines) > 0 {
        rob.SendText(robot.FileHelper, strings.Join(lines, "\n"))
    } else {
        rob.SendText(robot.FileHelper, "无")
    }
}

// house_api.go:625-633
func (that *House) GetApplicationsList() []string {
    // 从内存 map 中读取
    apps := that._getApplications()
    var lines []string
    for _, app := range apps {
        lines = append(lines, app.String())
    }
    return lines
}
```

**特点**：
- ✅ 从内存中读取（`applyInfos sync.Map`）
- ✅ 申请数据由游戏服务器推送到内存
- ✅ 查询操作本身不需要请求游戏接口

#### 新项目实现
- **API**: `POST /shops/game-applications/list`
- **参数**: `{"house_gid": 123}`
- **数据来源**: 
  - 游戏服务器推送到 Plaza 内存
  - 后端从 Plaza 读取
  - 类似老项目的实现方式

---

### 4. 通过申请（CmdAggreeApplication）✅ **需要游戏接口**

#### 老项目实现
**文件**: `service_house_command_handler.go:932-955`

```go
func (that *Service) handleCmdAggreeApplication(manager *Manager, txt string) {
    gid, err := strconv.Atoi(ps[0])
    
    // 调用 house 的 RespondApplication
    if err := manager.house.RespondApplication(gid, true); err == nil {
        rob.SendText(robot.FileHelper, "完成")
    }
}

// house_api.go:635-646
func (that *House) RespondApplication(gid int, agree bool) error {
    // 从内存获取申请信息
    val, ok := that.applyInfos.Load(gid)
    applyInfo := val.(*plaza.ApplyInfo)
    
    // 发送到游戏服务器 ✅
    that.session.RespondApplication(applyInfo, agree)
    
    // 从内存删除
    that.applyInfos.Delete(applyInfo.ApplierGid)
    return nil
}

// session_api.go:33-38
func (that *Session) RespondApplication(applyInfo *ApplyInfo, agree bool) {
    // 发送TCP命令到游戏服务器 ✅
    that._87cmdQueue.Push(&GameCommand{
        Pack: CmdRespondApplication(that.userID, that.userPwd, 
            applyInfo.MessageId, applyInfo.HouseGid, 
            applyInfo.AplierId, agree),
        Type: CmdTypeRespondApply,
    })
}

// tcpcmd.go:136-139
func CmdRespondApplication(userId int, pwd string, msgId int, 
    houseGid int, applierGid int, agree bool) *Packer {
    packer := &Packer{}
    // 构造游戏协议包 ✅
    packer.SetCmd(MDM_GA_LOGIC_SERVICE, SUB_GA_APPLY_RESPOND)
    // ... 设置参数
    return packer
}
```

**特点**：
- ❌ **需要发送TCP命令到游戏服务器**
- ❌ 使用游戏协议：`MDM_GA_LOGIC_SERVICE` + `SUB_GA_APPLY_RESPOND`
- ❌ 必须通过游戏接口才能完成审批

#### 新项目实现
- **API**: `POST /shops/game-applications/approve`
- **参数**: `{"house_gid": 123, "message_id": 456}`
- **后端处理**: 
  - ✅ **需要调用游戏服务器接口**
  - ✅ 发送 TCP 命令到游戏服务器
  - ✅ 保存审批记录到 `game_shop_application_log`

---

### 5. 拒绝申请（CmdRefuseApplication）✅ **需要游戏接口**

#### 老项目实现
**文件**: `service_house_command_handler.go:957-980`

```go
func (that *Service) handleCmdRefuseApplication(manager *Manager, txt string) {
    gid, err := strconv.Atoi(ps[0])
    
    // 调用 RespondApplication，agree=false
    if err := manager.house.RespondApplication(gid, false); err == nil {
        rob.SendText(robot.FileHelper, "完成")
    }
}
```

**特点**：
- ❌ **与通过申请相同，需要游戏接口**
- ❌ 只是 `agree` 参数不同
- ❌ 底层都是调用 `session.RespondApplication()`

#### 新项目实现
- **API**: `POST /shops/game-applications/reject`
- **参数**: `{"house_gid": 123, "message_id": 456}`
- **后端处理**: 与通过申请相同，需要调用游戏接口

---

### 6. 上分（CmdAddUserScore）❌ 不需要游戏接口

#### 老项目实现
**文件**: `service_user_command_handler.go:501-558`

```go
func (that *Service) handleCmdAddUserScore(manager *Manager, to string, txt string) {
    // 1. 获取好友信息
    f, ok := rob.GetFriendByUsername(to)
    
    // 2. 查找玩家ID
    id, ok := manager.house._getPlayerIdByWxKey(f.Wxkey)
    
    // 3. 检查圈子
    group, ok := manager.house._getPlayGroupById(id)
    
    // 4. 解析金额
    num, err := strconv.Atoi(ps[0])
    
    // 5. 充值（本地操作）
    bal, err := manager.house.RechargePlayer(manager.model.PlayGroup, f.Wxkey, num, false)
}

// house_api.go:841-871
func (that *House) RechargePlayer(group string, wxKey string, number int, force bool) (float64, error) {
    // 1. 查找玩家ID
    id, ok := that._getPlayerIdByWxKey(wxKey)
    
    // 2. 检查是否在房间（下分时）
    if number < 0 && !force {
        if _, ok = that.userOnTableMap.Load(id); ok {
            return 0, errors.New("房间中")
        }
    }
    
    // 3. 更新余额（本地数据库）✅
    bal, ok := that._settlePlayerBalance(id, number*100)
    
    // 4. 保存充值记录
    that.db.Save(&model.TRechargeRecord{
        HouseGid:    that.model.GameId,
        PlayerId:    id,
        PlayGroup:   group,
        Amount:      number * 100,
        RechargedAt: time.Now().Unix(),
        Balance:     bal,
    })
    
    return float64(bal) / 100.0, nil
}

// house_utils.go:69-88
func (that *House) _settlePlayerBalance(id int, delta int) (int, bool) {
    // 加锁防止并发
    locker := that._getPlayerSettleLocker(id)
    locker.Lock()
    defer locker.Unlock()
    
    // 查询余额
    var player model.TPlayer
    that.db.Select("balance").Where("id=?", id).First(&player)
    
    // 更新余额（纯数据库操作）✅
    player.Balance += delta
    that.db.Model(&model.TPlayer{}).Where("id=?", id).Update("balance", player.Balance)
    
    return player.Balance, true
}
```

**特点**：
- ✅ **纯本地数据库操作**
- ✅ 只更新 `t_player.balance` 字段
- ✅ 保存充值记录到 `t_recharge_record`
- ✅ **不需要通知游戏服务器**

#### 新项目实现
- **API**: `POST /members/credit/deposit`
- **参数**: `{"house_gid": 123, "member_id": 456, "amount": 10000, "biz_no": "xxx"}`
- **后端处理**: 
  - ✅ 更新 `game_member.balance`
  - ✅ 保存记录到充值表
  - ✅ 不需要调用游戏接口

---

### 7. 下分（CmdReduceUserScore）❌ 不需要游戏接口

#### 老项目实现
**文件**: `service_user_command_handler.go:560-620`

```go
func (that *Service) handleCmdReduceUserScore(manager *Manager, to string, txt string) {
    // 与上分类似，只是金额为负数
    bal, err := manager.house.RechargePlayer(manager.model.PlayGroup, f.Wxkey, -num, false)
}
```

**特点**：
- ✅ 与上分完全相同，只是金额为负
- ✅ 会检查用户是否在房间中（`userOnTableMap`）
- ✅ 如果在房间，拒绝下分

#### 新项目实现
- **API**: `POST /members/credit/withdraw`
- **参数**: 与上分相同
- **后端处理**: 与上分相同，只是金额为负

---

## 📊 总结表

### 需要游戏接口的命令 ✅

| 命令 | 老项目 | 新项目 | 游戏协议 |
|------|-------|--------|---------|
| 通过申请 | `RespondApplication` | `/shops/game-applications/approve` | `MDM_GA_LOGIC_SERVICE` + `SUB_GA_APPLY_RESPOND` |
| 拒绝申请 | `RespondApplication` | `/shops/game-applications/reject` | `MDM_GA_LOGIC_SERVICE` + `SUB_GA_APPLY_RESPOND` |

### 不需要游戏接口的命令 ❌

| 命令 | 老项目 | 新项目 | 操作类型 |
|------|-------|--------|---------|
| 额度 | `SetGameCredit` | `/shops/fees/update` | 数据库配置 |
| 分运 | `ShareFee=true` | `/shops/sharefee/set` | 数据库配置 |
| 取消分运 | `ShareFee=false` | `/shops/sharefee/set` | 数据库配置 |
| 申请列表 | `GetApplicationsList` | `/shops/game-applications/list` | 内存读取 |
| 上分 | `RechargePlayer(+)` | `/members/credit/deposit` | 数据库余额 |
| 下分 | `RechargePlayer(-)` | `/members/credit/withdraw` | 数据库余额 |

---

## 🔧 新项目实现建议

### 1. 申请审批功能（需要游戏接口）

**关键点**：
- ✅ 新项目已实现 API：`/shops/game-applications/approve` 和 `/reject`
- ⚠️ 需要确认后端是否已实现向游戏服务器发送TCP命令
- ⚠️ 需要游戏服务器支持审批协议

**实现步骤**：
1. 前端调用审批 API
2. 后端接收审批请求
3. **后端构造游戏协议包**（TCP命令）
4. **发送到游戏服务器**
5. 保存审批记录到 `game_shop_application_log`
6. 返回结果给前端

### 2. 其他功能（不需要游戏接口）

**关键点**：
- ✅ 所有功能都已在新项目中实现
- ✅ 前端页面已完成
- ✅ 后端API已完成
- ✅ 只需要正常的HTTP API调用

---

## ⚠️ 重点注意

### 申请审批是唯一需要游戏接口的功能！

**老项目流程**：
```
微信命令 "通过 123456"
    ↓
Robot解析命令
    ↓
House.RespondApplication()
    ↓
Session.RespondApplication()
    ↓
构造TCP协议包 CmdRespondApplication()
    ↓
发送到游戏服务器 (87端口)
    ↓
游戏服务器处理申请
    ↓
推送结果到客户端
```

**新项目需要实现**：
```
前端点击"通过"按钮
    ↓
调用 POST /shops/game-applications/approve
    ↓
后端 game_shop_application.go
    ↓
【需要实现】构造游戏TCP协议包
    ↓
【需要实现】发送到游戏服务器
    ↓
保存记录到 game_shop_application_log
    ↓
返回结果
```

### 其他功能都是纯后端操作

- **额度设置**: 只是保存 JSON 配置
- **分运费**: 只是修改布尔值
- **上分下分**: 只是更新余额字段
- **申请列表**: 从内存读取（游戏服务器推送的数据）

这些都不需要主动请求游戏接口！

---

## 🎯 实现优先级

### P0 - 必须实现（需要游戏接口）
1. ✅ 申请审批（通过/拒绝）- 需要TCP命令到游戏服务器

### P1 - 已完成（不需要游戏接口）
2. ✅ 额度设置 - 已有API
3. ✅ 分运费 - 已有API
4. ✅ 上分下分 - 已有API
5. ✅ 申请列表 - 已有API（从内存读取）

所有功能中，**只有申请审批需要调用游戏接口**！
