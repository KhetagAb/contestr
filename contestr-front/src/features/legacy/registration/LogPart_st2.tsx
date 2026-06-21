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
                        <h2>Ожидание подтверждения</h2>
                        <p>
                            Подтвердите регистрацию и нажмите «Готово».
                        </p>
                        <button onClick={submitHandler}>Готово</button>
                    </div>
                </section>
            </div>
        // </div>
    )
}

export default LogPart_st2