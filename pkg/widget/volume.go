package widget

type VolumeComponent struct {
	Level float32
	Muted bool
}

type VolumeChangeCommand struct {
	Level float32
}

type VolumeMuteCommand struct {
	Mute bool
}
