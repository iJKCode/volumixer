package widget

import (
	widgetv1 "github.com/ijkcode/volumixer-api/gen/go/widget/v1"
	"github.com/ijkcode/volumixer/pkg/core/component"
)

type InfoComponent struct {
	Name string
}

func init() {
	component.RegisterFatal(InfoComponentMarshaller{})
}

type InfoComponentMarshaller struct{}

var _ component.Marshaler[InfoComponent, *widgetv1.InfoComponent] = (*InfoComponentMarshaller)(nil)

func (m InfoComponentMarshaller) MarshalProto(cmp InfoComponent) (*widgetv1.InfoComponent, error) {
	return &widgetv1.InfoComponent{
		Name: cmp.Name,
	}, nil
}

func (m InfoComponentMarshaller) UnmarshalProto(msg *widgetv1.InfoComponent) (InfoComponent, error) {
	return InfoComponent{
		Name: msg.Name,
	}, nil
}
