package plaza

import (
	"battle-tiles/internal/consts"
	"battle-tiles/internal/dal/vo/game"
	"strings"
)

func CmdAccountLogon(account string, pwd string) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_MB_LOGON, consts.SUB_MB_LOGON_ACCOUNTS_LUA)
	packer.PushWord(consts.INVALID_ITEM)
	packer.PushWord(consts.MARKET_ID)
	packer.PushByte(consts.DEVICE_TYPE)
	packer.PushDWord(consts.APP_VERSION)
	packer.PushDWord(consts.CLIENT_VERSION)
	packer.PushDWord(consts.STATION_ID)
	packer.PushDWord(0)
	packer.PushString(pwd, consts.LEN_MD5)
	packer.PushString(account, consts.LEN_ACCOUNTS)
	packer.PushString(consts.MachineID(), consts.LEN_MACHINE_ID)
	packer.PushString("", consts.LEN_MOBILE_PHONE)

	packer.PaddingTo(249)
	return packer
}

func CmdForbidMember(user uint32, pwd string, group uint32, member uint32, forbid bool) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GA_LOGIC_SERVICE, consts.SUB_GA_UPDATE_MEMBER)
	packer.PushDWord(group)
	packer.PushDWord(member)

	packer.PushDWord(user)
	packer.PushString(pwd, consts.LEN_PASSWORD)

	packer.PushWord(consts.UPMEMBER_KIND_RIGHT)
	if forbid {
		packer.PushDWord(0x1000)
	} else {
		packer.PushDWord(0)
	}
	packer.PaddingTo(84)

	return packer
}

func CmdMobileLogon(mobile string, pwd string) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_MB_LOGON, consts.SUB_MB_LOGON_MOBILEPHONE)
	packer.PushWord(consts.INVALID_ITEM)
	packer.PushWord(consts.MARKET_ID)
	packer.PushByte(consts.DEVICE_TYPE)

	packer.PushDWord(consts.APP_VERSION)
	packer.PushDWord(consts.CLIENT_VERSION)

	packer.PushDWord(consts.STATION_ID)
	packer.PushDWord(0)

	packer.PushString(strings.ToUpper(pwd), consts.LEN_PASSWORD)
	packer.PushString(mobile, consts.LEN_MOBILE_PHONE)
	packer.PushString(strings.ToUpper(consts.MachineID()), consts.LEN_MACHINE_ID)

	packer.PaddingTo(185)

	return packer
}

func CmdMsgServerEnterMsg(dwUserID uint32) *game.Packer {
	packer := &game.Packer{}

	packer.SetCmd(consts.MDM_GA_MESSAGE_SERVICE, consts.SUB_GA_ENTER_MESSAGE)
	packer.PushDWord(dwUserID)
	packer.PushWord(consts.STATION_ID)
	packer.PaddingTo(8)

	return packer
}

func CmdLogonServer(userID uint32, pwdMD5 string) *game.Packer {
	packer := &game.Packer{}

	packer.SetCmd(consts.MDM_GA_LOGIC_SERVICE, consts.SUB_GA_LOGON_SERVER)
	packer.PushDWord(userID)
	packer.PushDWord(consts.STATION_ID)
	packer.PushDWord(4)
	packer.PushString(strings.ToUpper(pwdMD5), consts.LEN_PASSWORD)

	packer.PaddingTo(78)

	return packer
}

func CmdGroupService(userID uint32, groupID uint32) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GA_GROUP_SERVICE, consts.SUB_GA_ENTER_GROUP)
	packer.PushDWord(userID)
	packer.PushDWord(groupID)
	packer.PaddingTo(8)
	return packer
}

func CmdUserStandUp(tableID, chairID int) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GR_USER, consts.SUB_GR_USER_STANDUP)
	packer.PushWord(uint16(tableID))
	packer.PushWord(uint16(chairID))
	packer.PushByte(0)
	packer.PaddingTo(5)
	return packer
}

func CmdHeartBeat() *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(0, 1)
	packer.PaddingTo(0)
	return packer
}

func CmdDismissRoom(userID int, pwdMD5 string, kindID, mappedNum int) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GA_BATTLE_SERVICE, consts.SUB_GA_DISMISS_TABLE)
	packer.PushWord(uint16(kindID))
	packer.PushDWord(uint32(mappedNum))
	packer.PushDWord(uint32(userID))
	packer.PushString(strings.ToUpper(pwdMD5), consts.LEN_PASSWORD)
	packer.PaddingTo(76)
	return packer
}

func CmdRespondApplication(userId int, pwd string, msgId int, houseGid int, applierGid int, agree bool) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GA_LOGIC_SERVICE, consts.SUB_GA_APPLY_RESPOND)
	packer.PushDWord(uint32(msgId))
	packer.PushDWord(uint32(userId))
	packer.PushString(strings.ToUpper(pwd), consts.LEN_PASSWORD)
	packer.PushDWord(uint32(houseGid))
	packer.PushDWord(uint32(applierGid))
	if agree {
		packer.PushByte(byte(consts.APPLY_STATUS_AGREE))
	} else {
		packer.PushByte(byte(consts.APPLY_STATUS_REFUSE))
	}
	packer.PaddingTo(83)
	return packer
}

func CmdDeleteMember(userID int, pwdMD5 string, houseGid int, memId int) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GA_LOGIC_SERVICE, consts.SUB_GA_DELETE_MEMBER)
	packer.PushDWord(uint32(houseGid))
	packer.PushDWord(uint32(memId))
	packer.PushDWord(uint32(userID))
	packer.PushString(strings.ToUpper(pwdMD5), consts.LEN_PASSWORD)

	packer.PaddingTo(78)
	return packer
}

func CmdQueryTable(tabMappedNum int) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GA_BATTLE_SERVICE, consts.SUB_GA_QUERY_TABLE)
	packer.PushDWord(uint32(0))
	packer.PushDWord(uint32(tabMappedNum))
	packer.PaddingTo(6)
	return packer
}

func CmdQueryDiamond(userID int) *game.Packer {
	packer := &game.Packer{}
	packer.SetCmd(consts.MDM_GP_USER_SERVICE, consts.SUB_GP_QUERY_WEALTH_LUA)
	packer.PushDWord(uint32(userID))
	packer.PaddingTo(4)
	return packer
}
