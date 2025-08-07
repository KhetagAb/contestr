import './App.css';
import Tables from "./cont_compon/Tables.tsx"
import {Sidebar} from "./cont_compon/SideBar.tsx";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient()


export default function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <Sidebar/>
            <Tables/>
        </QueryClientProvider>
    )
}


