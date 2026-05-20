/** Единый шаблон расписания: шесть туров с нарастающей длительностью. */
export const DURATION_TEMPLATE_MINUTES = [15, 20, 25, 30, 40, 60] as const;

export const DURATION_TEMPLATE_LABEL = DURATION_TEMPLATE_MINUTES.map((m) => `${m}`).join(
    " — ",
) + " мин.";
