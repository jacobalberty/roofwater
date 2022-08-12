package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jacobalberty/roofwater/service/utils"
	"go.uber.org/zap"
)

type Valve struct {
	IP net.IP
}

func (v Valve) RWPulse(ctx context.Context, d time.Duration) {
	ctx, span := utils.Tracer.Start(ctx, "RWPulse")
	defer span.End()
	utils.Logger.Ctx(ctx).Info("Pulsing valve",
		zap.String("ip", v.IP.String()),
		zap.Duration("duration", d),
	)
	_, err := http.Get(fmt.Sprintf("http://%s/cm?cmnd=Power%%20On", v.IP.String()))
	if err != nil {
		utils.Logger.Ctx(ctx).Error("Failed to pulse valve",
			zap.Error(err),
		)
		return
	}
	defer func() {
		_, err := http.Get(fmt.Sprintf("http://%s/cm?cmnd=Power%%20Off", v.IP.String()))
		if err != nil {
			utils.Logger.Ctx(ctx).Error("Failed to turn off valve",
				zap.Error(err),
			)
		}
	}()
	time.Sleep(d)
}
