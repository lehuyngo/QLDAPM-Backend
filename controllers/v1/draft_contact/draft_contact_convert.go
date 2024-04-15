package draft_contact

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ConvertDraftContact(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.ConvertDraftContactRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	uuid := c.Param("uuid")
	draftContact, err := h.DraftContact.ReadByUUID(ctx, uuid)
	if err != nil {
		log.For(c).Error("[convert-draft-contact] query draft contact info failed", log.Field("uuid", uuid), 
			log.Field("user_id", userID), log.Field("draft_contact_uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.GetOrganization().GetID() != draftContact.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	// Add contact
	contact, err := h.addContact(ctx, user, req, draftContact)
	if err != nil {
		log.For(c).Error("[convert-draft-contact] add contact failed", log.Field("uuid", uuid), 
					log.Field("user_id", userID), log.Field("draft_contact_uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Add client
	client, err := h.addClient(c, user, req, draftContact)
	if err != nil {
		log.For(c).Error("[convert-draft-contact] add client failed", log.Field("uuid", uuid), 
					log.Field("user_id", userID), log.Field("draft_contact_uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Update contact id for draft
	err = h.DraftContact.UpdateContact(ctx, draftContact.ID, contact.ID)
	if err != nil {
		log.For(c).Error("[convert-draft-contact] update contact id for draft failed", log.Field("uuid", uuid), 
					log.Field("user_id", userID), log.Field("contact_id", contact.ID), log.Field("draft_contact_uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Add contact for client
	var data []entities.ClientContact
	data = append(data, entities.ClientContact{
		ClientID: client.ID,
		ContactID: contact.ID,
		CreatedBy: userID,
	})
	err = h.ContactClient.AddBatch(ctx, data)
	if err != nil {
		if !define.IsErrDuplicateKey(err) {
			log.For(c).Error("[create-contact-client] create bacth contact failed", log.Field("user_id", userID), 
				log.Field("contact_uuid", contact.UUID), log.Field("client_uuid", client.UUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	// Tags for contact
	err = h.addContactTag(c, user, req.Tags, contact)
	if err != nil {
		log.For(c).Error("[create-contact-client] add tag for contact failed", log.Field("user_id", userID), 
			log.Field("contact_uuid", contact.UUID), log.Field("client_uuid", client.UUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Tags for client
	err = h.addClientTag(c, user, req.ClientTags, client)
	if err != nil {
		log.For(c).Error("[create-contact-client] add tag for client failed", log.Field("user_id", userID), 
			log.Field("contact_uuid", contact.UUID), log.Field("client_uuid", client.UUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.ConvertDraftContactResponse{
		ContactUUID: contact.UUID,
		ClientUUID: client.UUID,
	})
}

func (h Handler) addContact(ctx context.Context, user *entities.User, req *apis.ConvertDraftContactRequest, draftContact *entities.DraftContact) (*entities.Contact, error) {
	contact, err := h.Contact.ReadByEmail(ctx, user.OrganizationID, req.Email)
	// if exist contact, we update
	if err == nil {
		if user.GetOrganization().GetID() == contact.OrganizationID {
			contact.UpdatedBy = user.ID
			contact.FullName = req.FullName
			contact.Phone = req.Phone
			contact.Email = req.Email
			contact.NameCardID = draftContact.NameCardID
			contact.NameCard = &entities.File{
				UUID:				draftContact.GetNameCard().GetUUID(),
				OriginalName:		draftContact.GetNameCard().GetOriginalName(),
				RelativePath:		draftContact.GetNameCard().GetRelativePath(),
				RelativeThumbnail:	draftContact.GetNameCard().GetRelativeThumbnail(),
				Ext:				draftContact.GetNameCard().GetExt(),
				CreatedBy:			draftContact.GetNameCard().CreatedBy,
			}
			err := h.Contact.Update(ctx, contact)
			if err != nil {
				return nil, err
			}

			return contact, nil
		}
	}

	// if not exist contact, we add new
	contact = &entities.Contact{
		Base: &entities.Base{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
		FullName:		req.FullName,
		ShortName:		req.FullName,
		JobTitle:		"",
		Phone:			req.Phone,
		Email:			req.Email,
		OrganizationID: user.OrganizationID,
		Gender: 		entities.Male,
		Birthday: 		0,
		NameCardID: 	draftContact.NameCardID,
		NameCard: 		&entities.File{
			UUID:				draftContact.GetNameCard().GetUUID(),
			OriginalName:		draftContact.GetNameCard().GetOriginalName(),
			RelativePath:		draftContact.GetNameCard().GetRelativePath(),
			RelativeThumbnail:	draftContact.GetNameCard().GetRelativeThumbnail(),
			Ext:				draftContact.GetNameCard().GetExt(),
			CreatedBy:			draftContact.GetNameCard().CreatedBy,
		},
	}
	_, err = h.Contact.Create(ctx, contact)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (h Handler) addContactTag(c *gin.Context, user *entities.User, tagList string, contact *entities.Contact) error {
	ctx := c.Request.Context()
	tags := strings.Split(tagList, ",")
	for _, val := range tags {
		parts := strings.Split(val, ":")
		if len(parts) < 2 {
			continue
		}
		tagName := parts[0]
		tagColor := parts[1]

		tag, err := h.ContactTag.ReadByName(ctx, user.OrganizationID, tagName)
		if err != nil {
			tag = &entities.ContactTag{
				Name: tagName,
				Color: tagColor,
				OrganizationID: user.OrganizationID,
			}
			tag.Contacts = make([]*entities.ContactContactTag, 0)
			tag.Contacts = append(tag.Contacts, &entities.ContactContactTag{
				ContactID: contact.ID,
				Color: tagColor,
				CreatedBy: user.ID,
			})
			_, err = h.ContactTag.Create(ctx, tag)
			if err != nil {
				log.For(c).Error("[create-draft-contact-tag] insert database failed", log.Field("user_id", user.ID), log.Field("contact_uuid", contact.UUID), log.Err(err))
				return err
			}

			continue
		}
		
		// only add relationship
		err = h.ContactTag.Add(ctx, &entities.ContactContactTag{
			ContactID: contact.ID,
			TagID: tag.ID,
			Color: tag.Color,
			CreatedBy: user.ID,
		})
		if define.IsErrDuplicateKey(err) {
			continue
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) addClient(c *gin.Context, user *entities.User, req *apis.ConvertDraftContactRequest, draftContact *entities.DraftContact) (*entities.Client, error) {
	ctx := c.Request.Context()
	logo, errLogo := services.UploadImage(c, "company_logo", 125)
	client, err := h.Client.ReadByWebsite(ctx, user.OrganizationID, req.ClientWebsite)
	// if exist contact, we update
	if err == nil {
		if user.GetOrganization().GetID() == client.OrganizationID {
			client.UpdatedBy = user.ID
			client.FullName = req.ClientName
			client.Address = req.ClientAddress
			client.Website = req.ClientWebsite
			if errLogo == nil {
				client.LogoID = logo.UUID
				client.Logo = &entities.File{
					UUID:				logo.UUID,
					OriginalName:		logo.OriginalName,
					RelativePath:		logo.RelativePathFile,
					RelativeThumbnail:	logo.Thumbnail,
					Ext:				logo.FileExt,
					CreatedBy:			user.ID,
				}
			} else if draftContact.GetCompanyLogo() != nil {
				client.LogoID = draftContact.GetCompanyLogo().GetUUID()
				client.Logo = &entities.File{
					UUID:				draftContact.GetCompanyLogo().GetUUID(),
					OriginalName:		draftContact.GetCompanyLogo().GetOriginalName(),
					RelativePath:		draftContact.GetCompanyLogo().GetRelativePath(),
					RelativeThumbnail:	draftContact.GetCompanyLogo().GetRelativeThumbnail(),
					Ext:				draftContact.GetCompanyLogo().GetExt(),
					CreatedBy:			user.ID,
				}
			}
			err := h.Client.Update(ctx, client)
			if err != nil {
				return nil, err
			}

			return client, nil
		}
	}

	// if not exist contact, we add new
	client = &entities.Client{
		Base: entities.Base{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
		FullName:		req.ClientName,
		Address: 		req.ClientAddress,
		Website:		req.ClientWebsite,
		OrganizationID: user.OrganizationID,
	}
	if errLogo == nil {
		client.LogoID = logo.UUID
		client.Logo = &entities.File{
			UUID:				logo.UUID,
			OriginalName:		logo.OriginalName,
			RelativePath:		logo.RelativePathFile,
			RelativeThumbnail:	logo.Thumbnail,
			Ext:				logo.FileExt,
			CreatedBy:			user.ID,
		}
	}

	_, err = h.Client.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (h Handler) addClientTag(c *gin.Context, user *entities.User, tagList string, client *entities.Client) error {
	ctx := c.Request.Context()
	tags := strings.Split(tagList, ",")
	for _, val := range tags {
		parts := strings.Split(val, ":")
		if len(parts) < 2 {
			continue
		}
		tagName := parts[0]
		tagColor := parts[1]

		tag, err := h.ClientTag.ReadByName(ctx, user.OrganizationID, tagName)
		if err != nil {
			tag = &entities.ClientTag{
				Name: tagName,
				Color: tagColor,
				OrganizationID: user.OrganizationID,
			}
			tag.Clients = make([]*entities.ClientClientTag, 0)
			tag.Clients = append(tag.Clients, &entities.ClientClientTag{
				ClientID: client.ID,
				Color: tagColor,
				CreatedBy: user.ID,
			})

			_, err = h.ClientTag.Create(ctx, tag)
			if err != nil {
				log.For(c).Error("[create-draft-contact-client-tag] insert database failed", log.Field("user_id", user.ID), log.Field("client_uuid", client.UUID), log.Err(err))
				return err
			}

			continue
		}
		
		err = h.ClientTag.Add(ctx, &entities.ClientClientTag{
			ClientID: client.ID,
			TagID: tag.ID,
			Color: tag.Color,
			CreatedBy: user.ID,
		})
		if define.IsErrDuplicateKey(err) {
			continue
		}
		if err != nil {
			return err
		}
	}

	return nil
}