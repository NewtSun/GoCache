/**
 * @Author : NewtSun
 * @Date : 2023/4/17 16:23
 * @Description :
 **/

package controller

import pb "GoCache/cachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
