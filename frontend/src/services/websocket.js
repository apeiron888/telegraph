class WebSocketService {
    constructor() {
        this.ws = null;
        this.reconnectInterval = null;
        this.messageHandlers = [];
        this.isConnecting = false;
    }

    connect(token) {
        if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
            return;
        }

        this.isConnecting = true;
        const wsUrl = `ws://localhost:8080/api/v1/ws`;

        try {
            this.ws = new WebSocket(wsUrl);

            // Send auth token after connection opens
            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.isConnecting = false;

                // Send authorization (token is already in localStorage, server reads it from HTTP upgrade)
                if (this.reconnectInterval) {
                    clearInterval(this.reconnectInterval);
                    this.reconnectInterval = null;
                }
            };

            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    this.messageHandlers.forEach(handler => handler(data));
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.isConnecting = false;
            };

            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.isConnecting = false;
                this.scheduleReconnect(token);
            };
        } catch (error) {
            console.error('Error creating WebSocket:', error);
            this.isConnecting = false;
            this.scheduleReconnect(token);
        }
    }

    scheduleReconnect(token) {
        if (!this.reconnectInterval) {
            this.reconnectInterval = setInterval(() => {
                if (!this.isConnecting && (!this.ws || this.ws.readyState === WebSocket.CLOSED)) {
                    console.log('Attempting to reconnect WebSocket...');
                    this.connect(token);
                }
            }, 5000); // Retry every 5 seconds
        }
    }

    disconnect() {
        if (this.reconnectInterval) {
            clearInterval(this.reconnectInterval);
            this.reconnectInterval = null;
        }
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    onMessage(handler) {
        this.messageHandlers.push(handler);
        return () => {
            this.messageHandlers = this.messageHandlers.filter(h => h !== handler);
        };
    }
}

export const wsService = new WebSocketService();
