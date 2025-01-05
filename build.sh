#!/bin/bash

echo "🚀 Starting web deployment..."
sleep 2

echo "🛑 Stopping web service..."
sudo systemctl stop adminBac.service
if [ $? -ne 0 ]; then
    echo "❌ Failed to stop the service! Exiting..."
    exit 1
fi

echo "🗑️ Removing the old main file..."
sudo rm /usr/local/bin/main
if [ $? -ne 0 ]; then
    echo "❌ Failed to remove the old file!"
    exit 1
fi

echo "⚙️ Building the Go file..."
sleep 2
go build main.go
if [ $? -ne 0 ]; then
    echo "❌ Failed to build the Go file!"
    exit 1
fi

echo "📂 Moving the new build file..."
sleep 2
sudo mv main /usr/local/bin/
if [ $? -ne 0 ]; then
    echo "❌ Failed to move the new file!"
    exit 1
fi

echo "🔌 Enabling the service..."
sudo systemctl enable adminBac.service

echo "▶️ Starting the service..."
sleep 2
sudo systemctl start adminBac.service
if [ $? -ne 0 ]; then
    echo "❌ Failed to start the service!"
    exit 1
fi

echo "📊 Checking the service status..."
sleep 2
sudo systemctl status adminBac.service

echo "✅ Web deployment completed successfully! 🎉"
