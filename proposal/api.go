// Copyright (c) 2021 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package proposal

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/discovery/gorest"
	"github.com/rs/zerolog/log"
)

type API struct {
	service *Service
}

func NewAPI(service *Service) *API {
	return &API{service: service}
}

// Ping godoc.
// @Summary Ping
// @Description Ping
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func (a *API) Ping(c *gin.Context) {
	c.JSON(200, PingResponse{"pong"})
}

type PingResponse struct {
	Message string `json:"message"`
}

// Proposals list proposals.
// @Summary List proposals
// @Description List proposals
// @Param from query string false "Consumer country"
// @Param provider_id query string false "Provider ID"
// @Param service_type query string false "Service type"
// @Param country query string false "Provider country"
// @Param ip_type query string false "IP type (residential, datacenter, etc.)"
// @Param access_policy query string false "Access policy. When empty, returns only public proposals (default). Use * to return all."
// @Param access_policy_source query string false "Access policy source"
// @Param price_gib_max query number false "Maximum price per GiB. When empty, will not filter by it. Price is set in ethereum wei."
// @Param price_hour_max query number false "Maximum price per hour. When empty, will not filter by it. Price is set in ethereum wei."
// @Param compatibility_min query number false "Minimum compatibility. When empty, will not filter by it."
// @Param compatibility_max query number false "Maximum compatibility. When empty, will not filter by it."
// @Param quality_min query number false "Minimal quality threshold. When empty will be defaulted to 0. Quality ranges from [0.0; 3.0]"
// @Accept json
// @Product json
// @Success 200 {array} v2.Proposal
// @Router /proposals [get]
func (a *API) Proposals(c *gin.Context) {
	opts := ListOpts{
		providerID:         c.Query("provider_id"),
		from:               c.Query("from"),
		serviceType:        c.Query("service_type"),
		country:            c.Query("country"),
		accessPolicy:       c.Query("access_policy"),
		accessPolicySource: c.Query("access_policy_source"),
		ipType:             c.Query("ip_type"),
	}

	priceGiBMax, _ := strconv.ParseInt(c.Query("price_gib_max"), 10, 64)
	opts.priceGiBMax = priceGiBMax

	priceHourMax, _ := strconv.ParseInt(c.Query("price_hour_max"), 10, 64)
	opts.priceHourMax = priceHourMax

	compatibilityMin, _ := strconv.ParseInt(c.Query("compatibility_min"), 10, 16)
	opts.compatibilityMin = int(compatibilityMin)

	compatibilityMax, _ := strconv.ParseInt(c.Query("compatibility_max"), 10, 16)
	opts.compatibilityMax = int(compatibilityMax)

	qlb, _ := strconv.ParseFloat(c.Query("quality_min"), 32)
	opts.qualityMin = float32(qlb)

	proposals, err := a.service.List(opts)

	if err != nil {
		log.Err(err).Msg("Failed to list proposals")
		c.JSON(500, gorest.Err500)
		return
	}

	c.JSON(200, proposals)
}

func (a *API) RegisterRoutes(r gin.IRoutes) {
	r.GET("/ping", a.Ping)
	r.GET("/proposals", a.Proposals)
}
