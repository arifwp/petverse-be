# Chat Module

Owns private and community chat behavior: conversations, participants, messages, delivery state, and real-time fanout.

Treat real-time transport as an adapter. The core chat use cases should not know whether delivery happens through WebSocket, Server-Sent Events, or a queue.
