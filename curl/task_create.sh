curl -X POST http://localhost:3002/api/v1/tasks \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUR0wgU29sdXRpb25zIiwiZXhwIjoiMjAyNC0wNi0xMFQxNToyMDozMi41OTM0MzkrMDc6MDAiLCJ1c2VyX2lkIjo5fQ.OQlR2YyP5XGbG2okEyC22bK9IynjGLO8YpWNykK5v0w" \
    -F "title=1129" \
    -F "status=2" \
    -F "priority=2" \
    -F "label=2" \
    -F "due_date=456" \
    -F "estimated_hours=456" \
    -F "description=456" \
    -F "project_uuid=24c75496-d1f7-46f4-b15f-9a960545cb02" \
    -F "assignee_uuids=8490ffc4-8f9f-4526-831f-4e8e7552eee6,8481bdc4-ec84-4a6f-8291-257a7fe86e21" \
    -F "attach_file_1=@C:/Users/akirakuma5/Desktop/text.txt" \
    -F "attach_file_2=@C:/Users/akirakuma5/Downloads/welcome_img.png"