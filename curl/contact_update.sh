curl -X PUT http://localhost:3002/api/v1/contacts/8f09c853-93ef-4b76-a0e2-9de815aee980 \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNS0xMlQxNzo1Nzo0My40MTgyNzQrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.D-WNe7osPWv9IO8T-Op0pPeqKilM4WHPia5KVlK0C5U" \
    -F "image=@/Users/huynhtrongtien/Downloads/logo_1.jpeg" \
    -F "fullname=Fullname1_update" \
    -F "shortname=Shortname1_update" \
    -F "phone=Phone1_update" \
    -F "email=Email1_update" \
    -F "job_title=CompanySize1_update" \
    -F "gender=0" \
    -F "birth_day=1700576026"