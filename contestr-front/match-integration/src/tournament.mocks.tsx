// tournament.mocks.ts
import type { Activity, TeamList } from '../client/types.gen.ts';

export const mockActivities: Activity[] = [
    {
        id: 1,
        title: "Турнир по футболу",
        description: "Еженедельный турнир для любителей",
        creator: {
            coreId: 1,
            tgId: 123456789
        }
    }
];

export const mockTeams: TeamList = [
    {
        id: 1,  // Скорпионы
        name: "Скорпионы",
        captain: {
            coreId: 2,
            tgId: 111111111
        },
        members: [
            { coreId: 4, tgId: 222222222 },
            { coreId: 5, tgId: 333333333 },
            { coreId: 6, tgId: 444444444 }
        ]
    },
    {
        id: 2,  // Молот
        name: "Молот",
        captain: {
            coreId: 3,
            tgId: 555555555
        },
        members: [
            { coreId: 7, tgId: 666666666 },
            { coreId: 8, tgId: 777777777 }
        ]
    },
    {
        id: 3,  // Витязи
        name: "Витязи",
        captain: {
            coreId: 10,
            tgId: 888888888
        },
        members: [
            { coreId: 11, tgId: 999999999 },
            { coreId: 12, tgId: 101010101 }
        ]
    }
];