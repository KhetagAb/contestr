import {LogPart_st1} from "./LogPart_st1.tsx";
import {LogPart_st3} from "./LogPart_st3.tsx";
import {LogPart_st2} from "./LogPart_st2.tsx";
import {useState} from "react";

const Register = () => {
    const [step, setStep] = useState(1)

    const goToStep = (newStep: number) => {
        setStep(newStep)
    }

    return (
        <div>
            {step === 1 && <LogPart_st1 onNext={() => goToStep(2)} />}
            {step === 2 && <LogPart_st2 onNext={() => goToStep(3)} onError={() => goToStep(99)} />}
            {step === 3 && <div>Регистрация завершена</div>}
            {step === 99 && <LogPart_st3 onRetry={() => goToStep(1)} />}
                </div>
    )
}

export default Register;