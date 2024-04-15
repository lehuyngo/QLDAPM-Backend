curl -X POST http://localhost:3002/api/v1/batch-mails \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNi0xMFQxNToyMDozMi41OTM0MzkrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.OQlR2YyP5XGbG2okEyC22bK9IynjGLO8YpWNykK5v0w" \
    -F "subject=TestMail" \
    -F "content=ThisIsTestMailAsCCWithAttachFile" \
    -F "receiver_contact_uuids=43a21a29-350c-421a-bdca-554f4f5717ce,b168a4b2-8169-4e9b-b160-be1c7d5e5f85,5f485fad-82f1-435d-ae23-ad2a8c9f14b1" \
    -F "cc_user_uuids=74935198-a191-4437-a049-a2b04a0015d8" \
    -F "cc_mail_addresses=huynhtrongtien89@gmail.com" \
    -F "attach_file_1=@/Users/huynhtrongtien/Downloads/mail_create.txt"