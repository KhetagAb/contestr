import './App.css';
import {Sidebar} from "./cont_compon/SideBar.tsx";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { client } from './client/client.gen.ts';
// import TournamentSchedule from "../match-integration/src/Tournament-2.tsx";
import parcelsData from "./cont_compon/mocks/parcelsData.json"
import RecentParcelsTable from "./cont_compon/recentParcelsTable.tsx";
const queryClient = new QueryClient()

client.setConfig({
    baseUrl: "/"
})


export default function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <Sidebar/>
            {/*<Tables/>*/}
            <RecentParcelsTable data={parcelsData}  contestStartTime={0}/>
            {/*<TournamentSchedule/>*/}
        </QueryClientProvider>
    )
}


