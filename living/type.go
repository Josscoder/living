package living

import (
	"time"

	"github.com/bedrock-gophers/nbtconv"
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
	nbt := make(map[string]any)

	nbt["Pos"] = nbtconv.Vec3ToFloat32Slice(data.Pos)

	nbt["Motion"] = nbtconv.Vec3ToFloat32Slice(data.Vel)

	nbt["Yaw"] = float32(data.Rot[0])
	nbt["Pitch"] = float32(data.Rot[1])
	nbt["NameTag"] = data.Name
	nbt["Fire"] = int16(data.FireDuration.Seconds() * 20)
	nbt["Age"] = int16(data.Age / (time.Second * 20))

	if original, ok := data.Data.(map[string]any); ok {
		for k, v := range original {
			nbt[k] = v
			println("k = %s, v = %s", k, v)
		}
	}

	return nbt
}
