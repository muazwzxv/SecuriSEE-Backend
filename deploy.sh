
# Move to /root/go
mv Oracle-Hackathon-BE /root/go
mv config.yml /root/go

# Restart the application service
#systemctl restart Go-Backend.service
#systemctl enable Go-Backend.service
systemctl daemon-reload

# Reload nginx
nginx -t && nginx -s reload