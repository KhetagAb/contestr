import "./App.css";
import { Sidebar } from "./cont_compon/SideBar.tsx";
import Tables from "./cont_compon/Tables.tsx";
import AdminLogin from "./AdminLogin.tsx";
import { AdminSessionProvider, useAdminSession } from "./AdminSessionContext.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { client } from "./client/client.gen.ts";

const queryClient = new QueryClient();

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
            {isAdminPage ? <AdminLogin /> : <Tables />}
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
