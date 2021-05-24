Copy this file into `/etc/systemd/system` as root:

```bash
sudo cp environment-monitor.service /etc/systemd/system/environment-monitor.service
```

Start the service using the following command:
```bash
sudo systemctl start environment-monitor.service
```

Check the service status using:
```bash
sudo systemctl status environment-monitor.service
```

Stop using the folloging command:
```bash
sudo systemctl start environment-monitor.service
```

If everything goes well, set the service to start automatically on reboot by using this command:
```bash
sudo systemctl enable environment-monitor.service
```

