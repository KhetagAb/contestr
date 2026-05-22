import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from "react";

export const CONTEST_PHASE_FOCUS_PATH = "/phase";

export const CONTEST_HOME_QUERY = "home";

type AppPathContextValue = {
    path: string;
    search: string;
    navigate: (to: string) => void;
};

const AppPathContext = createContext<AppPathContextValue | null>(null);

function readLocation() {
    return {
        path: window.location.pathname,
        search: window.location.search,
    };
}

export function AppPathProvider({ children }: { children: ReactNode }) {
    const [location, setLocation] = useState(readLocation);

    useEffect(() => {
        const sync = () => setLocation(readLocation());
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
        setLocation({ path: pathPart, search: searchPart });
    }, []);

    return (
        <AppPathContext.Provider
            value={{ path: location.path, search: location.search, navigate }}
        >
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
