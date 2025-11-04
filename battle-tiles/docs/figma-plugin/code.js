async function main() {
  await figma.loadFontAsync({ family: 'Inter', style: 'Regular' });
  await figma.loadFontAsync({ family: 'Inter', style: 'Bold' });

  // 注意：Starter 计划限制页面数量，避免 createPage，这里复用 currentPage 并在其下创建/复用一个画板 Frame
  const page = figma.currentPage;
  let board = page.findOne(n => n.type === 'FRAME' && n.name === 'H5 Prototype');
  if (!board) {
    board = figma.createFrame();
    board.name = 'H5 Prototype';
    board.fills = [{ type: 'SOLID', color: { r: 0.965, g: 0.97, b: 0.985 } }];
    page.appendChild(board);
  }
  // 清空旧内容
  if (board && 'children' in board) {
    const children = [...board.children];
    children.forEach(c => c.remove());
  }
  // 允许显示画板外的列标题
  if ('clipsContent' in board) board.clipsContent = false;

  const W = 375, H = 812;
  const APPBAR_H = 56, TABBAR_H = 60, PADDING = 12;
  const GAP_X = 420, GAP_Y = 920;

  const screens = [
    { role: 'admin', id: 'page-home', title: '首页（本店）', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'admin', id: 'page-shop', title: '房间', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'admin', id: 'page-room-detail', title: '房间详情', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'admin', id: 'page-members', title: '成员', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'admin', id: 'page-funds', title: '积分', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'admin', id: 'page-admin-apps', title: '入馆申请', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'admin', id: 'page-admin-my', title: '我的', tabs: ['首页','房间','成员','积分','我的'] },
    { role: 'root', id: 'root-home', title: '首页（全局）', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-shops', title: '店铺选择', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-rooms', title: '房间（选店后）', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-room-detail', title: '房间详情', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-ctrl', title: '中控账号', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-my', title: '我的', tabs: ['首页','店铺','中控','我的'] },
    // 超管具备管理员全部权限 → 补充管理员视图到 root 列
    { role: 'root', id: 'root-members', title: '成员（选店或本店）', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-funds', title: '积分（选店或本店）', tabs: ['首页','店铺','中控','我的'] },
    { role: 'root', id: 'root-admin-apps', title: '入馆申请', tabs: ['首页','店铺','中控','我的'] },
    { role: 'user', id: 'user-accounts', title: '我的游戏账号', tabs: ['账号','钱包','我的'] },
    { role: 'user', id: 'user-wallet', title: '我的余额与流水', tabs: ['账号','钱包','我的'] },
    { role: 'user', id: 'user-my', title: '我的', tabs: ['账号','钱包','我的'] },
  ];

  function rect(parent, x, y, w, h, fill, name) {
    const r = figma.createRectangle();
    r.resize(w, h);
    r.x = x; r.y = y; r.name = name || 'Rect';
    r.fills = [{ type: 'SOLID', color: fill || { r: 1, g: 1, b: 1 } }];
    parent.appendChild(r);
    return r;
  }
  function text(parent, x, y, str, size, bold, color, name) {
    const t = figma.createText();
    t.name = name || 'Text';
    t.x = x; t.y = y;
    t.characters = str;
    t.fontName = { family: 'Inter', style: bold ? 'Bold' : 'Regular' };
    t.fontSize = size || 14;
    t.fills = [{ type: 'SOLID', color: color || { r: 0.12, g: 0.14, b: 0.16 } }];
    parent.appendChild(t);
    return t;
  }
  function card(parent, yTop, title, lines) {
    const x = PADDING, w = W - PADDING * 2;
    const lineHeight = 20;
    const headerH = 30;
    const contentH = Math.max(40, ((lines && lines.length) || 0) * lineHeight + 12);
    const h = headerH + contentH + 12;
    const cardRect = rect(parent, x, yTop, w, h, { r: 1, g: 1, b: 1 }, 'Card');
    cardRect.strokes = [{ type: 'SOLID', color: { r: 0.94, g: 0.95, b: 0.96 } }];
    cardRect.strokeWeight = 1;
    text(parent, x + 12, yTop + 10, title, 13, false, { r: 0.42, g: 0.45, b: 0.50 }, 'CardTitle');
    (lines || []).forEach((ln, i) => {
      text(parent, x + 12, yTop + headerH + 8 + i * lineHeight, ln, 14, false);
    });
    return yTop + h + 10;
  }
  // UI 基础组件
  function chip(parent, x, y, textStr, fill) {
    const padX = 8, h = 22;
    const t = text(parent, x + padX, y + 4, textStr, 11, true, { r: 0.25, g: 0.29, b: 0.36 });
    const w = Math.round(t.width + padX * 2);
    const bg = rect(parent, x, y, w, h, fill || { r: 0.93, g: 0.95, b: 1.0 }, 'Chip');
    bg.cornerRadius = 6;
    t.x = x + padX; t.y = y + (h - t.height) / 2 - 1;
    return { w, h };
  }
  function button(parent, x, y, label, variant) {
    const isPrimary = variant === 'primary';
    const fill = isPrimary ? { r: 0.24, g: 0.49, b: 1.0 } : { r: 1, g: 1, b: 1 };
    const textColor = isPrimary ? { r: 1, g: 1, b: 1 } : { r: 0.12, g: 0.14, b: 0.16 };
    const border = isPrimary ? null : { type: 'SOLID', color: { r: 0.9, g: 0.92, b: 0.96 } };
    const padX = 14, padY = 10;
    const t = text(parent, x + padX, y + padY, label, 12, true, textColor);
    const w = Math.round(t.width + padX * 2), h = Math.round(t.height + padY * 2);
    const bg = rect(parent, x, y, w, h, fill, 'Btn');
    bg.cornerRadius = 10; if (border) { bg.strokes = [border]; bg.strokeWeight = 1; }
    t.x = x + padX; t.y = y + padY - 2;
    return { w, h };
  }
  function input(parent, x, y, w, placeholder) {
    const h = 40; const bg = rect(parent, x, y, w, h, { r: 1, g: 1, b: 1 }, 'Input');
    bg.strokes = [{ type: 'SOLID', color: { r: 0.9, g: 0.92, b: 0.96 } }]; bg.strokeWeight = 1; bg.cornerRadius = 10;
    text(parent, x + 12, y + 12, placeholder, 12, false, { r: 0.5, g: 0.52, b: 0.56 });
  }
  function listCard(parent, yTop, title, metaLines, chips, actions) {
    const x = PADDING, w = W - PADDING * 2, h = 92;
    const cardRect = rect(parent, x, yTop, w, h, { r: 1, g: 1, b: 1 }, 'ListCard');
    cardRect.strokes = [{ type: 'SOLID', color: { r: 0.94, g: 0.95, b: 0.96 } }]; cardRect.strokeWeight = 1; cardRect.cornerRadius = 12;
    text(parent, x + 12, yTop + 10, title, 15, true);
    (metaLines || []).forEach((ln, i) => text(parent, x + 12, yTop + 34 + i * 16, ln, 12, false, { r: 0.42, g: 0.45, b: 0.50 }));
    let cx = x + w - 12;
    (actions || []).reverse().forEach(label => { const b = button(parent, cx - 76, yTop + 50, label, 'ghost'); cx -= 84; });
    let chipX = x + 12; (chips || []).forEach(txt => { const c = chip(parent, chipX, yTop + 60, txt); chipX += c.w + 8; });
    return yTop + h + 10;
  }
  // 玩法名称映射
  function kindName(id) {
    switch (id) {
      case 60: return '丁二红';
      case 110: return '红二十';
      case 57: return '三人断勾卡';
      case 115: return '跑得快';
      case 95: return '斗十四';
      case 114: return '断勾卡';
      case 150: return '红中';
      default: return 'unknown';
    }
  }
  function appbar(frame, title) {
    const bar = rect(frame, 0, 0, W, APPBAR_H, { r: 1, g: 1, b: 1 }, 'AppBar');
    bar.strokes = [{ type: 'SOLID', color: { r: 0.93, g: 0.94, b: 0.95 } }];
    bar.strokeWeight = 1;
    text(frame, 16, 16, title, 18, true);
  }
  function tabbar(frame, tabs) {
    const y = H - TABBAR_H;
    const bar = rect(frame, 0, y, W, TABBAR_H, { r: 1, g: 1, b: 1 }, 'TabBar');
    bar.strokes = [{ type: 'SOLID', color: { r: 0.93, g: 0.94, b: 0.95 } }];
    bar.strokeWeight = 1;
    const colW = W / Math.max(1, tabs.length);
    tabs.forEach((t, i) => {
      text(frame, Math.round(i * colW + colW / 2 - (t.length * 6)), y + 20, t, 11, false, { r: 0.24, g: 0.49, b: 1.0 });
    });
  }
  function roleBadge(frame, role) {
    var label = role === 'admin' ? '管理员可见' : (role === 'root' ? '超管可见' : '普通用户可见');
    var color = role === 'admin' ? { r: 0.24, g: 0.49, b: 1.0 } : (role === 'root' ? { r: 0.56, g: 0.27, b: 0.68 } : { r: 0.16, g: 0.66, b: 0.36 });
    const t = text(frame, 0, 12, label, 11, true, { r: 1, g: 1, b: 1 }, 'RoleBadgeText');
    // 计算宽度后放置背景
    const padX = 8, padY = 4;
    const w = Math.round(t.width + padX * 2), h = 22;
    const x = W - 12 - w, y = 12;
    const bg = rect(frame, x, y, w, h, color, 'RoleBadge');
    bg.cornerRadius = 6;
    // 移动文字到背景内部
    t.x = x + padX; t.y = y + (h - t.height) / 2 - 1;
  }
  function roleHeader(boardFrame, role, x) {
    var label = role === 'admin' ? '管理员（Admin）' : (role === 'root' ? '超级管理员（Root）' : '普通用户（User）');
    var color = role === 'admin' ? { r: 0.24, g: 0.49, b: 1.0 } : (role === 'root' ? { r: 0.56, g: 0.27, b: 0.68 } : { r: 0.16, g: 0.66, b: 0.36 });
    const tag = rect(boardFrame, x, -44, 160, 30, color, 'RoleHeader');
    tag.cornerRadius = 8;
    text(boardFrame, x + 10, -41, label, 13, true, { r: 1, g: 1, b: 1 }, 'RoleHeaderText');
  }
  function renderContent(frame, screen) {
    let y = APPBAR_H + 12;
    const gray = { r: 0.42, g: 0.45, b: 0.50 };
    if (screen.id === 'page-home') {
      // 根据 resp.StatsVO 设计首页卡片
      y = card(frame, y, '时间范围', [
        'house_gid: number',
        'range_start: ISO8601',
        'range_end: ISO8601',
      ]);
      y = card(frame, y, '会话', [
        'session.active（当前在线会话数）',
      ]);
      y = card(frame, y, '钱包汇总', [
        'wallet.balance_total（余额总和）',
        'wallet.members（成员总数）',
        'wallet.low_balance_members（低余额成员数）',
      ]);
      y = card(frame, y, '流水汇总', [
        'ledger.income（上分总额）',
        'ledger.payout（下分总额）',
        'ledger.adjust（调整）',
        'ledger.net（净变动）',
        'ledger.records（流水条数）',
        'ledger.members_involved（参与成员数）',
      ]);
      y = card(frame, y, '统计接口', [
        'POST /stats/today｜/stats/yesterday｜/stats/week｜/stats/lastweek',
      ]);
    } else if (screen.id === 'page-shop') {
      input(frame, PADDING, y, W - PADDING * 2 - 80, '搜索 table_id / mapped_num'); button(frame, W - PADDING - 70, y, '查询', 'primary'); y += 56;
      chip(frame, PADDING, y, '玩法 #150'); chip(frame, PADDING + 90, y, '底 8'); y += 32;
      y = listCard(frame, y, '房间 #177', ['映射号 755765 · 圈 60870'], ['玩法 #150', '底 10'], ['查桌','解桌','详情']);
      y = listCard(frame, y, '房间 #176', ['映射号 440446 · 圈 60870'], ['玩法 #150', '底 8'], ['查桌','解桌','详情']);
      button(frame, PADDING, y, '全选', 'ghost'); button(frame, W - PADDING - 92, y, '批量解桌', 'primary'); y += 54;
    } else if (screen.id === 'page-members') {
      y = card(frame, y, '成员列表', ['搜ID/昵称/圈号', '行内：余额、操作(上/下/强下)']);
      y = card(frame, y, '接口映射（成员）', [
        'POST /shops/members/list',
        'POST /shops/members/pull',
        'POST /shops/members/kick',
        'POST /shops/members/logout',
        'POST /shops/diamond/query',
      ]);
      y = card(frame, y, '成员详情（示例）', ['余额：¥12,300', '最近流水：+500 / -200 ...']);
      y = card(frame, y, '钱包响应（WalletVO）', [
        'house_gid: int32',
        'member_id: int32',
        'balance: int32',
        'forbid: bool',
        'limit_min: int32',
        'updated_at: time',
        'updated_by: int32',
      ]);
      y = card(frame, y, '流水响应（LedgerVO）', [
        'id: int32, house_gid: int32, member_id: int32',
        'change_amount: int32, balance_before/after: int32',
        'type: int32(1上分/2下分/3强下/4调整)',
        'reason: string, operator_user_id: int32, biz_no: string',
        'created_at: time',
      ]);
    } else if (screen.id === 'page-funds') {
      y = card(frame, y, '积分操作（接口）', [
        'POST /members/credit/deposit',
        'POST /members/credit/withdraw',
        'POST /members/credit/force_withdraw',
        'PATCH /members/limit',
      ]);
      y = card(frame, y, '钱包/流水（接口）', [
        'POST /members/wallet/get',
        'POST /members/wallet/list',
        'POST /members/ledger/list',
      ]);
      y = card(frame, y, '钱包响应（WalletVO）', [
        'house_gid, member_id, balance, forbid, limit_min',
        'updated_at, updated_by',
      ]);
      y = card(frame, y, '流水响应（LedgerVO）', [
        'id, house_gid, member_id, change_amount, balance_before/after',
        'type(1上/2下/3强/4调), reason, operator_user_id, biz_no, created_at',
      ]);
    } else if (screen.id === 'page-admin-apps') {
      y = card(frame, y, '入馆申请（接口）', [
        'POST /shops/applications/list',
        'POST /applications/approve',
        'POST /applications/reject',
      ]);
      y = card(frame, y, '申请列表响应（ApplicationsVO）', [
        'items: ApplicationItemVO[]',
        'ApplicationItemVO: id, status, applier_id, applier_gid, applier_name, house_gid, created_at(ms)',
      ]);
    } else if (screen.id === 'page-admin-my') {
      y = card(frame, y, '我的', ['值班开关/通知配置/会话状态/退出']);
    } else if (screen.id === 'root-home') {
      y = card(frame, y, '概览（示例）', ['可展示全局 active 会话、全局资金趋势等']);
      y = card(frame, y, '运维接口', [
        'GET /ops/plaza/metrics',
        'GET /ops/plaza/health',
      ]);
    } else if (screen.id === 'root-shops') {
      input(frame, PADDING, y, W - PADDING * 2 - 80, '搜索店铺名或 #house_gid'); button(frame, W - PADDING - 70, y, '搜索', 'primary'); y += 56;
      y = listCard(frame, y, '天天部落_XD  #58959', ['会话：在线 · 管理员 2 · 中控 1'], ['店铺'], ['进入房间']);
      y = listCard(frame, y, '星河会馆  #60001', ['会话：离线 · 管理员 1 · 中控 0'], ['店铺'], ['进入房间']);
      y = card(frame, y, '店铺管理相关接口', [
        'POST /shops/admins (指派/撤销/列表)',
        'POST /shops/groups (forbid/unforbid/bind/unbind/delete)',
        'POST /shops/applications/list',
      ]);
      y = card(frame, y, '管理员响应（ShopAdminVO）', [
        'id: int32, house_gid: int32, user_id: int32, role: admin|operator',
      ]);
    } else if (screen.id === 'root-rooms') {
      chip(frame, PADDING, y, '当前：天天部落_XD #58959', { r: 0.86, g: 0.93, b: 1.0 }); button(frame, W - PADDING - 80, y, '切换店铺', 'ghost'); y += 34;
      input(frame, PADDING, y, W - PADDING * 2 - 80, '搜索 table_id / mapped_num'); button(frame, W - PADDING - 70, y, '查询', 'primary'); y += 56;
      chip(frame, PADDING, y, '玩法 #115'); chip(frame, PADDING + 90, y, '底 5'); y += 32;
      y = listCard(frame, y, '房间 #53', ['映射号 556387 · 圈 60870'], ['玩法 #115', '底 5'], ['查桌','解桌','详情']);
      y = listCard(frame, y, '房间 #6', ['映射号 247984 · 圈 60870'], ['玩法 #110', '底 1'], ['查桌','解桌','详情']);
      y = listCard(frame, y, '房间 #0（未分配）', ['映射号 744165 · 圈 60870'], ['玩法 #95', '底 10'], ['查桌','解桌','详情']);
      button(frame, PADDING, y, '全选', 'ghost'); button(frame, W - PADDING - 92, y, '批量解桌', 'primary'); y += 54;
    } else if (screen.id === 'root-members') {
      y = card(frame, y, '成员（与管理员一致）', [
        '接口：/shops/members/list|pull|kick|logout、/shops/diamond/query',
        '显示余额/操作（上/下/强下）、详情进入流水',
      ]);
      y = card(frame, y, 'WalletVO / LedgerVO', [
        'wallet：house_gid, member_id, balance, forbid, limit_min, updated_*',
        'ledger：id, change_amount, balance_before/after, type, reason, operator, biz_no, created_at',
      ]);
    } else if (screen.id === 'root-funds') {
      y = card(frame, y, '积分操作（与管理员一致）', [
        '接口：/members/credit/deposit|withdraw|force_withdraw、PATCH /members/limit',
        '钱包/流水：/members/wallet/get|list、/members/ledger/list',
      ]);
    } else if (screen.id === 'root-admin-apps') {
      y = card(frame, y, '入馆申请（与管理员一致）', [
        '接口：/shops/applications/list、/applications/approve|reject',
        '响应：items: ApplicationItemVO[]',
      ]);
    } else if (screen.id === 'root-ctrl') {
      y = card(frame, y, '中控账号接口', [
        'POST /shops/ctrlAccounts',
        'POST /shops/ctrlAccounts/bind',
        'DELETE /shops/ctrlAccounts/bind',
        'POST /shops/ctrlAccounts/list',
        'POST /shops/ctrlAccounts/listAll',
      ]);
      y = card(frame, y, '列表响应（CtrlAccountVO/CtrlAccountAllVO）', [
        'CtrlAccountVO: id, house_gid, login_mode, identifier, status',
        'CtrlAccountAllVO: id, login_mode, identifier, status, last_verify_at, houses[]',
      ]);
    } else if (screen.id === 'root-my') {
      y = card(frame, y, '我的', ['账户信息 / 权限说明 / 退出登录']);
    } else if (screen.id === 'user-accounts') {
      y = card(frame, y, '我的游戏账号（接口）', [
        'POST /game/accounts/verify',
        'POST /game/accounts',
        'GET /game/accounts/me',
        'DELETE /game/accounts/me',
      ]);
      y = card(frame, y, '账号响应（AccountVO）', [
        'id: int32, account: string, nickname: string',
        'is_default: bool, status: int32, login_mode: string',
      ]);
      y = card(frame, y, '会话（接口）', [
        'POST /game/accounts/sessionStart',
        'POST /game/accounts/sessionStop',
      ]);
    } else if (screen.id === 'user-wallet') {
      y = card(frame, y, '我的钱包/流水（接口）', [
        'POST /members/wallet/get',
        'POST /members/ledger/list',
      ]);
    } else if (screen.id === 'user-my') {
      y = card(frame, y, '我的', ['申请入馆（与管理端审批联动）', '退出登录']);
    } else if (screen.id === 'page-game') {
      y = card(frame, y, '会话控制', ['开始会话：选择类型/店铺/备注', '停止会话：选择进行中会话']);
    }
  }

  const colOf = (role) => (role === 'admin' ? 0 : role === 'root' ? 1 : 2);
  const counters = { admin: 0, root: 0, user: 0 };

  // 预估画板尺寸：三列（admin/root/user），行数为各自数量最大值
  const countByRole = screens.reduce((m, s) => (m[s.role] = (m[s.role] || 0) + 1, m), {});
  const maxRows = Math.max(countByRole.admin || 0, countByRole.root || 0, countByRole.user || 0);
  const boardWidth = 2 * GAP_X + W; // 三列：x=0, GAP_X, 2*GAP_X，宽度以第三列+W 估算
  const boardHeight = (Math.max(1, maxRows) - 1) * GAP_Y + H;
  board.resize(boardWidth, boardHeight);
  // 放置列标题
  roleHeader(board, 'admin', colOf('admin') * GAP_X);
  roleHeader(board, 'root', colOf('root') * GAP_X);
  roleHeader(board, 'user', colOf('user') * GAP_X);

  screens.forEach((s) => {
    const frame = figma.createFrame();
    frame.name = `${s.role}/${s.id}`;
    frame.resize(W, H);
    frame.x = colOf(s.role) * GAP_X;
    frame.y = counters[s.role] * GAP_Y;
    counters[s.role] += 1;
    frame.fills = [{ type: 'SOLID', color: { r: 0.965, g: 0.97, b: 0.985 } }];
    board.appendChild(frame);
    appbar(frame, s.title);
    renderContent(frame, s);
    tabbar(frame, s.tabs);
    roleBadge(frame, s.role);
  });

  figma.notify('已生成 H5 原型界面（三种角色）');
  figma.closePlugin();
}

main();



