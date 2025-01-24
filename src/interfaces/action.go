package interfaces

type ServiceAction interface {
	FeatureUpload() error
	FeatureDownload() error
	FeatureRemove() error
	FeatureClear() error
	GetActionType() string
}
