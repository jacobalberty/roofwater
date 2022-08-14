package service

import (
	"context"
	"time"

	"github.com/jacobalberty/roofwater/service/client/tasmota"
	"github.com/jacobalberty/roofwater/service/utils"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Valve struct {
	Addr       string
	MQTTConfig *MQTTConfig
	tClient    *tasmota.Client
}

func (v *Valve) RWPulse(ctx context.Context, d time.Duration) {
	ctx, span := utils.Tracer.Start(ctx, "RWPulse")
	defer span.End()
	utils.Logger.Ctx(ctx).Info("Pulsing valve",
		zap.String("addr", v.Addr),
		zap.Duration("duration", d),
	)
	if v.tClient == nil {
		if v.MQTTConfig == nil {
			v.tClient = &tasmota.Client{
				Type: tasmota.ClientTypeWeb,
				Addr: v.Addr,
			}
		} else {
			v.tClient = &tasmota.Client{
				Type: tasmota.ClientTypeMQTT,
				MQTTConfig: tasmota.MQTTConfig{
					BrokerUrl:         v.MQTTConfig.URL,
					Topic:             v.Addr,
					Username:          v.MQTTConfig.User,
					Password:          []byte(v.MQTTConfig.Pass),
					ClientID:          "roofwater",
					KeepAlive:         v.MQTTConfig.KeepAlive,
					ConnectRetryDelay: v.MQTTConfig.ConnectRetryDelay,
					Timeout:           v.MQTTConfig.Timeout,
				},
			}
		}
	}
	tCommand := v.tClient.Command().Power(tasmota.PowerOn).Delay(d).Power(tasmota.PowerOff)
	err := tCommand.Execute(ctx, tCommand)
	if err != nil {
		span.SetStatus(codes.Error, "RWPulse failed")
		span.RecordError(err)
		utils.Logger.Ctx(ctx).Error("Failed to pulse valve",
			zap.String("addr", v.Addr),
			zap.Duration("duration", d),
			zap.Error(err),
		)
	}
}
