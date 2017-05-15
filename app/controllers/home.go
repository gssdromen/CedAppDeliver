package controllers

import "github.com/revel/revel"

// Home ...
type Home struct {
	GormController
}

// Index ...
func (c Home) Index() revel.Result {
	var params = make(map[string]string)

	params["cerAddr"] = "public/myCer/WuYinJun-CA.crt"

	return c.Render(params)
}
