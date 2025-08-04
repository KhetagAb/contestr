// import {AdaptiveLay} from "./cont_compon/AdaptiveLay.tsx";
import './cont_compon/Test.css';
import {MainPart} from "./cont_compon/MainPart.tsx";
import Register from "./cont_compon/Registration.tsx";
// import LogPart_st2 from "./cont_compon/LogPart_st2.tsx";
import {Sidebar} from "./cont_compon/SideBar.tsx";
// import LogPart_st3 from "./cont_compon/LogPart_st3.tsx";
// import {MainWrapper} from "./cont_compon/Wrapper.tsx";

export default function App() {
    return (
        <>
            {/*<AdaptiveLay/>*/}
            <Sidebar/>
            <MainPart/>
            {/*<LogPart_st3 onRetry={function(): void {*/}
            {/*    throw new Error('Function not implemented.');*/}
            {/*} }/>*/}
            <Register/>
            {/*<LogPart_st2 onNext={function(): void {*/}
            {/*    throw new Error('Function not implemented.');*/}
            {/*} } onError={function(): void {*/}
            {/*    throw new Error('Function not implemented.');*/}
            {/*} }/>*/}
            {/*<MainWrapper/>*/}
        </>
    )
}


