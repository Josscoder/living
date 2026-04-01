package living

import (
	"fmt"
	"time"

	"github.com/bedrock-gophers/nbtconv"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
)

type NopLivingType struct{}

func (n NopLivingType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	ld := normalizeLivingData(handle, data)

	l := &Living{
		livingData: ld,
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
	if original, ok := data.Data.(map[string]any); ok {
		for k, v := range original {
			nbt[k] = v
		}
	}

	nbt["Pos"] = nbtconv.Vec3ToFloat32Slice(data.Pos)

	nbt["Motion"] = nbtconv.Vec3ToFloat32Slice(data.Vel)

	nbt["Yaw"] = float32(data.Rot[0])
	nbt["Pitch"] = float32(data.Rot[1])
	nbt["NameTag"] = data.Name

	if ld, ok := data.Data.(*livingData); ok {
		nbt["Fire"] = int16(ld.fireTicks)
		nbt["Age"] = int16(ld.age / (time.Second * 20))
		nbt["Variant"] = ld.variant
		nbt["MarkVariant"] = ld.markVariant
	} else {
		nbt["Fire"] = int16(data.FireDuration.Seconds() * 20)
		nbt["Age"] = int16(data.Age / (time.Second * 20))
	}

	return nbt
}

func normalizeLivingData(handle *world.EntityHandle, data *world.EntityData) *livingData {
	ld := defaultLivingData(handle.Type())

	switch original := data.Data.(type) {
	case *livingData:
		ld = original
		if ld.entityType == nil {
			ld.entityType = handle.Type()
		}
	case map[string]any:
		ageTicks := nbtconv.Int64(original, "Age")
		if ageTicks == 0 {
			ageTicks = int64(nbtconv.Int16(original, "Age"))
		}
		if ageTicks == 0 {
			ageTicks = int64(nbtconv.Int32(original, "Age"))
		}
		ld.age = time.Duration(ageTicks) * time.Second / 20

		fireTicks := nbtconv.Int64(original, "Fire")
		if fireTicks == 0 {
			fireTicks = int64(nbtconv.Int16(original, "Fire"))
		}
		if fireTicks == 0 {
			fireTicks = int64(nbtconv.Int32(original, "Fire"))
		}
		ld.fireTicks = fireTicks

		ld.variant = nbtconv.Int32(original, "Variant")
		ld.markVariant = nbtconv.Int32(original, "MarkVariant")
	case nil:
		// Keep defaults.
	default:
		fmt.Printf("living: tipo inesperado en EntityData.Data (%T), usando defaults\n", original)
	}

	if ld.mc == nil {
		ld.mc = &entity.MovementComputer{}
	}
	if ld.HealthManager == nil {
		ld.HealthManager = entity.NewHealthManager(20, 20)
	}
	if ld.effects == nil {
		ld.effects = make(map[effect.Type]effect.Effect)
	}
	if ld.handler == nil {
		ld.handler = NopHandler{}
	}
	if ld.drops == nil {
		ld.drops = func(func(Drop) bool) {}
	}
	if ld.scale == 0 {
		ld.scale = 1
	}

	data.Data = ld
	return ld
}

func defaultLivingData(entityType world.EntityType) *livingData {
	return &livingData{
		entityType:    entityType,
		mc:            &entity.MovementComputer{},
		HealthManager: entity.NewHealthManager(20, 20),
		effects:       make(map[effect.Type]effect.Effect),
		handler:       NopHandler{},
		drops:         func(func(Drop) bool) {},
		scale:         1,
	}
}
