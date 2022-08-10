package service

import (
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

func (v Valve) RWPulse(d time.Duration) {
	utils.Logger.Info("Pulsing valve",
		zap.String("ip", v.IP.String()),
		zap.Duration("duration", d),
	)

	http.Get(fmt.Sprintf("http://%s/cm?cmnd=Power%%20On", v.IP.String()))
	defer http.Get(fmt.Sprintf("http://%s/cm?cmnd=Power%%20Off", v.IP.String()))
	time.Sleep(d)
}
