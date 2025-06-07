import {reviews} from "../cinemas/assets/mock.tsx";

export const Cinema = ({cinema}) => {
    const {name, movies} = cinema;
    return (
        <div>
            <h2>{name}</h2>
            {/*<Cinema cinema={cinema}/>*/}

            <div>
                <h2>reviews</h2>
                {reviews.map((review) => (
                    <>
                        <div>
                            <h4>{review.user.id} {review.user.name} </h4>
                        </div>
                        <div>
                            {review.text}
                            <p>Оценка фильма: {review.rating}</p>
                        </div>
                    </>
                ))}
            </div>

        </div>
    )
};

// struct point {
//     int x, y;
// }
//
// int get_count(point a) {
//     return a.x;
// }


//