package plaza

import (
	"battle-bot/internal/plaza/game"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// 读完整一帧：前4字节头(0:ver,1:chk,2..3:len[LE]) + body
func readFullPacket(conn net.Conn) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	total := int(binary.LittleEndian.Uint16(header[2:4]))
	if total < 4 || total > 64*1024 {
		return nil, fmt.Errorf("invalid length: %d", total)
	}
	buf := make([]byte, total)
	copy(buf, header)
	if _, err := io.ReadFull(conn, buf[4:]); err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return buf, nil
}

func dial82WithCtx(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
	d := &net.Dialer{Timeout: timeout}
	con, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	if tcp, ok := con.(*net.TCPConn); ok {
		_ = tcp.SetKeepAlive(true)
		_ = tcp.SetKeepAlivePeriod(30 * time.Second)
	}
	return con, nil
}

func GetUserInfoByAccountCtx(ctx context.Context, server82 string, account, pwdMD5 string) (*game.UserLogonInfo, error) {
	cmd := CmdAccountLogon(account, pwdMD5)

	enc := Encoder{}
	out := enc.Encrypt(cmd)

	con, err := dial82WithCtx(ctx, server82, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("dial82: %w", err)
	}
	defer con.Close()

	// 读写整体设置超时（也可用 Set{Read,Write}Deadline）
	_ = con.SetDeadline(time.Now().Add(8 * time.Second))

	if _, err = con.Write(out); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	reply, err := readFullPacket(con)
	if err != nil {
		return nil, err
	}

	pk, err := enc.Decrypt(reply)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	// 约定：subCmdID==100 登录成功（你原来的逻辑）
	if pk.Head.Cmd.SubCmdID != 100 {
		f := ParseLogonFailure(pk.Data())
		if f != nil {
			return nil, fmt.Errorf("login failed: %s", f.Desc)
		}
		return nil, fmt.Errorf("login failed: subCmdID=%d", pk.Head.Cmd.SubCmdID)
	}
	return ParseUserLogon(pk.Data()), nil
}

func GetUserInfoByMobileCtx(ctx context.Context, server82 string, mobile, pwdMD5 string) (*game.UserLogonInfo, error) {
	cmd := CmdMobileLogon(mobile, pwdMD5)

	enc := Encoder{}
	out := enc.Encrypt(cmd)

	con, err := dial82WithCtx(ctx, server82, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("dial82: %w", err)
	}
	defer con.Close()

	_ = con.SetDeadline(time.Now().Add(8 * time.Second))

	if _, err = con.Write(out); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	reply, err := readFullPacket(con)
	if err != nil {
		return nil, err
	}
	pk, err := enc.Decrypt(reply)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	if pk.Head.Cmd.SubCmdID != 100 {
		f := ParseLogonFailure(pk.Data())
		if f != nil {
			return nil, fmt.Errorf("login failed: %s", f.Desc)
		}
		return nil, fmt.Errorf("login failed: subCmdID=%d", pk.Head.Cmd.SubCmdID)
	}
	return ParseUserLogon(pk.Data()), nil
}
