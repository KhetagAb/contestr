export const SOCKET_EVENTS = {
    contestJoin: "contest:join",
    announcementEmit: "announcement:emit",
    announcement: "announcement",
};
export function contestRoom(contestId) {
    return `contest:${contestId}`;
}
