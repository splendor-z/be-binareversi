package db

import (
	"be-binareversi/model"

	"time"
)

func CreateRoom(room *model.Room) error {
	return DB.Create(room).Error
}

func GetRoomByID(id string) (*model.Room, error) {
	var room model.Room
	err := DB.First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func UpdateRoom(room *model.Room) error {
	return DB.Save(room).Error
}

func DeleteRoom(id string) error {
	return DB.Delete(&model.Room{}, "id = ?", id).Error
}

func DeleteOldRooms(thresholdMinutes int) error {
	threshold := time.Now().Add(-time.Duration(thresholdMinutes) * time.Minute)
	return DB.Where("created_at < ?", threshold).Delete(&model.Room{}).Error
}
