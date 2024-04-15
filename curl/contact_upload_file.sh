curl -X POST http://localhost:3002/api/v1/contacts/5f485fad-82f1-435d-ae23-ad2a8c9f14b1/files \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNS0xMlQxNzo1Nzo0My40MTgyNzQrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.D-WNe7osPWv9IO8T-Op0pPeqKilM4WHPia5KVlK0C5U" \
    -F "file=@/Users/huynhtrongtien/Downloads/logo_1.jpeg"