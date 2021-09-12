
# Build application
GOOS=linux GOARCH=amd64 go build

# Move to /root/go
mv Oracle-Hackathon-BE /root/go

# Restart the application service
systemctl restart Go-Backend.service
systemctl enable Go-Backend.service

