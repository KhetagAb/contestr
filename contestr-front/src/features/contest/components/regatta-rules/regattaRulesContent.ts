export type RegattaRulesSectionId = "structure" | "bonuses";

export type RegattaRulesSection = {
    id: RegattaRulesSectionId;
    title: string;
    items: readonly string[];
};

/** Текст правил для участников (общие формулировки, без привязки к числам контеста). */
export const REGATTA_RULES_SECTIONS: readonly RegattaRulesSection[] = [
    {
        id: "structure",
        title: "Как устроена регата",
        items: [
            "Соревнование идёт турами по расписанию; между турами могут быть перерывы.",
            "В каждом туре участники делятся на группы — вы соревнуетесь за дополнительные баллы с соперниками из вашей команды.",
        ],
    },
    {
        id: "bonuses",
        title: "Начисление очков",
        items: [
            "Дополнительные баллы за решение во время тура — если задача зачтена до конца тура.",
            "Бонус за первенство в группе — если вы первым в группе набрали максимум по задаче до конца тура.",
        ],
    },
];
