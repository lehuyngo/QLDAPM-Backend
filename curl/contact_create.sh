curl -X POST http://localhost:3002/api/v1/contacts \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNS0xMlQxNzo1Nzo0My40MTgyNzQrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.D-WNe7osPWv9IO8T-Op0pPeqKilM4WHPia5KVlK0C5U" \
    -F "namecard=@/Users/huynhtrongtien/Downloads/logo_1.jpeg" \
    -F "avatar=@/Users/huynhtrongtien/Downloads/logo_1.jpeg" \
    -F "fullname=Fullname20" \
    -F "shortname=Shortname20" \
    -F "phone=Phone20" \
    -F "email=Email20" \
    -F "job_title=SE20" \
    -F "gender=0" \
    -F "client_uuids=["abc","def","ghk"]" \
    -F "birth_day=1700575740000"