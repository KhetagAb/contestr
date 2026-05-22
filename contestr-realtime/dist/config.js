function parseCorsOrigin(raw) {
    const value = raw?.trim() || "*";
    if (value === "*") {
        return "*";
    }
    return value.split(",").map((s) => s.trim()).filter(Boolean);
}
export const config = {
    port: parseInt(process.env.PORT ?? "3001", 10),
    corsOrigin: parseCorsOrigin(process.env.CORS_ORIGIN),
    jwtSecret: process.env.APP_ADMIN_JWT_SECRET?.trim() ?? "",
};
export function jwtEnabled() {
    return config.jwtSecret.length > 0;
}
