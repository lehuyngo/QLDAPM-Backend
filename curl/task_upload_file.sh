curl -X POST http://localhost:3002/api/v1/tasks/5f42a4ee-0aef-47c8-8e82-6e6d2d87f157/files \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNS0xMlQxNzo1Nzo0My40MTgyNzQrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.D-WNe7osPWv9IO8T-Op0pPeqKilM4WHPia5KVlK0C5U" \
    -F "file=@C:/Users/akirakuma5/Downloads/a.png"