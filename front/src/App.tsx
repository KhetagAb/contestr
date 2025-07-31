import './cont_compon/Test.css';
// import {MainPart} from "./cont_compon/MainPart.tsx";
// import Register from "./cont_compon/Registration.tsx";
import {Sidebar} from "./cont_compon/SideBar.tsx";
import ResultsTable from "./cont_compon/tables.tsx";
// import LogPart_st2 from "./cont_compon/LogPart_st2.tsx";
// import {LogPart_st1} from "./cont_compon/LogPart_st1.tsx";

export function App() {

    return (
        <>
            <Sidebar/>
            <ResultsTable />
            {/*<MainPart/>*/}
            {/*/!*<Register/>*!/*/}
            {/*/!*<LogPart_st1 onNext={function(): void {*!/*/}
            {/*/!*    throw new Error('Function not implemented.');*!/*/}
            {/*/!*} }/>*!/*/}
            {/*<LogPart_st2 onNext={function(): void {*/}
            {/*    throw new Error('Function not implemented.');*/}
            {/*} } onError={function(): void {*/}
            {/*    throw new Error('Function not implemented.');*/}
            {/*} }/>*/}
        </>
    )
}


