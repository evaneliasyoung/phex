package phaser

type Frame struct {
	FileName         string `json:"filename"`
	Frame            Rect   `json:"frame"`
	Rotated          bool   `json:"rotated"`
	SourceSize       Size   `json:"sourceSize"`
	SpriteSourceSize Rect   `json:"spriteSourceSize"`
	Trimmed          bool   `json:"trimmed"`
}

type Texture struct {
	Format string  `json:"format"`
	Frames []Frame `json:"frames"`
	Image  string  `json:"image"`
	Scale  float64 `json:"scale"`
	Size   Size    `json:"size"`
}

type Atlas struct {
	Meta     map[string]string `json:"meta"`
	Textures []Texture         `json:"textures"`
}
