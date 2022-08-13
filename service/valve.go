package service

import (
	"context"
	"net"
	"time"

	"github.com/jacobalberty/roofwater/service/client/tasmota"
	"github.com/jacobalberty/roofwater/service/utils"
	"go.opentelemetry.io/otel/codes"
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
	tClient := tasmota.Client{
		Type: tasmota.ClientTypeWeb,
		IP:   v.IP,
	}
	tCommand := tasmota.Command{}.Power(tasmota.PowerOn).Delay(d).Power(tasmota.PowerOff)
	err := tClient.Execute(tCommand)
	if err != nil {
		span.SetStatus(codes.Error, "RWPulse failed")
		span.RecordError(err)
		utils.Logger.Ctx(ctx).Error("Failed to pulse valve",
			zap.String("ip", v.IP.String()),
			zap.Duration("duration", d),
			zap.Error(err),
		)
	}
}
