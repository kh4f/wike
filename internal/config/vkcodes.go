package config

var VKCodeMap = map[string]uint16{
	"VK_F13":         0x7C,
	"VK_CAPITAL":     0x14,
	"VK_VOLUME_UP":   0xAF,
	"VK_VOLUME_DOWN": 0xAE,
}

var RevVKCodeMap = func() map[uint16]string {
	m := make(map[uint16]string)
	for k, v := range VKCodeMap {
		m[v] = k
	}
	return m
}()
