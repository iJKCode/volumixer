package pulse

import (
	"github.com/jfreymuth/pulse/proto"
)

func request[Response any, Request proto.RequestArgs, RespPtr interface {
	proto.Reply
	*Response
}](c *Client, request Request) (*Response, error) {
	response := new(Response)
	err := c.api.Request(request, RespPtr(response))
	return response, err
}

func requestList[Item any, Response ~[]*Item, Request proto.RequestArgs, RespPtr interface {
	proto.Reply
	*Response
}](c *Client, request Request) ([]*Item, error) {
	response := new(Response)
	err := c.api.Request(request, RespPtr(response))
	if err != nil || response == nil {
		return make([]*Item, 0), err
	}
	return *response, err
}

func (c *Client) SetEventCallback(callback func(event any)) {
	c.api.Callback = callback
}

func (c *Client) SetEventSubscription(mask proto.SubscriptionMask) error {
	return c.Command(&proto.Subscribe{
		Mask: mask,
	})
}

func (c *Client) GetClientVersion() proto.Version {
	return c.api.Version()
}

func (c *Client) GetServerInfo() (*proto.GetServerInfoReply, error) {
	return request[proto.GetServerInfoReply](c, &proto.GetServerInfo{})
}

func (c *Client) Command(cmd proto.RequestArgs) error {
	return c.api.Request(cmd, nil)
}
