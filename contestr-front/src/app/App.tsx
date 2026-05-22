import "@/app/styles/App.css";
import { useCallback, useEffect, useState } from "react";
import { QueryClientProvider } from "@tanstack/react-query";
import { client } from "@/client/client.gen";
import { AppPathProvider, CONTEST_PHASE_FOCUS_PATH, useAppPath } from "@/app/AppPath";
import AdminLogin from "@/features/admin/pages/AdminLogin";
import ContestPhaseFocusPage from "@/features/contest/pages/ContestPhaseFocusPage";
import ContestStandingsPage from "@/features/contest/pages/ContestStandingsPage";
import { Sidebar } from "@/features/contest/components/sidebar/Sidebar";
import { AdminSessionProvider, useAdminSession } from "@/app/providers/AdminSessionContext";
import { queryClient } from "@/app/providers/queryClient";
import { FollowedParticipantProvider } from "@/features/contest/follow/FollowedParticipantContext";
import { RegattaRulesModal } from "@/features/contest/components/regatta-rules/RegattaRulesModal";
import {
    hasSeenRegattaRules,
    markRegattaRulesSeen,
} from "@/features/contest/components/regatta-rules/regattaRulesStorage";

const baseUrl = localStorage.getItem("baseUrl") ?? "/";

client.setConfig({
    baseUrl: baseUrl,
});

function RegattaRulesLayer() {
    const [open, setOpen] = useState(false);

    useEffect(() => {
        if (!hasSeenRegattaRules()) {
            setOpen(true);
        }
    }, []);

    const onClose = useCallback(() => {
        markRegattaRulesSeen();
        setOpen(false);
    }, []);

    return <RegattaRulesModal open={open} onClose={onClose} />;
}

function AppShell() {
    const { sidebarSession } = useAdminSession();
    const { path } = useAppPath();
    const isAdminPage = path.startsWith("/admin");
    const isPhaseFocusPage = path === CONTEST_PHASE_FOCUS_PATH;

    let mainContent = <ContestStandingsPage />;
    if (isAdminPage) {
        mainContent = <AdminLogin />;
    } else if (isPhaseFocusPage) {
        mainContent = <ContestPhaseFocusPage />;
    }

    return (
        <FollowedParticipantProvider>
            <Sidebar adminSession={sidebarSession} />
            {mainContent}
            <RegattaRulesLayer />
        </FollowedParticipantProvider>
    );
}

export default function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <AdminSessionProvider>
                <AppPathProvider>
                    <AppShell />
                </AppPathProvider>
            </AdminSessionProvider>
        </QueryClientProvider>
    );
}
