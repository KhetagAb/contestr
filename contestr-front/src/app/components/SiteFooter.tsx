import styles from "./SiteFooter.module.css";

export function SiteFooter() {
    return (
        <footer className={styles.footer}>
            <p className={styles.text}>
                <span className={styles.brand}>Regatta</span>
                <span className={styles.separator} aria-hidden="true">
                    ·
                </span>
                Khetag Dzestelov, Daniil Golov · 2026
            </p>
        </footer>
    );
}
