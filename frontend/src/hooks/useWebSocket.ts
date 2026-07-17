import { useEffect, useRef, useState } from 'react';

export function useWebSocket(onMessage?: (data: unknown) => void) {
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = import.meta.env.VITE_WS_URL || `${proto}://${window.location.hostname}:8080/ws`;
    const ws = new WebSocket(host);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);
    ws.onmessage = (evt) => {
      try {
        const parsed = JSON.parse(evt.data);
        onMessage?.(parsed);
      } catch {
        onMessage?.(evt.data);
      }
    };

    return () => ws.close();
  }, [onMessage]);

  return { connected, socket: wsRef.current };
}
