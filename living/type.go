package living

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type NopLivingType struct{}

func (n NopLivingType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	l := &Living{
		livingData: data.Data.(*livingData),
		tx:         tx,
		handle:     handle,
		data:       data,
	}

	return l
}

func (NopLivingType) EncodeEntity() string {
	panic("implement me")
}

func (NopLivingType) BBox(e world.Entity) cube.BBox {
	return cube.BBox{}
}

func (NopLivingType) DecodeNBT(m map[string]any, data *world.EntityData) {
	data.Data = m
}

func (NopLivingType) EncodeNBT(data *world.EntityData) map[string]any {
	return map[string]any{
		"Pos":     []float32{float32(data.Pos[0]), float32(data.Pos[1]), float32(data.Pos[2])},
		"Motion":  []float32{float32(data.Vel[0]), float32(data.Vel[1]), float32(data.Vel[2])},
		"Yaw":     float32(data.Rot[0]),
		"Pitch":   float32(data.Rot[1]),
		"NameTag": data.Name,
	}
}
