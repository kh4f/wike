package config

const (
	VK_CAPITAL     = 0x14
	VK_F13         = 0x7C
	VK_VOLUME_DOWN = 0xAE
	VK_VOLUME_UP   = 0xAF
)

var VKCodeMap = map[string]uint16{
	"VK_F13":         VK_F13,
	"VK_CAPITAL":     VK_CAPITAL,
	"VK_VOLUME_UP":   VK_VOLUME_UP,
	"VK_VOLUME_DOWN": VK_VOLUME_DOWN,
}

var RevVKCodeMap = func() map[uint16]string {
	m := make(map[uint16]string)
	for k, v := range VKCodeMap {
		m[v] = k
	}
	return m
}()
