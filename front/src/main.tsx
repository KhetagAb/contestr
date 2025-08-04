import {createRoot} from "react-dom/client";
import App from "./App";
import './cont_compon/Test.css';


const container = document.getElementById("root");

const root = createRoot(container!);

root.render(
    <>
        <div className="masked-bg"></div>
        <div className="masked-dg-2"></div>
        <div id="MyDiv">
            <App/>
        </div>
    </>
)
;