import { verifyAdminToken } from "../auth.js";
import { SOCKET_EVENTS, contestRoom, } from "../protocol.js";
function isValidContestId(id) {
    return typeof id === "number" && Number.isFinite(id) && id > 0;
}
function isAnnouncementPayload(raw) {
    if (!raw || typeof raw !== "object") {
        return false;
    }
    const p = raw;
    if (typeof p.id !== "string" || !p.id.trim()) {
        return false;
    }
    if (typeof p.title !== "string" || !p.title.trim()) {
        return false;
    }
    if (typeof p.caption !== "string") {
        return false;
    }
    if (p.titleVariant !== "tour" && p.titleVariant !== "points") {
        return false;
    }
    if (typeof p.confetti !== "boolean") {
        return false;
    }
    if (p.audience !== "all" && p.audience !== "participant") {
        return false;
    }
    if (p.audience === "participant" && typeof p.participantId !== "string") {
        return false;
    }
    if (p.holdMs != null && (typeof p.holdMs !== "number" || p.holdMs <= 0)) {
        return false;
    }
    return true;
}
export function registerSocketHandlers(io) {
    io.on("connection", (socket) => {
        socket.on(SOCKET_EVENTS.contestJoin, (data) => {
            if (!isValidContestId(data?.contestId)) {
                return;
            }
            socket.join(contestRoom(data.contestId));
        });
        socket.on(SOCKET_EVENTS.announcementEmit, (data, ack) => {
            const reply = (result) => {
                if (typeof ack === "function") {
                    ack(result);
                }
            };
            if (!isValidContestId(data?.contestId)) {
                reply({ ok: false, error: "invalid contestId" });
                return;
            }
            const auth = verifyAdminToken(data?.token ?? "");
            if (!auth.ok) {
                reply({ ok: false, error: auth.error });
                return;
            }
            if (!isAnnouncementPayload(data?.payload)) {
                reply({ ok: false, error: "invalid announcement payload" });
                return;
            }
            socket.join(contestRoom(data.contestId));
            io.to(contestRoom(data.contestId)).emit(SOCKET_EVENTS.announcement, data.payload);
            reply({ ok: true });
        });
    });
}
