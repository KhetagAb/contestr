import jwt from "jsonwebtoken";
import { config, jwtEnabled } from "./config.js";
export function verifyAdminToken(tokenRaw) {
    if (!jwtEnabled()) {
        return { ok: false, error: "admin authentication is disabled on server" };
    }
    const token = tokenRaw.replace(/^Bearer\s+/i, "").trim();
    if (!token) {
        return { ok: false, error: "missing token" };
    }
    try {
        const decoded = jwt.verify(token, config.jwtSecret, {
            algorithms: ["HS256"],
        });
        const username = decoded.username ?? decoded.sub;
        if (!username) {
            return { ok: false, error: "invalid token claims" };
        }
        return { ok: true };
    }
    catch {
        return { ok: false, error: "invalid or expired token" };
    }
}
