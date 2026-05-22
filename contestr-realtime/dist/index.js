import express from "express";
import { createServer } from "node:http";
import { Server } from "socket.io";
import { config } from "./config.js";
import { registerSocketHandlers } from "./handlers/connection.js";
const app = express();
app.get("/health", (_req, res) => {
    res.json({ ok: true });
});
const httpServer = createServer(app);
const io = new Server(httpServer, {
    cors: {
        origin: config.corsOrigin,
        credentials: true,
    },
    path: "/socket.io/",
});
registerSocketHandlers(io);
httpServer.listen(config.port, () => {
    console.log(`contestr-realtime listening on :${config.port}`);
});
