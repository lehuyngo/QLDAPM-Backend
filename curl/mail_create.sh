curl -X POST http://localhost:3002/api/v1/mails \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNi0yNlQxMDozMzowMi4wNDc1NDErMDc6MDAiLCJ1c2VyX2lkIjo5fQ.W9bvbOPDy3RbGKs0UxWrQXU_3ZxvTY0MY1oiCPlpHYQ" \
    -F "subject=TestMail" \
    -F "content=MailContentWithAttachFiles https://www.google.com/" \
    -F "receiver_contact_uuids=8f09c853-93ef-4b76-a0e2-9de815aee980" \
    -F "receiver_mail_addresses=ngqubao.tgl@gmail.com" \
    -F "cc_mail_addresses=huynhtrongtien89@gmail.com" \
    -F "attach_file_1=@/Users/huynhtrongtien/Downloads/mail_create.txt"