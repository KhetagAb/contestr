import './App.css';
import {Sidebar} from "./cont_compon/SideBar.tsx";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { client } from './client/client.gen.ts';
// import TournamentSchedule from "../match-integration/src/Tournament-2.tsx";
import RecentParcelsTable from "./cont_compon/recentParcelsTable.tsx";
const queryClient = new QueryClient()

const baseUrl = localStorage.getItem("baseUrl") ?? "/"

client.setConfig({
    baseUrl: baseUrl
})


export default function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <Sidebar/>
            {/*<Tables/>*/}
            <RecentParcelsTable contestStartTime={0}/>
            {/*<TournamentSchedule/>*/}
        </QueryClientProvider>
    )
}


