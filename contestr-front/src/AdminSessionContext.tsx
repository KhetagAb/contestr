import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useMemo,
    useState,
    type ReactNode,
} from "react";
import { adminAuthHeaders, clearAdminToken, getAdminToken } from "./adminAuth";
import type { AdminSidebarSession } from "./cont_compon/SideBar.tsx";

type MeResponse = {
    username: string;
    ok: boolean;
};

type ErrorResponse = {
    message: string;
};

type AdminSessionContextValue = {
    username: string | null;
    isLoading: boolean;
    sidebarSession: AdminSidebarSession | null;
    setUsername: (username: string | null) => void;
    logout: () => void;
    refreshSession: () => Promise<void>;
};

const AdminSessionContext = createContext<AdminSessionContextValue | null>(null);

async function fetchAdminMe(): Promise<{ ok: true; username: string } | { ok: false }> {
    const meRes = await fetch("/api/admin/me", { headers: adminAuthHeaders() });
    const meBody = (await meRes.json()) as MeResponse & ErrorResponse;
    if (!meRes.ok || !meBody.ok) {
        return { ok: false };
    }
    return { ok: true, username: meBody.username };
}

export function AdminSessionProvider({ children }: { children: ReactNode }) {
    const [username, setUsername] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const refreshSession = useCallback(async () => {
        const token = getAdminToken();
        if (!token) {
            setUsername(null);
            setIsLoading(false);
            return;
        }

        setIsLoading(true);
        try {
            const result = await fetchAdminMe();
            if (result.ok) {
                setUsername(result.username);
            } else {
                clearAdminToken();
                setUsername(null);
            }
        } catch {
            clearAdminToken();
            setUsername(null);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        void refreshSession();
    }, [refreshSession]);

    const logout = useCallback(() => {
        clearAdminToken();
        setUsername(null);
    }, []);

    const sidebarSession = useMemo<AdminSidebarSession | null>(
        () => (username ? { username, onLogout: logout } : null),
        [username, logout]
    );

    const value = useMemo(
        () => ({
            username,
            isLoading,
            sidebarSession,
            setUsername,
            logout,
            refreshSession,
        }),
        [username, isLoading, sidebarSession, logout, refreshSession]
    );

    return <AdminSessionContext.Provider value={value}>{children}</AdminSessionContext.Provider>;
}

export function useAdminSession() {
    const context = useContext(AdminSessionContext);
    if (!context) {
        throw new Error("useAdminSession must be used within AdminSessionProvider");
    }
    return context;
}
