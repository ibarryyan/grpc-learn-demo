package roundrobin

import (
	"sync/atomic"
	
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/grpcrand"
)

const Name = "barry_lb"

var logger = grpclog.Component("barry_lb")

func NewBarryBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &barryPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(NewBarryBuilder())
}

type barryPickerBuilder struct {
}

func (b *barryPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("barryPicker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &barryPicker{
		subConn: scs,
		next:    int32(grpcrand.Intn(len(scs))),
	}
}

type barryPicker struct {
	subConn []balancer.SubConn
	next    int32
}

func (p *barryPicker) Pick(opts balancer.PickInfo) (balancer.PickResult, error) {
	var bs balancer.SubConn
	// 下一次请求的索引
	nextIndex := atomic.AddInt32(&p.next, 1)
	// 可以提供的连接数
	subConnLen := len(p.subConn)

	if nextIndex%2 == 0 && subConnLen == 2 {
		bs = p.subConn[subConnLen-1]
	} else {
		bs = p.subConn[0]
	}

	return balancer.PickResult{
		SubConn: bs,
	}, nil
}
