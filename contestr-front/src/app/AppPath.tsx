import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from "react";

export const CONTEST_PHASE_FOCUS_PATH = "/phase";

type AppPathContextValue = {
    path: string;
    navigate: (to: string) => void;
};

const AppPathContext = createContext<AppPathContextValue | null>(null);

export function AppPathProvider({ children }: { children: ReactNode }) {
    const [path, setPath] = useState(() => window.location.pathname);

    useEffect(() => {
        const sync = () => setPath(window.location.pathname);
        window.addEventListener("popstate", sync);
        return () => window.removeEventListener("popstate", sync);
    }, []);

    const navigate = useCallback((to: string) => {
        const pathPart = to.split("?")[0].split("#")[0] || "/";
        const searchPart = to.includes("?")
            ? to.slice(to.indexOf("?"), to.includes("#") ? to.indexOf("#") : undefined)
            : window.location.search;
        const hashPart = to.includes("#")
            ? to.slice(to.indexOf("#"))
            : window.location.hash;
        const href = `${pathPart}${searchPart}${hashPart}`;

        if (
            window.location.pathname === pathPart &&
            window.location.search === searchPart &&
            window.location.hash === hashPart
        ) {
            return;
        }
        window.history.pushState({}, "", href);
        setPath(pathPart);
    }, []);

    return (
        <AppPathContext.Provider value={{ path, navigate }}>
            {children}
        </AppPathContext.Provider>
    );
}

export function useAppPath(): AppPathContextValue {
    const value = useContext(AppPathContext);
    if (!value) {
        throw new Error("useAppPath must be used within AppPathProvider");
    }
    return value;
}
