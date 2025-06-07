type Props = {
    onNext: () => void
    onError: () => void
}

export const LogPart_st2: React.FC<Props> = ({onNext, onError}) => {
    const submitHandler = () => {
        const success = true
        if (success) {
            onNext()
        } else {
            onError()
        }
    }


    return (
        // <div className={"container_row"}>
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
                        <h2>Ждем подтверждения в Telegram</h2>
                        <p>
                            Поделитесь контактом с ботом в Telegram c аккаунта с номером +7 (999) 999-99-99 и нажмите
                            «Готово»
                        </p>
                        <p>Перейти в <a>Telegram @contestrBot</a></p>
                        <button onClick={submitHandler}>Готово</button>
                    </div>
                </section>
            </div>
        // </div>
    )
}

export default LogPart_st2