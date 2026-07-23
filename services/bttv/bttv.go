package bttv

type BTTVService struct{

}

func (bttv *BTTVService) GetGlobalEmotes() []BTTVEmote {
	panic("not implemented")
}

func (bttv *BTTVService) GetUser(platform, channel string) *BTTVUser {
	panic("not implemented")
}
