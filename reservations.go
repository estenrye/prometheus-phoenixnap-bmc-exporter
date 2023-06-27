package exporter

import (
	"fmt"
	"time"

	"github.com/phoenixnap/go-sdk-bmc/billingapi"
	log "github.com/sirupsen/logrus"
)

type ReservationStatValue struct {
	ProductCode           string
	LocationCode          string
	IsAutoRenewed         bool
	IsActive              bool
	IsAssigned            bool
	IsRenewingInSevenDays bool
	IsExpiringInSevenDays bool

	MonthlyCostTotal float64
	Total            float64
}

func (r *ReservationStatValue) Init(reservation billingapi.Reservation) {
	r.ProductCode = reservation.GetProductCode()
	r.LocationCode = string(reservation.GetLocation())
	r.IsAutoRenewed = reservation.GetAutoRenew()
	r.IsAssigned = reservation.HasAssignedResourceId()

	if tP, vset := reservation.GetEndDateTimeOk(); vset {
		t := *tP
		remainingDuration := t.UTC().Sub(time.Now().UTC()).Seconds()
		r.IsActive = remainingDuration >= 0
		r.IsExpiringInSevenDays = remainingDuration > 0 && remainingDuration <= 7*24*60*60
	} else {
		r.IsActive = true
		r.IsExpiringInSevenDays = false
	}

	if tP, vset := reservation.GetNextRenewalDateTimeOk(); vset {
		t := *tP
		remainingDuration := t.UTC().Sub(time.Now().UTC()).Seconds()
		r.IsRenewingInSevenDays = remainingDuration > 0 && remainingDuration <= 7*24*60*60
	} else {
		r.IsRenewingInSevenDays = false
	}

	r.MonthlyCostTotal = float64(reservation.GetPrice())
	r.Total = 1
}

func (r ReservationStatValue) GetKey() string {
	return fmt.Sprintf(
		"%s-%s-%v-%v-%v",
		r.GetProductCodeLabel(),
		r.GetLocationLabel(),
		r.GetStatusLabel(),
		r.GetAssignmentLabel(),
		r.GetAutoRenewalLabel(),
	)
}

func (r ReservationStatValue) GetProductCodeLabel() string {
	return r.ProductCode
}

func (r ReservationStatValue) GetLocationLabel() string {
	return r.LocationCode
}

func (r ReservationStatValue) GetAutoRenewalLabel() string {
	if r.IsAutoRenewed {
		return "enabled"
	}
	return "disabled"
}

func (r ReservationStatValue) GetAssignmentLabel() string {
	if r.IsAssigned {
		return "assigned"
	}
	return "unassigned"
}

func (r ReservationStatValue) GetStatusLabel() string {
	if r.IsActive {
		if r.IsExpiringInSevenDays {
			return "expiring"
		}
		if r.IsRenewingInSevenDays {
			return "renewing"
		}
		return "active"
	}
	return "expired"
}

func (r *ReservationStatValue) Increment(reservation billingapi.Reservation) {
	r.MonthlyCostTotal += float64(reservation.GetPrice())
	r.Total += 1
}

type ReservationStats struct {
	TotalReservations map[string]ReservationStatValue
}

func GetBmcReservations(config BmcApiConfiguration) ([]ReservationStats, error) {
	apiClient := getBillingApiClient(config.ToClientCredentials())

	var stats []ReservationStats

	var res ReservationStats
	res.TotalReservations = make(map[string]ReservationStatValue)

	resp, r, err := apiClient.ReservationsApi.ReservationsGet(getContext()).ProductCategory(billingapi.SERVER).Execute()
	if err != nil {
		log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `ReservationsApi.ReservationsGet`.")
		return stats, err
	}

	for _, reservation := range resp {
		if reservation.GetProductCode() == "bandwidth" {
			continue
		}

		var v ReservationStatValue
		v.Init(reservation)
		if val, ok := res.TotalReservations[v.GetKey()]; ok {
			val.Increment(reservation)
			res.TotalReservations[v.GetKey()] = val
		} else {
			res.TotalReservations[v.GetKey()] = v
		}
	}

	stats = append(stats, res)
	log.WithField("reservations", stats).Info("Get Reservation Metrics")

	return stats, nil
}
