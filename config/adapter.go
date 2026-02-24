package config

type Adapter struct {
	StreamProvideSave string `mapstructure:"stream-provide-save" json:"stream-provide-save" yaml:"stream-provide-save"`
	StreamTsSave      string `mapstructure:"stream-ts-save" json:"stream-ts-save" yaml:"stream-ts-save"`
	CutterVideoUpload string `mapstructure:"cutter-video-upload" json:"cutter-video-upload" yaml:"cutter-video-upload"`
}
