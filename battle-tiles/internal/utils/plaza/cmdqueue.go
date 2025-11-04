package plaza

import (
	"battle-tiles/internal/dal/vo/game"
	"sync"
)

const (
	CmdTypeForbid         = 0
	CmdTypeGetGroupMember = 1
	CmdTypeDismissTable   = 3
	CmdTypeRespondApply   = 4
	CmdTypeDeleteMember   = 5
	CmdTypeQueryTable     = 6
	CmdTypeQueryDiamond   = 7
)

type GameCommand struct {
	Pack   *game.Packer
	Type   int
	Key    string
	Member int
	//Flag   bool
}

type GameCmdQueue struct {
	mut      sync.RWMutex
	commands []*GameCommand
}

func (g *GameCmdQueue) Push(cmd *GameCommand) {
	g.mut.Lock()
	found := false
	for _, c := range g.commands {
		if c.Key == cmd.Key {
			found = true
			break
		}
	}
	if !found {
		g.commands = append(g.commands, cmd)
	}
	g.mut.Unlock()
}

func (g *GameCmdQueue) AddHead(cmd *GameCommand) {
	g.mut.Lock()
	var tmp []*GameCommand
	tmp = append(tmp, cmd)
	tmp = append(tmp, g.commands...)
	g.commands = tmp
	g.mut.Unlock()
}

func (g *GameCmdQueue) Remove(key string) {
	g.mut.Lock()
	defer g.mut.Unlock()

	var tmp []*GameCommand
	for _, cmd := range g.commands {
		if cmd.Key == key {
			continue
		}
		tmp = append(tmp, cmd)
	}
	g.commands = tmp
}

func (g *GameCmdQueue) Pop() *GameCommand {
	g.mut.Lock()
	defer g.mut.Unlock()
	if len(g.commands) == 0 {
		return nil
	}

	cmd := g.commands[0]
	if len(g.commands) != 0 {
		g.commands = g.commands[1:]
	}

	return cmd
}

func (g *GameCmdQueue) Last() *GameCommand {
	g.mut.Lock()
	defer g.mut.Unlock()
	if len(g.commands) == 0 {
		return nil
	}

	return g.commands[len(g.commands)-1]
}

func (g *GameCmdQueue) Top() *GameCommand {
	g.mut.Lock()
	defer g.mut.Unlock()
	if len(g.commands) == 0 {
		return nil
	}

	return g.commands[0]
}

func (g *GameCmdQueue) Clear() {
	g.mut.Lock()
	g.commands = []*GameCommand{}
	g.mut.Unlock()
}
