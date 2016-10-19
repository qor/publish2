package version

type Visible struct {
	PublishReady bool
}

func (visible *Visible) SetPublishReady(b bool) {
	visible.PublishReady = b
}

type PublishReadyInterface interface {
	SetPublishReady(bool)
}
