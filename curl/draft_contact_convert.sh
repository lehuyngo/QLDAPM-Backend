curl -X POST http://localhost:3002/api/v1/converted-draft-contacts/e5989484-4ff7-4da3-97b4-de64142c7b92 \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNS0xMlQxNzo1Nzo0My40MTgyNzQrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.D-WNe7osPWv9IO8T-Op0pPeqKilM4WHPia5KVlK0C5U" \
    -F "company_logo=@/Users/huynhtrongtien/Downloads/namecard_01.png" \
    -F "fullname=DraftContactsE" \
    -F "phone=55555" \
    -F "email=d@google.com" \
    -F "tags=new-contact:FF00FF,young-staff:AA00CC" \
    -F "client_name=CompanyE" \
    -F "client_website=google.com" \
    -F "client_address=123EEEE" \
    -F "client_tags=new-company:1100FF,technical-company:AA00CC"