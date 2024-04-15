package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/account"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/batch_mail"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_attach_file"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_contact"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_contact_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_note"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_note_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_project"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_project_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_tag"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/client_tag_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_attach_file"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_client"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_client_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_mail_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_mail_shortclick"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_note"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_note_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_project"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_project_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_tag"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contact_tag_activity"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/contributor"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/draft_contact"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/draft_contact_report"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/mail"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/meeting"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/meeting_highlight"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/meeting_note"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/project"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/receive_mail_history"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/report_read_mail"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/static_file"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/task"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/task_assignee"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/task_attach_file"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/task_comment"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/url_access"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/controllers/v1/user"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/jwt"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func StartHTTPServer() {

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "https://sky-crm.click", "http://sky-crm.click", "https://crm.tgl-cloud.com", "https://crm-dev.tgl-cloud.com", "https://crm-stag.tgl-cloud.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return (origin == "http://localhost:3000") || (origin == "http://localhost:3001") || (origin == "https://sky-crm.click") || (origin == "http://sky-crm.click") || (origin == "https://crm.tgl-cloud.com") || (origin == "https://crm-dev.tgl-cloud.com") || (origin == "https://crm-stag.tgl-cloud.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	router.LoadHTMLGlob("webapp/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// apis

	accountHandler := account.NewHandler()
	urlAccessHandler := url_access.NewHandler()
	receivedMailHistoryHandler := receive_mail_history.NewHandler()

	router.POST("/api/v1/auth", accountHandler.Authenticate)
	router.POST("/api/v1/register", accountHandler.Register)
	router.GET("/api/v1/me", middlewares.Authenticate(), accountHandler.Me)
	router.POST("/api/v1/customer-access", urlAccessHandler.AccessURL)
	router.GET("/d/:token", urlAccessHandler.RedirectShortenLink)

	clientHandler := client.NewHandler()
	clientActivityHandler := client_activity.NewHandler()
	clientNoteHandler := client_note.NewHandler()
	clientTagHandler := client_tag.NewHandler()
	clientContactHandler := client_contact.NewHandler()
	clientAttachFileHandler := client_attach_file.NewHandler()
	clientProjectHandler := client_project.NewHandler()
	clientProjectActivityHandler := client_project_activity.NewHandler()
	clientTagActivityHandler := client_tag_activity.NewHandler()
	clientNoteActivityHandler := client_note_activity.NewHandler()
	clientContactActivityHandler := client_contact_activity.NewHandler()

	contactHandler := contact.NewHandler()
	contactNoteHandler := contact_note.NewHandler()
	contactTagHandler := contact_tag.NewHandler()
	contactClientHandler := contact_client.NewHandler()
	contactClientActivityHandler := contact_client_activity.NewHandler()
	contactAttachFileHandler := contact_attach_file.NewHandler()
	contactProjectHandler := contact_project.NewHandler()
	contactProjecActivityHandler := contact_project_activity.NewHandler()
	contactTagActivityHandler := contact_tag_activity.NewHandler()
	contactNoteActivityHandler := contact_note_activity.NewHandler()
	contactMailActivityHandler := contact_mail_activity.NewHandler()
	contactActivityHandler := contact_activity.NewHandler()
	contactMailShortClickHandler := contact_mail_shortclick.NewHandler()

	mailHandler := mail.NewHandler()
	batchMailHandler := batch_mail.NewHandler()

	projectHandler := project.NewHandler()

	draftContactHandler := draft_contact.NewHandler()
	draftContactReportHandler := draft_contact_report.NewHandler()

	taskHandler := task.NewHandler()
	taskAssigneeHandler := task_assignee.NewHandler()
	taskAttachFileHandler := task_attach_file.NewHandler()
	taskCommentHandler := task_comment.NewHandler()

	// reports
	reportReadMail := report_read_mail.NewHandler()

	userHandler := user.NewHandler()

	meetingHandler := meeting.NewHandler()
	meetingNoteHandler := meeting_note.NewHandler()
	contributorHandler := contributor.NewHandler()

	meetingHightligtHandler := meeting_highlight.NewHandler()

	staticFileHandler := static_file.NewHandler()

	apiV1 := router.Group("/api/v1", middlewares.Authenticate())
	{
		staticFile := apiV1.Group("/static-files")
		{
			staticFile.POST("", staticFileHandler.CreateFile)
			staticFile.GET("/:uuid/:file_name", staticFileHandler.DownloadFile)
		}
		activeClient := apiV1.Group("/clients")
		{
			activeClient.POST("", clientHandler.CreateClient)
			activeClient.GET("/:uuid", clientHandler.ReadClient)
			activeClient.GET("", clientHandler.ListClient)
			activeClient.PUT("/:uuid", clientHandler.UpdateClient)
			activeClient.DELETE("/:uuid", clientHandler.InTrashClient)

			activeClient.GET("/:uuid/logo/:file_uuid", clientHandler.ReadClientLogo)

			activeClient.POST("/:uuid/notes", clientNoteHandler.CreateClientNote)
			activeClient.GET("/:uuid/notes/:note_uuid", clientNoteHandler.ReadClientNote)
			activeClient.GET("/:uuid/notes", clientNoteHandler.ListClientNote)
			activeClient.PUT("/:uuid/notes/:note_uuid", clientNoteHandler.UpdateClientNote)
			activeClient.DELETE("/:uuid/notes/:note_uuid", clientNoteHandler.DeleteClientNote)

			activeClient.POST("/:uuid/tags", clientTagHandler.CreateClientTag)
			activeClient.DELETE("/:uuid/tags/:tag_uuid", clientTagHandler.DeleteClientTag)

			activeClient.GET("/:uuid/contacts", clientContactHandler.ListClientContact)
			activeClient.POST("/:uuid/contacts", clientContactHandler.CreateClientContact)
			activeClient.DELETE("/:uuid/contacts/:contact_uuid", clientContactHandler.DeleteClientContact)

			activeClient.GET("/:uuid/files", clientAttachFileHandler.ListClientAttachFile)
			activeClient.POST("/:uuid/files", clientAttachFileHandler.CreateClientAttachFile)
			activeClient.DELETE("/:uuid/files/:file_uuid", clientAttachFileHandler.DeleteClientAttachFile)
			activeClient.GET("/:uuid/downloaded-files/:file_uuid/:file_name", clientAttachFileHandler.DownloadClientAttachFile)

			activeClient.GET("/:uuid/projects", clientProjectHandler.ListClientProject)
			activeClient.POST("/:uuid/projects", clientProjectHandler.CreateClientProject)
			activeClient.DELETE("/:uuid/projects/:project_uuid", clientProjectHandler.DeleteClientProject)

			activeClient.GET("/:uuid/client_activities", clientActivityHandler.ListClientActivity)
			activeClient.GET("/:uuid/project_activities", clientProjectActivityHandler.ListClientProjectActivity)
			activeClient.GET("/:uuid/tag_activities", clientTagActivityHandler.ListClientTagActivity)
			activeClient.GET("/:uuid/note_activities", clientNoteActivityHandler.ListClientNoteActivity)
			activeClient.GET("/:uuid/contact_activities", clientContactActivityHandler.ListClientContactActivity)
		}

		clientTag := apiV1.Group("/client-tags")
		{
			clientTag.GET("", clientTagHandler.ListClientTag)
			clientTag.GET("/:uuid", clientTagHandler.ReadClientTag)
		}

		deletedClient := apiV1.Group("/deleted-clients")
		{
			deletedClient.POST("/:uuid", clientHandler.RestoreClient)
			deletedClient.DELETE("/:uuid", clientHandler.DeleteClient)
		}

		activeContact := apiV1.Group("/contacts")
		{
			activeContact.POST("", contactHandler.CreateContact)
			activeContact.GET("/:uuid", contactHandler.ReadContact)
			activeContact.GET("", contactHandler.ListContact)
			activeContact.PUT("/:uuid", contactHandler.UpdateContact)
			activeContact.DELETE("/:uuid", contactHandler.InTrashContact)

			activeContact.GET("/:uuid/avatar/:file_uuid", contactHandler.ReadContactAvatar)
			activeContact.GET("/:uuid/namecard/:file_uuid", contactHandler.ReadContactNameCard)

			activeContact.POST("/:uuid/notes", contactNoteHandler.CreateContactNote)
			activeContact.GET("/:uuid/notes/:note_uuid", contactNoteHandler.ReadContactNote)
			activeContact.GET("/:uuid/notes", contactNoteHandler.ListContactNote)
			activeContact.PUT("/:uuid/notes/:note_uuid", contactNoteHandler.UpdateContactNote)
			activeContact.DELETE("/:uuid/notes/:note_uuid", contactNoteHandler.DeleteContactNote)

			activeContact.POST("/:uuid/tags", contactTagHandler.CreateContactTag)
			activeContact.DELETE("/:uuid/tags/:tag_uuid", contactTagHandler.DeleteContactTag)

			activeContact.GET("/:uuid/clients", contactClientHandler.ListContactClient)
			activeContact.POST("/:uuid/clients", contactClientHandler.CreateContactClient)
			activeContact.DELETE("/:uuid/clients/:client_uuid", contactClientHandler.DeleteContactClient)

			activeContact.GET("/:uuid/files", contactAttachFileHandler.ListContactAttachFile)
			activeContact.POST("/:uuid/files", contactAttachFileHandler.CreateContactAttachFile)
			activeContact.DELETE("/:uuid/files/:file_uuid", contactAttachFileHandler.DeleteContactAttachFile)
			activeContact.GET("/:uuid/downloaded-files/:file_uuid/:file_name", contactAttachFileHandler.DownloadContactAttachFile)

			activeContact.GET("/:uuid/projects", contactProjectHandler.ListContactProject)
			activeContact.POST("/:uuid/projects", contactProjectHandler.CreateContactProject)
			activeContact.DELETE("/:uuid/projects/:project_uuid", contactProjectHandler.DeleteContactProject)

			activeContact.POST("/:uuid/mail-shortclicks", contactMailShortClickHandler.CreateContactMailShortClick)
			activeContact.GET("/:uuid/mail-shortclicks", contactMailShortClickHandler.ListContactMailShortClick)

			activeContact.GET("/:uuid/tag_activities", contactTagActivityHandler.ListContactTagActivity)
			activeContact.GET("/:uuid/note_activities", contactNoteActivityHandler.ListContactNoteActivity)
			activeContact.GET("/:uuid/mail_activities", contactMailActivityHandler.ListContactMailActivity)
			activeContact.GET("/:uuid/project_activities", contactProjecActivityHandler.ListContactProjectActivity)
			activeContact.GET("/:uuid/client_activities", contactClientActivityHandler.ListContactClientActivity)
			activeContact.GET("/:uuid/contact_activities", contactActivityHandler.ListContactActivity)
		}

		contactTag := apiV1.Group("/contact-tags")
		{
			contactTag.GET("", contactTagHandler.ListContactTag)
			contactTag.GET("/:uuid", contactTagHandler.ReadContactTag)
		}

		deletedContact := apiV1.Group("/deleted-contacts")
		{
			deletedContact.POST("/:uuid", contactHandler.RestoreContact)
			deletedContact.DELETE("/:uuid", contactHandler.DeleteContact)
		}

		activeProject := apiV1.Group("/projects")
		{
			activeProject.POST("", projectHandler.CreateProject)
			activeProject.GET("/:uuid", projectHandler.ReadProject)
			activeProject.GET("/:uuid/highlights", meetingHightligtHandler.ListHighlight)
			activeProject.GET("", projectHandler.ListProject)
			activeProject.PUT("/:uuid", projectHandler.UpdateProject)
			activeProject.PUT("/:uuid/status", projectHandler.UpdateProjectStatus)
			activeProject.DELETE("/:uuid", projectHandler.InTrashProject)
			activeProject.POST("/:uuid/meetings", meetingHandler.CreateMeeting)
			activeProject.GET("/:uuid/meetings", meetingHandler.ListMeeting)
			activeProject.GET("/:uuid/meeting-notes", meetingNoteHandler.ListMeetingNote)
			activeProject.GET("/:uuid/tasks", taskHandler.ListTaskByProject)
			activeProject.POST("/:uuid/meetings/:meeting_uuid/meeting-notes", meetingNoteHandler.CreateMeetingNote)
			activeProject.PUT("/:uuid/meeting-notes/:note_uuid", meetingNoteHandler.UpdateMeetingNote)
			activeProject.DELETE("/:uuid/meeting-notes/:note_uuid", meetingNoteHandler.DeleteMeetingNote)
		}
		activeHighlight := apiV1.Group("/meeting-notes")
		{
			activeHighlight.POST("/:uuid/contributors", contributorHandler.CreateContributor)
			activeHighlight.DELETE("/:uuid/contributors/:contributor_uuid", contributorHandler.DeleteContributor)
			activeHighlight.POST("/:uuid/batch-contributors", contributorHandler.CreateContributorBatch)
			activeHighlight.DELETE("/:uuid/batch-contributors", contributorHandler.DeleteContributorBatch)
			activeHighlight.POST("/:uuid/highlights", meetingHightligtHandler.CreateHighlight)
			activeHighlight.DELETE("/:uuid/highlights/:highlight_uuid", meetingHightligtHandler.DeleteHighlight)
			activeHighlight.DELETE("/:uuid/batch-highlights", meetingHightligtHandler.DeleteHighlightBatch)
		}

		deletedProject := apiV1.Group("/deleted-projects")
		{
			deletedProject.POST("/:uuid", projectHandler.RestoreProject)
			deletedProject.DELETE("/:uuid", projectHandler.DeleteProject)
		}

		projectStatus := apiV1.Group("/project-statuses")
		{
			projectStatus.GET("", projectHandler.ListProjectStatus)
		}

		draftContact := apiV1.Group("/draft-contacts")
		{
			draftContact.POST("", draftContactHandler.CreateDraftContact)
			draftContact.GET("/:uuid", draftContactHandler.ReadDraftContact)
			draftContact.GET("", draftContactHandler.ListDraftContact)
			draftContact.DELETE("/:uuid", draftContactHandler.DeleteDraftContact)
			draftContact.GET("/:uuid/name-card/:file_uuid", draftContactHandler.DownloadDraftContactNameCard)
			draftContact.GET("/:uuid/company-logo/:file_uuid", draftContactHandler.DownloadDraftContactCompanyLogo)
		}
		draftContactReport := apiV1.Group("/draft-contact-report")
		{
			draftContactReport.GET("", draftContactReportHandler.GetDraftContactReport)

		}

		restoredDraftContact := apiV1.Group("/restored-draft-contacts")
		{
			restoredDraftContact.POST("", draftContactHandler.RestoreDraftContact)
		}

		convertedDraftContact := apiV1.Group("/converted-draft-contacts")
		{
			convertedDraftContact.POST("/:uuid", draftContactHandler.ConvertDraftContact)
		}

		task := apiV1.Group("/tasks")
		{
			task.POST("", taskHandler.CreateTask)
			task.GET("", taskHandler.ListTask)
			task.GET("/:uuid", taskHandler.ReadTask)
			task.PUT("/:uuid", taskHandler.UpdateTask)
			task.PUT("/:uuid/status", taskHandler.UpdateTaskStatus)
			task.DELETE("/:uuid", taskHandler.DeleteTask)

			task.POST("/:uuid/assignees", taskAssigneeHandler.CreateTaskAssignee)
			task.DELETE("/:uuid/assignees/:assignee_uuid", taskAssigneeHandler.DeleteTaskAssignee)

			task.POST("/:uuid/files", taskAttachFileHandler.CreateTaskAttachFile)
			task.DELETE("/:uuid/files/:file_uuid", taskAttachFileHandler.DeleteTaskAttachFile)
			task.GET("/:uuid/downloaded-files/:file_uuid/:file_name", taskAttachFileHandler.DownloadTaskAttachFile)

			task.POST("/:uuid/comments", taskCommentHandler.CreateTaskComment)
			task.GET("/:uuid/comments", taskCommentHandler.ListTaskComment)
		}

		mail := apiV1.Group("/mails")
		{
			mail.POST("", mailHandler.CreateMail)
			mail.GET("/:uuid", mailHandler.ReadMail)
		}

		receivedMailHistory := apiV1.Group("/received-mail-histories")
		{
			receivedMailHistory.GET("", receivedMailHistoryHandler.ListReceiveMailHistory)
		}

		report := apiV1.Group("/reports")
		{
			report.POST("/read-mail-timeline", reportReadMail.ReportTimelineReadMail)
		}

		batchMail := apiV1.Group("/batch-mails")
		{
			batchMail.POST("", batchMailHandler.CreateBatchMail)
			batchMail.GET("", batchMailHandler.ListBatchMail)
			batchMail.GET("/:uuid", batchMailHandler.ReadBatchMail)
		}

		user := apiV1.Group("/users")
		{
			user.GET("", userHandler.ListUser)
		}

	}

	host := fmt.Sprintf("localhost:%d", viper.GetInt("service.port"))

	log.Bg().Info("[start service] start service", log.Field("host", host))

	router.Run(host)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func main() {
	services.InitViper()
	log.Setup()

	services.InitServices()

	clients.MySQLClient, _ = clients.NewMySQLClient()
	clients.AutoMigrate()

	middlewares.InitMiddlewares()
	jwt.InitModels()

	StartHTTPServer()
}
