import "@/app/styles/App.css";
import { QueryClientProvider } from "@tanstack/react-query";
import { client } from "@/client/client.gen";
import AdminLogin from "@/features/admin/pages/AdminLogin";
import ContestStandingsPage from "@/features/contest/pages/ContestStandingsPage";
import { Sidebar } from "@/features/contest/components/sidebar/Sidebar";
import { AdminSessionProvider, useAdminSession } from "@/app/providers/AdminSessionContext";
import { queryClient } from "@/app/providers/queryClient";

const baseUrl = localStorage.getItem("baseUrl") ?? "/";

client.setConfig({
    baseUrl: baseUrl,
});

function AppShell() {
    const { sidebarSession } = useAdminSession();
    const isAdminPage = window.location.pathname.startsWith("/admin");

    return (
        <>
            <Sidebar adminSession={sidebarSession} />
            {isAdminPage ? <AdminLogin /> : <ContestStandingsPage />}
        </>
    );
}

export default function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <AdminSessionProvider>
                <AppShell />
            </AdminSessionProvider>
        </QueryClientProvider>
    );
}
