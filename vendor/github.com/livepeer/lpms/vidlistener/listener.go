package vidlistener

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/livepeer/lpms/segmenter"
	"github.com/livepeer/lpms/stream"
	joy4rtmp "github.com/nareix/joy4/format/rtmp"
)

var segOptions = segmenter.SegmenterOptions{SegLength: time.Second * 2}

type LocalStream struct {
	StreamID  string
	Timestamp int64
}

type VidListener struct {
	RtmpServer *joy4rtmp.Server
	FfmpegPath string
}

func (self *VidListener) HandleRTMPPublish(
	makeStreamID func(url *url.URL) (strmID string),
	gotStream func(url *url.URL, rtmpStrm stream.RTMPVideoStream) error,
	endStream func(url *url.URL, rtmpStrm stream.RTMPVideoStream) error) {

	self.RtmpServer.HandlePublish = func(conn *joy4rtmp.Conn) {
		glog.V(2).Infof("RTMP server got upstream: %v", conn.URL)

		s := stream.NewBasicRTMPVideoStream(makeStreamID(conn.URL))
		ctx, cancel := context.WithCancel(context.Background())
		ec := make(chan error)
		go func() { ec <- s.WriteRTMPToStream(ctx, conn) }()

		err := gotStream(conn.URL, s)
		if err != nil {
			glog.Errorf("Error RTMP gotStream handler: %v", err)
			cancel()
			return
		}

		select {
		case err := <-ec:
			endStream(conn.URL, s)
			if err != io.EOF {
				glog.Errorf("Got error writing RTMP: %v", err)
			}
			cancel()
		}
	}
}
