import './App.css';
// import Tables from "./cont_compon/Tables.tsx"
import {Sidebar} from "./cont_compon/SideBar.tsx";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { client } from './client/client.gen.ts';
// import TournamentSchedule from "../match-integration/src/Tournament-2.tsx";
import Tables from "./cont_compon/Tables.tsx"
import parcelsData from "./cont_compon/mocks/parcelsData.json"
import RecentParcelsTable from "./cont_compon/recentParcelsTable.tsx";
const queryClient = new QueryClient()

client.setConfig({
    baseUrl: "http://contestr.d.lksh.ru:8080"
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


