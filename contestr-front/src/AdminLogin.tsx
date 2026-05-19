import { useState, type FormEvent } from "react";
import { adminAuthHeaders, clearAdminToken, setAdminToken } from "./adminAuth";
import { useAdminSession } from "./AdminSessionContext.tsx";
import TimetablePage from "./admin-timetable/TimetablePage";
import "./App.css";
import "./AdminLogin.css";

type LoginResponse = {
    access_token: string;
    expires_in: number;
};

type MeResponse = {
    username: string;
    ok: boolean;
};

type ErrorResponse = {
    message: string;
};

const AUTH_MESSAGE_TRANSLATIONS: Record<string, string> = {
    "admin authentication is disabled": "Авторизация администратора отключена",
    "failed to issue token": "Не удалось выпустить токен",
    "internal server error": "Внутренняя ошибка сервера",
    "invalid or expired token": "Токен некорректен или истёк",
    "invalid request body": "Некорректное тело запроса",
    "invalid username or password": "Неверный логин или пароль",
    "missing authorization header": "Не передан заголовок авторизации",
    "unauthorized": "Не авторизован",
};

function translateAdminAuthMessage(message?: string, fallback = "Ошибка авторизации") {
    const value = message?.trim();
    if (!value) {
        return fallback;
    }

    if (/[А-Яа-яЁё]/.test(value)) {
        return value;
    }

    return AUTH_MESSAGE_TRANSLATIONS[value.toLowerCase()] ?? fallback;
}

async function fetchAdminMe(): Promise<{ ok: true; username: string } | { ok: false; message: string }> {
    const meRes = await fetch("/api/admin/me", { headers: adminAuthHeaders() });
    const meBody = (await meRes.json()) as MeResponse & ErrorResponse;
    if (!meRes.ok || !meBody.ok) {
        return { ok: false, message: translateAdminAuthMessage(meBody.message, "Сессия недействительна") };
    }
    return { ok: true, username: meBody.username };
}

export default function AdminLogin() {
    const { username, setUsername } = useAdminSession();
    const [loginUsername, setLoginUsername] = useState("");
    const [password, setPassword] = useState("");
    const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
    const [message, setMessage] = useState("");

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setStatus("loading");
        setMessage("");

        try {
            const loginRes = await fetch("/api/admin/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username: loginUsername, password }),
            });

            const loginBody = (await loginRes.json()) as LoginResponse & ErrorResponse;
            if (!loginRes.ok) {
                setStatus("error");
                setMessage(translateAdminAuthMessage(loginBody.message));
                return;
            }

            setAdminToken(loginBody.access_token);

            const meResult = await fetchAdminMe();
            if (!meResult.ok) {
                clearAdminToken();
                setStatus("error");
                setMessage(meResult.message);
                return;
            }

            setUsername(meResult.username);
            setStatus("success");
        } catch {
            setStatus("error");
            setMessage("Сервер недоступен");
        }
    };

    return (
        <div className={`admin-login-content${username ? " admin-login-content--console" : ""}`}>
            <div className={`admin-login-panel ${username ? "admin-console-panel" : ""}`}>
                {username ? (
                    <TimetablePage />
                ) : (
                    <>
                        <h1 className="admin-login-title">Авторизация</h1>
                        <form onSubmit={handleSubmit} className="admin-login-form">
                            <div className="admin-login-form-box">
                                {message && status === "error" && (
                                    <p
                                        className="admin-login-message admin-login-message--error"
                                        role="alert"
                                    >
                                        Ошибка: {message}
                                    </p>
                                )}
                                <label className="admin-login-field">
                                    <span>Логин</span>
                                    <input
                                        type="text"
                                        value={loginUsername}
                                        onChange={(e) => setLoginUsername(e.target.value)}
                                        autoComplete="username"
                                        required
                                    />
                                </label>
                                <label className="admin-login-field">
                                    <span>Пароль</span>
                                    <input
                                        type="password"
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        autoComplete="current-password"
                                        required
                                    />
                                </label>
                                <button type="submit" className="admin-login-btn" disabled={status === "loading"}>
                                    {status === "loading" ? "Проверка…" : "Войти"}
                                </button>
                            </div>
                        </form>
                    </>
                )}
            </div>
        </div>
    );
}
