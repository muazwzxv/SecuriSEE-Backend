
# Move to /root/go
mv Oracle-Hackathon-BE /root/go

# Restart the application service
systemctl restart Go-Backend.service
systemctl enable Go-Backend.service

# Reload nginx
nginx -t && nginx -s reload