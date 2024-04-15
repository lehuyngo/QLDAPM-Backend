package report_read_mail

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ReportTimelineReadMail(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.ReportTimelineReadMailRequest{}
	userID, _, _ := middlewares.ParseToken(c)
	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(req.TimeRanges) < 1 {
		log.For(c).Error("[report-timeline-read-mail] empty ranges", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, err)
		return
	}
	sort.Slice(req.TimeRanges, func(i, j int) bool {
		return req.TimeRanges[i].StartTime < req.TimeRanges[j].StartTime
	})

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[report-timeline-read-mail] query user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ReportTimelineReadMailResponse{}
	for _, val := range req.TimeRanges {
		resp.Data = append(resp.Data, &apis.ReportTimelineReadMailData{
			TimeRange: val,
		})
	}
	
	// find min max time
	minTime := req.TimeRanges[0].StartTime
	maxTime := req.TimeRanges[0].EndTime
	for _, val := range req.TimeRanges {
		if minTime > val.StartTime {
			minTime = val.StartTime
		}

		if maxTime < val.EndTime {
			maxTime = val.EndTime
		}
	}

	// mustnot read numbers form config, because of self-protect
	countQuery := 0
	maxCountQuery := 50
	recordPerQuery := 1000
	for countQuery < maxCountQuery {
		events, err := h.EventClickURL.List(ctx, user.OrganizationID, minTime, maxTime, countQuery * recordPerQuery, (countQuery + 1) * recordPerQuery - 1)
		if err != nil {
			log.For(c).Error("[report-shorten-link-timeline] query events failed", log.Field("user_id", userID), log.Field("count", countQuery), log.Field("record_per_query", recordPerQuery), log.Err(err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if len(events) < 1 {
			break
		}

		for _, val := range events {
			for i := range resp.Data {
				if (resp.Data[i].TimeRange.StartTime <= val.ReadTime) && (resp.Data[i].TimeRange.EndTime >= val.ReadTime) {
					if val.GetReceiver() == nil {
						continue
					}

					record := &apis.ReadMailRecord{
						Sender: 	&apis.User{
							UUID: val.GetSender().GetUUID(),
							DisplayName: val.GetSender().GetDisplayName(),
						},
						Receiver: 	&apis.Contact{
							UUID: val.GetReceiver().GetUUID(),
							FullName: val.GetReceiver().GetFullName(),
						},
						SendTime: val.SendTime,
						ReadTime: val.ReadTime,
						URL: val.OriginalURL,
					}

					if val.MailID > 0 {
						record.Mail = &apis.Mail{
							UUID: val.MailUUID,
							Subject: val.MailSubject,
						}
					}

					if val.BatchMailID > 0 {
						record.BatchMail = &apis.BatchMail{
							UUID: val.BatchMailUUID,
							Subject: val.BatchMailSubject,
						}
					}

					resp.Data[i].Records = append(resp.Data[i].Records, record)
				}
			}
		}
		countQuery = countQuery + 1
	}

	if countQuery >= maxCountQuery {
		log.For(c).Error("[report-shorten-link-timeline] out range data", log.Field("user_id", userID), log.Field("count", countQuery), log.Field("record_per_query", recordPerQuery))
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].TimeRange.StartTime < resp.Data[j].TimeRange.StartTime
	})

	log.For(c).Info("[report-shorten-link-timeline] process success", log.Field("user_id", userID))
	c.JSON(http.StatusOK, resp)
}
