package interfaces

type ActionBase interface {
	FeatureUpload() error
	FeatureDownload() error
	FeatureRemove() error
	FeatureClear() error
	GetActionType() string
}
