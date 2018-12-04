package main

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
)

func main() {
	listen := config.GetString("listen")
	d := cod.New()

	d.ListenAndServe(listen)
}
