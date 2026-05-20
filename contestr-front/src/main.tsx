import { createRoot } from "react-dom/client";
import "@/app/styles/index.css";
import App from "@/app/App";

const container = document.getElementById("root");

const root = createRoot(container!);

root.render(
    <div id="MyDiv">
        <App />
    </div>
);
