import {cinemas} from './assets/mock';
import {Cinema} from "../cinema/component.tsx";

export const Cinemas = () => {
    return <div>
        {cinemas.map((cinema) =>
            (<Cinema key = {cinema.id} cinema={cinema}/>)
        )}
    </div>
}
// export function Cinemas(){
//     return <div>
//         {cinemas.map((cinema) =>
//             (<Cinema key={cinema.id} cinema = {cinema} />
//             ))}
//     </div>
// } тоже что и сверху
