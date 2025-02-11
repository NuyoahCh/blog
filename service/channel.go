// @Author Gopher
// @Date 2025/2/11 08:39:00
// @Desc 创建操作频道 Service
package service

import (
	"blog/common/global"
	"blog/models"
)

type ChannelService struct {
}

// AddChannel add
func (u *ChannelService) AddChannel(channel models.Channel) int64 {
	return global.Db.Table("channel").Create(&channel).RowsAffected
}

// DelChannel del
func (u *ChannelService) DelChannel(id int) int64 {
	return global.Db.Delete(&models.Channel{}, id).RowsAffected
}

// UpdateChannel update
func (u *ChannelService) UpdateChannel(channel models.Channel) int64 {
	return global.Db.Updates(&channel).RowsAffected
}

// GetChannel get
func (u *ChannelService) GetChannel(id int) models.Channel {
	var channel models.Channel
	global.Db.First(&channel, id)
	return channel
}

// GetChannelList get channel list
func (u *ChannelService) GetChannelList() []models.Channel {
	channelList := make([]models.Channel, 0)
	global.Db.Find(&channelList)
	return channelList
}
