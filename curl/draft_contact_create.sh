curl -X POST https://sky-crm.click/api/v1/draft-contacts \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNi0wNlQyMjo1NjoyNC42NTYxMjQwMjdaIiwidXNlcl9pZCI6Nn0.yf-60aPd_ZCH0OwbFJ5IIdVPKXf2KjdCjcNpin_4WRc" \
    -F "name_card=@/Users/huynhtrongtien/Downloads/namecard/name_card_5.jpeg" \
    -F "company_logo=@/Users/huynhtrongtien/Downloads/namecard/logo_05.png" \
    -F "fullname=DraftContactsE" \
    -F "phone=55555" \
    -F "email=d@google.com" \
    -F "client_name=CompanyE" \
    -F "client_website=google.com" \
    -F "client_address=123EEEE"