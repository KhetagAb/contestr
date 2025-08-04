import { MainPart } from './MainPart'
import { LogPart_st1 } from './LogPart_st1'

export const MainWrapper = () => {
    const handleNext = () => {
        console.log('next')
    }

    return (
        <div className="wrapper">
            <MainPart />
            <LogPart_st1 onNext={handleNext} />
        </div>
    )
}
