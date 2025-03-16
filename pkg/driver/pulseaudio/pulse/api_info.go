package pulse

import "github.com/jfreymuth/pulse/proto"

func (c *Client) GetSinkInfoByIndex(index uint32) (*proto.GetSinkInfoReply, error) {
	return request[proto.GetSinkInfoReply](c, &proto.GetSinkInfo{SinkIndex: index})
}

func (c *Client) GetSinkInfoByName(name string) (*proto.GetSinkInfoReply, error) {
	return request[proto.GetSinkInfoReply](c, &proto.GetSinkInfo{SinkName: name})
}

func (c *Client) GetSinkInfoList() ([]*proto.GetSinkInfoReply, error) {
	return requestList[proto.GetSinkInfoReply, proto.GetSinkInfoListReply](c, &proto.GetSinkInfoList{})
}

func (c *Client) GetSourceInfoByIndex(index uint32) (*proto.GetSourceInfoReply, error) {
	return request[proto.GetSourceInfoReply](c, &proto.GetSourceInfo{SourceIndex: index})
}

func (c *Client) GetSourceInfoByName(name string) (*proto.GetSourceInfoReply, error) {
	return request[proto.GetSourceInfoReply](c, &proto.GetSourceInfo{SourceName: name})
}

func (c *Client) GetSourceInfoList() ([]*proto.GetSourceInfoReply, error) {
	return requestList[proto.GetSourceInfoReply, proto.GetSourceInfoListReply](c, &proto.GetSourceInfoList{})
}

func (c *Client) GetSinkInputInfoByIndex(index uint32) (*proto.GetSinkInputInfoReply, error) {
	return request[proto.GetSinkInputInfoReply](c, &proto.GetSinkInputInfo{SinkInputIndex: index})
}

func (c *Client) GetSinkInputInfoList() ([]*proto.GetSinkInputInfoReply, error) {
	return requestList[proto.GetSinkInputInfoReply, proto.GetSinkInputInfoListReply](c, &proto.GetSinkInputInfoList{})
}

func (c *Client) GetSourceOutputInfoByIndex(index uint32) (*proto.GetSourceOutputInfoReply, error) {
	return request[proto.GetSourceOutputInfoReply](c, &proto.GetSourceOutputInfo{SourceOutpuIndex: index})
}

func (c *Client) GetSourceOutputInfoList() ([]*proto.GetSourceOutputInfoReply, error) {
	return requestList[proto.GetSourceOutputInfoReply, proto.GetSourceOutputInfoListReply](c, &proto.GetSourceOutputInfoList{})
}
