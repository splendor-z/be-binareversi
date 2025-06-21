package db

import (
	"be-binareversi/model"
	"time"
)

// プレイヤーを作成
func CreatePlayer(player *model.Player) error {
	return DB.Create(player).Error
}

// IDでプレイヤーを取得
func GetPlayerByID(id string) (*model.Player, error) {
	var player model.Player
	err := DB.First(&player, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// プレイヤー情報を更新
func UpdatePlayer(player *model.Player) error {
	return DB.Save(player).Error
}

// プレイヤーを削除
func DeletePlayer(id string) error {
	return DB.Delete(&model.Player{}, "id = ?", id).Error
}

// 一定期間アクセスされていないプレイヤーを削除
func DeleteInactivePlayers(thresholdMinutes int) error {
	threshold := time.Now().Add(-time.Duration(thresholdMinutes) * time.Minute)
	return DB.Where("last_used_at < ?", threshold).Delete(&model.Player{}).Error
}
