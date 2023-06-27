package entities

import "mylife-tools-server/services/store"

var StoreDef = []interface{}{
	store.MakeCollectionBuilder[*Device]("devices", "devices", deviceDecode, deviceEncode),
}
