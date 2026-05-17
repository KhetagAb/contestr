import { useEffect, useState, type FormEvent } from "react";
import { adminAuthHeaders, clearAdminToken, getAdminToken, setAdminToken } from "./adminAuth";
import AdminTimetable from "./AdminTimetable";
import { Sidebar } from "./cont_compon/SideBar.tsx";
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

async function fetchAdminMe(): Promise<{ ok: true; username: string } | { ok: false; message: string }> {
    const meRes = await fetch("/api/admin/me", { headers: adminAuthHeaders() });
    const meBody = (await meRes.json()) as MeResponse & ErrorResponse;
    if (!meRes.ok || !meBody.ok) {
        return { ok: false, message: meBody.message ?? "Сессия недействительна" };
    }
    return { ok: true, username: meBody.username };
}

export default function AdminLogin() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [loggedInAs, setLoggedInAs] = useState<string | null>(null);
    const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
    const [message, setMessage] = useState("");

    useEffect(() => {
        const token = getAdminToken();
        if (!token) {
            return;
        }

        setStatus("loading");
        fetchAdminMe()
            .then((result) => {
                if (result.ok) {
                    setLoggedInAs(result.username);
                    setStatus("success");
                } else {
                    clearAdminToken();
                    setStatus("idle");
                }
            })
            .catch(() => {
                clearAdminToken();
                setStatus("idle");
            });
    }, []);

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setStatus("loading");
        setMessage("");

        try {
            const loginRes = await fetch("/api/admin/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
            });

            const loginBody = (await loginRes.json()) as LoginResponse & ErrorResponse;
            if (!loginRes.ok) {
                setStatus("error");
                setMessage(loginBody.message ?? "Ошибка авторизации");
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

            setLoggedInAs(meResult.username);
            setStatus("success");
        } catch {
            setStatus("error");
            setMessage("Сервер недоступен");
        }
    };

    const handleLogout = () => {
        clearAdminToken();
        setLoggedInAs(null);
        setStatus("idle");
        setMessage("");
        setPassword("");
    };

    return (
        <>
            <Sidebar />
            <div className="admin-login-content">
                <div className={`admin-login-panel ${loggedInAs ? "admin-console-panel" : ""}`}>
                    {loggedInAs ? (
                        <>
                            <div className="admin-login-form-box">
                                <p className="admin-login-greeting">
                                    Вы вошли как <span className="h1_pink">{loggedInAs}</span>
                                </p>
                                <button type="button" className="admin-login-btn" onClick={handleLogout}>
                                    Выйти
                                </button>
                                <a href="/" className="admin-login-link">
                                    На главную
                                </a>
                            </div>
                            <AdminTimetable />
                        </>
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
                                            value={username}
                                            onChange={(e) => setUsername(e.target.value)}
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
        </>
    );
}
