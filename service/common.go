package service

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/validate"
)

type (
	// Location location
	Location struct {
		IP       string `json:"ip,omitempty" valid:"ip"`
		Country  string `json:"country,omitempty" valid:"runelength(2|30),optional"`
		Province string `json:"province,omitempty" valid:"runelength(2|30),optional"`
		City     string `json:"city,omitempty" valid:"runelength(2|30),optional"`
		ISP      string `json:"isp,omitempty" valid:"runelength(2|30),optional"`
	}
)

// GetLocationByIP get location by ip address
func GetLocationByIP(ip string, c *cod.Context) (l *Location, err error) {
	url := config.GetString("ipLocation")
	d := helper.GetWithContext(url, c)
	d.Param("ip", ip)
	_, body, err := d.Do()
	if err != nil {
		return
	}
	l = new(Location)
	err = validate.Do(l, body)
	if err != nil {
		return
	}
	return
}
