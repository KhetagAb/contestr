import {useState} from "react";

type Props = {
    onNext: () => void;
}

export const LogPart_st1: React.FC<Props> = ({onNext}) => {
    const [phone] = useState('')

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (phone.includes('_')) {
            alert('Заполните номер полностью')
            return
        }

        onNext()
    }

    return (
        <div className="full_auth">
            <div className="mac-but">
                <span className="btn blue"></span>
                <span className="btn gray"></span>
                <span className="btn gray"></span>
            </div>

            <div className="auth">
                <h2>Авторизация</h2>
            </div>

            <section className="log_section">

                <div id="log" className="register">
                    <h2>Введите телефон</h2>
                    <form onSubmit={handleSubmit}>
                        <input></input>
                    </form>
                    <div className="container">
                        <input type="checkbox" id="cbx2" style={{display: 'none'}}/>
                        <label htmlFor="cbx2" className="check">
                            <svg width="18px" height="18px" viewBox="0 0 18 18">
                                <path
                                    d="M 1 9 L 1 9 c 0 -5 3 -8 8 -8 L 9 1 C 14 1 17 5 17 9 L 17 9 c 0 4 -4 8 -8 8 L 9 17 C 5 17 1 14 1 9 L 1 9 Z"/>
                                <polyline points="1 9 7 14 15 4"/>
                            </svg>
                            <span className="check-text">
                                Я согласен на обработку персональных данных
                            </span>
                        </label>
                    </div>
                    <button type="submit" className="TG">
                        Войти через Telegram
                    </button>
                </div>
            </section>
        </div>
    )
}

