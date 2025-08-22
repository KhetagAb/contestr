import styles from "./tour.module.css"
import { useState } from 'react'
import { ChevronDown, ChevronUp } from 'lucide-react'
import { PiStarFill } from "react-icons/pi";
import telegramLogo from "../public/telegramLogo.svg"
import { useActivities, useTeamsLazy } from './tournament.queries';
import type { Activity, Team, RegisterPlayerRequest } from '../client';

interface TournamentProps {
    playersData: RegisterPlayerRequest[]; // массив с данными всех игроков
}

const Tournament = ({ playersData }: TournamentProps) => {
    const [openId, setOpenId] = useState<number | null>(null);
    const FOOTBALL_SPORT_ID = 1;

    const { data: activities = [] } = useActivities(FOOTBALL_SPORT_ID);

    const toggleContainer = (id: number) => {
        setOpenId(openId === id ? null : id);
    };

    return (
        <div className={styles.tournament}>
            <h1>Спортивные активности</h1>
            <div className={styles.tournamentsRow}>
                {activities.map((activity) => (
                    <TournamentCard
                        key={activity.id}
                        activity={activity}
                        isOpen={openId === activity.id}
                        onToggle={() => toggleContainer(activity.id)}
                        playersData={playersData}
                    />
                ))}
            </div>
        </div>
    );
};

interface TournamentCardProps {
    activity: Activity;
    isOpen: boolean;
    onToggle: () => void;
    playersData: RegisterPlayerRequest[];
}

const TournamentCard = ({ activity, isOpen, onToggle, playersData }: TournamentCardProps) => {
    const { data: teams = [], isLoading: teamsLoading } = useTeamsLazy(activity.id, isOpen);

    return (
        <div className={styles.tourContainer}>
            <p className={styles.tourName}>{activity.title}</p>
            <p className={styles.tourInfo}>{activity.description || "Описание отсутствует"}</p>
            <div className={styles.tourAdmin}>
                <p>
                    Организатор активности: {activity.creator?.tgId || "Информация отсутствует"}
                    <img src={telegramLogo} className={styles.logo} alt="Telegram" />
                </p>
            </div>
            <div className={styles.toggleSection}>
                <button onClick={onToggle} className={styles.toggleButton}>
                    {isOpen ? <ChevronUp /> : <>Команды-участники <ChevronDown /></>}
                </button>
                {isOpen && (
                    <div className={`${styles.tourLists} ${isOpen ? styles.open : ""}`}>
                        {teamsLoading ? (
                            <p>Загрузка команд...</p>
                        ) : teams.length > 0 ? (
                            <TeamsList teams={teams} playersData={playersData} />
                        ) : (
                            <p>Нет зарегистрированных команд</p>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

interface TeamsListProps {
    teams: Team[];
    playersData: RegisterPlayerRequest[];
}

const TeamsList = ({ teams, playersData }: TeamsListProps) => {
    // функция для поиска имени по tg_id
    const getPlayerName = (tgId: number) => {
        const player = playersData.find(p => p.tg_id === tgId);
        return player?.name || "Игрок";
    };

    return (
        <div className={styles.teamsRow}>
            {teams.map((team) => (
                <div key={team.id} className={styles.teamCard}>
                    <h3>{team.id} — {team.name}</h3>
                    <p>
                        {getPlayerName(team.captain.tgId)}
                        <PiStarFill className={styles.logo} color={"#BF9298"} />
                    </p>
                    <ul>
                        {team.members && team.members.length > 0 ? (
                            team.members.map((member) => (
                                <li key={member.tgId}>{getPlayerName(member.tgId)}</li>
                            ))
                        ) : (
                            <li>Нет участников</li>
                        )}
                    </ul>
                </div>
            ))}
        </div>
    );
};

export default Tournament;
