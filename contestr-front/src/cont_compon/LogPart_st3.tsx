type Props = {
    onRetry: () => void,
}

export const LogPart_st3: React.FC<Props> = ({onRetry}) => {
    return (
        <div className="full_auth">
            <div className="mac-but">
                <span className="btn blue"></span>
                <span className="btn gray"></span>
                <span className="btn gray"></span>
            </div>
            <div className="auth">
                <h1>Авторизация</h1>
            </div>
            <section className="log_section">
                <div id={"log"} className="register">
                    <h2>Произошла ошибка</h2>
                    <p>Обновите страницу и повторите попытку</p>
                    <button onClick={onRetry}>Обновить</button>
                    {/*<image src = "attencion.png"></image>*/}
                </div>
            </section>
        </div>
    )
}

export default LogPart_st3;