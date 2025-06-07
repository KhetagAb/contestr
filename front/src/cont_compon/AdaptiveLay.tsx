import Register from "./Registration.tsx";
import {MainPart} from "./MainPart.tsx";
import {Desktop, Tablet, Mobile} from "./layout.tsx"

export const AdaptiveLay = () => {
    return (
        <>
            <Desktop>
                <div>
                    <MainPart/>
                    <Register/>
                </div>
            </Desktop>

            <Tablet>
                <div>
                    <MainPart/>
                    <Register/>
                </div>
            </Tablet>

            <Mobile>
                <div>
                    <MainPart/>
                    <Register/>
                </div>
            </Mobile>
        </>
    )
}