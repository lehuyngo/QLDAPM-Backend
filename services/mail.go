package services

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/url_parser"
)

type IMail interface {
	Create(ctx context.Context, data *entities.Mail) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Mail, error)
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Mail, error)
}

type Mail struct {
	Model         models.IMail
	ActivityModel models.ContactMailActivity
	TrackedURLModel models.ITrackedURL
}

func NewMail() IMail {
	return &Mail{
		Model:         models.Mail{},
		ActivityModel: models.ContactMailActivity{},
		TrackedURLModel: models.TrackedURL{},
	}
}

func (p *Mail) Create(ctx context.Context, data *entities.Mail) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	// add mail activity for each receiver
	for _, val := range data.Receivers {
		if val.ContactID > 0 {
			p.ActivityModel.Create(ctx, &entities.ContactMailActivity{
				UUID:      uuid.NewString(),
				CreatedBy: data.CreatedBy,
				MailID:    data.ID,
				ContactID: val.ContactID,
				Type:      entities.ActivityCreated,
			})
		}
	}

	return result, nil
}

func (p *Mail) ReadByUUID(ctx context.Context, uuid string) (*entities.Mail, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Mail) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Mail, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func UpgradeURLContent(ctx context.Context, contactID int64, content string, useShortenLink bool) (string, map[string]entities.EncodedURL, error) {
	urls := url_parser.ExtractURL(content)
	sort.Slice(urls, func(i, j int) bool {
		return len(urls[i]) > len(urls[j])
	})

	codes := make(map[string]entities.EncodedURL)
	result := content
	for _, val := range urls {
		trackingCode, err := GenerateToken(contactID, time.Now().UnixMicro())
		if err != nil {
			return "", nil, err
		}

		var newURL string
		if strings.Contains(val, "?") == true {
			newURL = val + "&cid=" + trackingCode
		} else {
			newURL = val + "?cid=" + trackingCode
		}

		shortenLink := Config.Domain.Domain + "d/" + trackingCode
		if useShortenLink == true {
			result = strings.ReplaceAll(result, val, shortenLink)
		} else {
			result = strings.ReplaceAll(result, val, newURL)
		}
		
		codes[trackingCode] = entities.EncodedURL{
			ShortenLink: shortenLink,
			OriginalURL: newURL,
		}
	}

	return result, codes, nil
}
