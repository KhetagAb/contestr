
// tournament.queries.ts
import { useQuery } from '@tanstack/react-query';
import type { Activity, TeamList } from '../client/types.gen.ts';
import { mockActivities, mockTeams } from './tournament.mocks';

export const useActivities = (sportSectionId: number) => {
    return useQuery<Activity[], Error>({
        queryKey: ['activities', sportSectionId],
        queryFn: async () => {
            console.log('🔄 Запрос активностей для секции:', sportSectionId);
            // Имитируем задержку сети
            await new Promise(resolve => setTimeout(resolve, 1000));
            return mockActivities;
        },
        enabled: !!sportSectionId,
    });
};

export const useTeamsLazy = (activityId: number, enabled: boolean = true) => {
    return useQuery<TeamList, Error>({
        queryKey: ['teams', activityId],
        queryFn: async () => {
            console.log('🔄 Запрос команд для активности:', activityId);
            // Имитируем задержку сети
            await new Promise(resolve => setTimeout(resolve, 800));
            return mockTeams;
        },
        enabled: !!activityId && enabled,
    });
};