package entity

type Player struct {
	Uid int64
}

func (p *Player) UserId() int64 {
	return p.Uid
}
