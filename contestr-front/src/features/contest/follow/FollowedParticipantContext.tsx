import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useMemo,
    useState,
    type ReactNode,
} from "react";
import { useSearchParam } from "react-use";
import {
    readFollowedParticipantId,
    writeFollowedParticipantId,
} from "./storage";

type FollowedParticipantContextValue = {
    contestId: number | null;
    followedParticipantId: string | null;
    setFollowedParticipantId: (participantId: string | null) => void;
    clearFollowedParticipant: () => void;
    isParticipantPickerOpen: boolean;
    openParticipantPicker: () => void;
    closeParticipantPicker: () => void;
};

const FollowedParticipantContext =
    createContext<FollowedParticipantContextValue | null>(null);

export function FollowedParticipantProvider({ children }: { children: ReactNode }) {
    const contestIdRaw = useSearchParam("contestId") || "";
    const contestId = parseInt(contestIdRaw, 10);
    const validContestId =
        Number.isFinite(contestId) && contestId > 0 ? contestId : null;

    const [followedParticipantId, setFollowedState] = useState<string | null>(null);
    const [isParticipantPickerOpen, setParticipantPickerOpen] = useState(false);

    useEffect(() => {
        if (validContestId == null) {
            setFollowedState(null);
            return;
        }
        setFollowedState(readFollowedParticipantId(validContestId));
    }, [validContestId]);

    const setFollowedParticipantId = useCallback(
        (participantId: string | null) => {
            if (validContestId == null) {
                return;
            }
            writeFollowedParticipantId(validContestId, participantId);
            setFollowedState(participantId?.trim() || null);
        },
        [validContestId],
    );

    const clearFollowedParticipant = useCallback(() => {
        setFollowedParticipantId(null);
    }, [setFollowedParticipantId]);

    const openParticipantPicker = useCallback(() => {
        setParticipantPickerOpen(true);
    }, []);

    const closeParticipantPicker = useCallback(() => {
        setParticipantPickerOpen(false);
    }, []);

    const value = useMemo(
        () => ({
            contestId: validContestId,
            followedParticipantId:
                validContestId != null ? followedParticipantId : null,
            setFollowedParticipantId,
            clearFollowedParticipant,
            isParticipantPickerOpen,
            openParticipantPicker,
            closeParticipantPicker,
        }),
        [
            validContestId,
            followedParticipantId,
            setFollowedParticipantId,
            clearFollowedParticipant,
            isParticipantPickerOpen,
            openParticipantPicker,
            closeParticipantPicker,
        ],
    );

    return (
        <FollowedParticipantContext.Provider value={value}>
            {children}
        </FollowedParticipantContext.Provider>
    );
}

export function useFollowedParticipant(): FollowedParticipantContextValue {
    const ctx = useContext(FollowedParticipantContext);
    if (!ctx) {
        throw new Error(
            "useFollowedParticipant must be used within FollowedParticipantProvider",
        );
    }
    return ctx;
}
